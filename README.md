# Voting System Project

This project implements a simple voting system using Django, Django REST Framework, and Swagger UI.
It provides both a traditional web interface and a RESTful API for interacting with the voting system.

## Features

### Web Interface

- User Registration
- User Login/Logout
- View Candidates
- Cast a Vote (one vote per user)

### REST API

- User Registration
- User Login (with token authentication)
- User Logout
- List Candidates
- Cast a Vote
- View Vote Results
- Interactive API documentation with Swagger UI and ReDoc

## REST API Endpoints

The API endpoints are accessible under the `/api/` prefix.

### API Documentation

- **Swagger UI:** `http://127.0.0.1:8000/swagger/`
  - Interactive API documentation where you can test the endpoints.
- **ReDoc:** `http://127.0.0.1:8000/redoc/`
  - Alternative API documentation view.

### Authentication

Most API endpoints require token authentication. After logging in via `/api/login/`, you will receive a token. Include this token in the `Authorization` header of subsequent requests as `Token <your_auth_token>`.

### Endpoints Details

- **Register User:**

  - `POST /api/register/`
  - **Body (JSON):** `{"username": "your_username", "email": "your_email@example.com", "password": "your_password"}`
  - Creates a new user and returns user details.

- **Login User:**

  - `POST /api/login/`
  - **Body (JSON):** `{"username": "your_username", "password": "your_password"}`
  - Authenticates a user and returns an authentication token.

- **Logout User:**

  - `POST /api/logout/`
  - **Headers:** `Authorization: Token <your_auth_token>`
  - Invalidates the user's authentication token.

- **List Candidates:**

  - `GET /api/candidates/`
  - **Headers:** `Authorization: Token <your_auth_token>`
  - Returns a list of all available candidates.

- **Cast a Vote:**

  - `POST /api/vote/`
  - **Headers:** `Authorization: Token <your_auth_token>`
  - **Body (JSON):** `{"candidate": <candidate_id>}`
  - Casts a vote for the specified candidate. Each user can only vote once.

- **View Vote Results:**
  - `GET /api/results/`
  - **Headers:** `Authorization: Token <your_auth_token>`
  - Returns a list of all votes cast.
# Voting_example_go_vanillaJS
