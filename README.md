# Voting System in Go

This project implements a simple voting system using Go with a REST API, token authentication, and Swagger documentation.

## Features

### REST API

- User registration with encrypted passwords (bcrypt)
- Login/Logout with token authentication
- Candidate listing
- Voting system (one vote per user)
- Results visualization
- Interactive documentation with Swagger UI

### Web Interface

- Modern interface with Bootstrap 5
- Registration and login
- Candidate visualization
- Intuitive voting system
- Real-time results with charts

## Installation and Execution

### Prerequisites

- Go 1.16 or higher

### Installation Steps

1. **Clone or create the project**

```bash
mkdir voting-system
cd voting-system
```

2. **Save the server code** in `main.go`

3. **Initialize the Go module and install dependencies**

```bash
go mod init voting-system
go get golang.org/x/crypto/bcrypt
```

4. **Run the server**

```bash
go run main.go
```

The server will start at `http://127.0.0.1:8000`

### Run the web interface

1. **Save the HTML file** as `index.html`

2. **Open in the browser** (any of these options):

   - Directly: Double-click on `index.html`
   - With Python server:
     ```bash
     python -m http.server 3000
     # Open http://localhost:3000
     ```
   - With Node.js:
     ```bash
     npx http-server -p 3000
     # Open http://localhost:3000
     ```

## API Documentation

### Swagger UI

Once the server is running, access:

- **Swagger UI:** http://127.0.0.1:8000/swagger/
- **OpenAPI Specification:** http://127.0.0.1:8000/api/swagger.json

## API Endpoints

### Authentication

#### Register User

```http
POST /api/register/
Content-Type: application/json

{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "securepass123"
}
```

**Successful response (201):**

```json
{
  "id": 1,
  "username": "john_doe",
  "email": "john@example.com",
  "has_voted": false
}
```

#### Log In

```http
POST /api/login/
Content-Type: application/json

{
  "username": "john_doe",
  "password": "securepass123"
}
```

**Successful response (200):**

```json
{
  "token": "a1b2c3d4e5f6...",
  "user": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "has_voted": false
  }
}
```

#### Log Out

```http
POST /api/logout/
Authorization: Token <your_token>
```

**Successful response (200):**

```json
{
  "message": "Successfully logged out"
}
```

### Voting

#### List Candidates

```http
GET /api/candidates/
Authorization: Token <your_token>
```

**Successful response (200):**

```json
[
  {
    "id": 1,
    "name": "Alice Johnson"
  },
  {
    "id": 2,
    "name": "Bob Smith"
  },
  {
    "id": 3,
    "name": "Charlie Brown"
  }
]
```

#### Cast Vote

```http
POST /api/vote/
Authorization: Token <your_token>
Content-Type: application/json

{
  "candidate": 1
}
```

**Successful response (201):**

```json
{
  "id": 1,
  "user_id": 1,
  "candidate_id": 1,
  "created_at": "2025-10-29T10:30:00Z"
}
```

**Error - Already voted (400):**

```json
{
  "error": "user has already voted"
}
```

#### View Results

```http
GET /api/results/
Authorization: Token <your_token>
```

**Successful response (200):**

```json
[
  {
    "id": 1,
    "user_id": 1,
    "candidate_id": 1,
    "created_at": "2025-10-29T10:30:00Z"
  },
  {
    "id": 2,
    "user_id": 2,
    "candidate_id": 2,
    "created_at": "2025-10-29T10:35:00Z"
  }
]
```

## Testing with cURL

### Full test flow

#### 1. Register a user

```bash
curl -X POST http://127.0.0.1:8000/api/register/ \
  -H "Content-Type: application/json" \
  -d
    "username": "alice",
    "email": "alice@example.com",
    "password": "password123"

```

#### 2. Log in

```bash
curl -X POST http://127.0.0.1:8000/api/login/ \
  -H "Content-Type: application/json" \
  -d
    "username": "alice",
    "password": "password123"

```

**Save the token you receive in the response**

#### 3. View candidates

```bash
curl -X GET http://127.0.0.1:8000/api/candidates/ \
  -H "Authorization: Token YOUR_TOKEN_HERE"
```

#### 4. Cast vote

```bash
curl -X POST http://127.0.0.1:8000/api/vote/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Token YOUR_TOKEN_HERE" \
  -d
    "candidate": 1

```

#### 5. View results

```bash
curl -X GET http://127.0.0.1:8000/api/results/ \
  -H "Authorization: Token YOUR_TOKEN_HERE"
```

#### 6. Log out

```bash
curl -X POST http://127.0.0.1:8000/api/logout/ \
  -H "Authorization: Token YOUR_TOKEN_HERE"
```

### Full test script

```bash
#!/bin/bash

API_URL="http://127.0.0.1:8000"

echo "1. Registering user..."
curl -s -X POST $API_URL/api/register/ \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"test123"}' | jq

echo -e "\n2. Logging in..."
RESPONSE=$(curl -s -X POST $API_URL/api/login/ \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"test123"}')

TOKEN=$(echo $RESPONSE | jq -r '.token')
echo "Token obtained: $TOKEN"

echo -e "\n3. Getting candidates..."
curl -s -X GET $API_URL/api/candidates/ \
  -H "Authorization: Token $TOKEN" | jq

echo -e "\n4. Voting for candidate 1..."
curl -s -X POST $API_URL/api/vote/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Token $TOKEN" \
  -d '{"candidate":1}' | jq

echo -e "\n5. Viewing results..."
curl -s -X GET $API_URL/api/results/ \
  -H "Authorization: Token $TOKEN" | jq

echo -e "\n6. Logging out..."
curl -s -X POST $API_URL/api/logout/ \
  -H "Authorization: Token $TOKEN" | jq
```

Save the script as `test.sh`, give it execution permissions, and run it:

```bash
chmod +x test.sh
./test.sh
```

**Note:** Requires `jq` to format the JSON. Install it with:

- Ubuntu/Debian: `sudo apt-get install jq`
- macOS: `brew install jq`
- Or remove `| jq` from the script

## Architecture

### Project Structure

```
voting-system/
├── main.go           # Go server code
├── index.html        # Web interface
├── go.mod            # Go dependencies
├── go.sum            # Dependency checksums
└── README.md         # This file
```

### Technologies Used

#### Backend

- **Go (Golang)** - Programming language
- **net/http** - Standard HTTP server
- **bcrypt** - Password encryption
- **sync** - Synchronization for thread-safety

#### Frontend

- **HTML5** - Structure
- **Bootstrap 5** - Styles and components
- **JavaScript (Vanilla)** - Interaction logic
- **Fetch API** - HTTP calls

### Security Features

- ✅ Passwords encrypted with bcrypt
- ✅ Token-based authentication
- ✅ Cryptographically secure random tokens (32 bytes)
- ✅ Token validation on protected endpoints
- ✅ One vote per user (backend validation)
- ✅ CORS enabled for development

## Database

The system uses **in-memory** storage with thread-safe data structures:

- `users` - Map of users by ID
- `candidates` - Map of candidates by ID (pre-populated)
- `votes` - Map of votes by ID
- `tokens` - Map of authentication tokens

**Note:** Data is lost when the server restarts. For production, consider integrating a database like PostgreSQL or MongoDB.

## Predefined Candidates

The system includes three default candidates:

1. **Alice Johnson** (ID: 1)
2. **Bob Smith** (ID: 2)
3. **Charlie Brown** (ID: 3)

## Current Limitations

- In-memory database (not persistent)
- No pagination in results
- No advanced email validation
- No password recovery
- No rate limiting
- No structured logs

## Future Improvements

- [ ] Integration with a real database (PostgreSQL, MySQL, MongoDB)
- [ ] Password recovery system
- [ ] Admin panel
- [ ] Advanced statistics and charts
- [ ] Results export (CSV, PDF)
- [ ] Rate limiting and attack protection
- [ ] Unit and integration tests
- [ ] Docker and docker-compose
- [ ] CI/CD pipeline
- [ ] Cloud deployment

## Contributions

Contributions are welcome. Please:

1. Fork the project
2. Create a branch for your feature (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is open source and available under the MIT License.

## Author

Developed as an educational project for a voting system with Go.

## Report Issues

If you find any bugs or have suggestions, please open an issue in the repository.

---

**Happy voting! **
