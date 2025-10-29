package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Models
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`
	HasVoted bool   `json:"has_voted"`
}

type Candidate struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Vote struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	CandidateID int       `json:"candidate_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type Token struct {
	UserID    int
	Token     string
	CreatedAt time.Time
}

// Request/Response structs
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type VoteRequest struct {
	CandidateID int `json:"candidate"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// Database (in-memory)
type Database struct {
	users      map[int]*User
	candidates map[int]*Candidate
	votes      map[int]*Vote
	tokens     map[string]*Token
	userIDSeq  int
	voteIDSeq  int
	mu         sync.RWMutex
}

func NewDatabase() *Database {
	db := &Database{
		users:      make(map[int]*User),
		candidates: make(map[int]*Candidate),
		votes:      make(map[int]*Vote),
		tokens:     make(map[string]*Token),
		userIDSeq:  0,
		voteIDSeq:  0,
	}

	// Seed candidates
	db.candidates[1] = &Candidate{ID: 1, Name: "Alice Johnson"}
	db.candidates[2] = &Candidate{ID: 2, Name: "Bob Smith"}
	db.candidates[3] = &Candidate{ID: 3, Name: "Charlie Brown"}

	return db
}

func (db *Database) CreateUser(username, email, password string) (*User, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Check if username exists
	for _, u := range db.users {
		if u.Username == username {
			return nil, fmt.Errorf("username already exists")
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	db.userIDSeq++
	user := &User{
		ID:       db.userIDSeq,
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		HasVoted: false,
	}
	db.users[user.ID] = user
	return user, nil
}

func (db *Database) AuthenticateUser(username, password string) (*User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	for _, u := range db.users {
		if u.Username == username {
			err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
			if err != nil {
				return nil, fmt.Errorf("invalid credentials")
			}
			return u, nil
		}
	}
	return nil, fmt.Errorf("invalid credentials")
}

func (db *Database) CreateToken(userID int) (string, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	tokenStr := hex.EncodeToString(tokenBytes)

	token := &Token{
		UserID:    userID,
		Token:     tokenStr,
		CreatedAt: time.Now(),
	}
	db.tokens[tokenStr] = token
	return tokenStr, nil
}

func (db *Database) ValidateToken(tokenStr string) (*User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	token, exists := db.tokens[tokenStr]
	if !exists {
		return nil, fmt.Errorf("invalid token")
	}

	user, exists := db.users[token.UserID]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (db *Database) DeleteToken(tokenStr string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.tokens, tokenStr)
	return nil
}

func (db *Database) GetCandidates() []*Candidate {
	db.mu.RLock()
	defer db.mu.RUnlock()

	candidates := make([]*Candidate, 0, len(db.candidates))
	for _, c := range db.candidates {
		candidates = append(candidates, c)
	}
	return candidates
}

func (db *Database) CastVote(userID, candidateID int) (*Vote, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	user, exists := db.users[userID]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	if user.HasVoted {
		return nil, fmt.Errorf("user has already voted")
	}

	_, exists = db.candidates[candidateID]
	if !exists {
		return nil, fmt.Errorf("candidate not found")
	}

	db.voteIDSeq++
	vote := &Vote{
		ID:          db.voteIDSeq,
		UserID:      userID,
		CandidateID: candidateID,
		CreatedAt:   time.Now(),
	}
	db.votes[vote.ID] = vote
	user.HasVoted = true

	return vote, nil
}

func (db *Database) GetResults() []*Vote {
	db.mu.RLock()
	defer db.mu.RUnlock()

	votes := make([]*Vote, 0, len(db.votes))
	for _, v := range db.votes {
		votes = append(votes, v)
	}
	return votes
}

// HTTP Handlers
type Server struct {
	db *Database
}

func NewServer() *Server {
	return &Server{
		db: NewDatabase(),
	}
}

func (s *Server) extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Token" {
		return ""
	}
	return parts[1]
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, ErrorResponse{Error: message})
}

func (s *Server) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "All fields are required")
		return
	}

	user, err := s.db.CreateUser(req.Username, req.Email, req.Password)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, user)
}

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := s.db.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	token, err := s.db.CreateToken(user.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create token")
		return
	}

	respondJSON(w, http.StatusOK, LoginResponse{Token: token, User: *user})
}

func (s *Server) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	token := s.extractToken(r)
	if token == "" {
		respondError(w, http.StatusUnauthorized, "No token provided")
		return
	}

	s.db.DeleteToken(token)
	respondJSON(w, http.StatusOK, map[string]string{"message": "Successfully logged out"})
}

func (s *Server) CandidatesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	token := s.extractToken(r)
	if token == "" {
		respondError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	_, err := s.db.ValidateToken(token)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	candidates := s.db.GetCandidates()
	respondJSON(w, http.StatusOK, candidates)
}

func (s *Server) VoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	token := s.extractToken(r)
	if token == "" {
		respondError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	user, err := s.db.ValidateToken(token)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	var req VoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	vote, err := s.db.CastVote(user.ID, req.CandidateID)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, vote)
}

func (s *Server) ResultsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	token := s.extractToken(r)
	if token == "" {
		respondError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	_, err := s.db.ValidateToken(token)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	votes := s.db.GetResults()
	respondJSON(w, http.StatusOK, votes)
}

func (s *Server) SwaggerHandler(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Voting System API</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui.css">
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui-bundle.js"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                url: '/api/swagger.json',
                dom_id: '#swagger-ui',
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIBundle.SwaggerUIStandalonePreset
                ]
            });
        };
    </script>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func (s *Server) SwaggerJSONHandler(w http.ResponseWriter, r *http.Request) {
	spec := map[string]interface{}{
		"openapi": "3.0.0",
		"info": map[string]interface{}{
			"title":       "Voting System API",
			"description": "API for a simple voting system",
			"version":     "1.0.0",
		},
		"servers": []map[string]string{
			{"url": "http://127.0.0.1:8000"},
		},
		"components": map[string]interface{}{
			"securitySchemes": map[string]interface{}{
				"TokenAuth": map[string]string{
					"type": "apiKey",
					"in":   "header",
					"name": "Authorization",
				},
			},
		},
		"paths": map[string]interface{}{
			"/api/register/": map[string]interface{}{
				"post": map[string]interface{}{
					"summary": "Register a new user",
					"requestBody": map[string]interface{}{
						"required": true,
						"content": map[string]interface{}{
							"application/json": map[string]interface{}{
								"schema": map[string]interface{}{
									"type": "object",
									"properties": map[string]interface{}{
										"username": map[string]string{"type": "string"},
										"email":    map[string]string{"type": "string"},
										"password": map[string]string{"type": "string"},
									},
								},
							},
						},
					},
					"responses": map[string]interface{}{
						"201": map[string]interface{}{"description": "User created"},
					},
				},
			},
			"/api/login/": map[string]interface{}{
				"post": map[string]interface{}{
					"summary": "Login user",
					"requestBody": map[string]interface{}{
						"required": true,
						"content": map[string]interface{}{
							"application/json": map[string]interface{}{
								"schema": map[string]interface{}{
									"type": "object",
									"properties": map[string]interface{}{
										"username": map[string]string{"type": "string"},
										"password": map[string]string{"type": "string"},
									},
								},
							},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{"description": "Login successful"},
					},
				},
			},
			"/api/logout/": map[string]interface{}{
				"post": map[string]interface{}{
					"summary":  "Logout user",
					"security": []map[string][]string{{"TokenAuth": {}}},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{"description": "Logout successful"},
					},
				},
			},
			"/api/candidates/": map[string]interface{}{
				"get": map[string]interface{}{
					"summary":  "List all candidates",
					"security": []map[string][]string{{"TokenAuth": {}}},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{"description": "List of candidates"},
					},
				},
			},
			"/api/vote/": map[string]interface{}{
				"post": map[string]interface{}{
					"summary":  "Cast a vote",
					"security": []map[string][]string{{"TokenAuth": {}}},
					"requestBody": map[string]interface{}{
						"required": true,
						"content": map[string]interface{}{
							"application/json": map[string]interface{}{
								"schema": map[string]interface{}{
									"type": "object",
									"properties": map[string]interface{}{
										"candidate": map[string]string{"type": "integer"},
									},
								},
							},
						},
					},
					"responses": map[string]interface{}{
						"201": map[string]interface{}{"description": "Vote cast successfully"},
					},
				},
			},
			"/api/results/": map[string]interface{}{
				"get": map[string]interface{}{
					"summary":  "View vote results",
					"security": []map[string][]string{{"TokenAuth": {}}},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{"description": "List of votes"},
					},
				},
			},
		},
	}
	respondJSON(w, http.StatusOK, spec)
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func main() {
	server := NewServer()

	http.HandleFunc("/api/register/", corsMiddleware(server.RegisterHandler))
	http.HandleFunc("/api/login/", corsMiddleware(server.LoginHandler))
	http.HandleFunc("/api/logout/", corsMiddleware(server.LogoutHandler))
	http.HandleFunc("/api/candidates/", corsMiddleware(server.CandidatesHandler))
	http.HandleFunc("/api/vote/", corsMiddleware(server.VoteHandler))
	http.HandleFunc("/api/results/", corsMiddleware(server.ResultsHandler))
	http.HandleFunc("/swagger/", server.SwaggerHandler)
	http.HandleFunc("/api/swagger.json", server.SwaggerJSONHandler)

	port := 8000
	log.Printf("Server starting on port %d", port)
	log.Printf("Swagger UI available at: http://127.0.0.1:%d/swagger/", port)
	log.Printf("API endpoints available at: http://127.0.0.1:%d/api/", port)

	if err := http.ListenAndServe(":"+strconv.Itoa(port), nil); err != nil {
		log.Fatal(err)
	}
}