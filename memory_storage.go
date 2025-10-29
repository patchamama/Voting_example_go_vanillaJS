package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// MemoryStorage is an in-memory implementation of the Storage interface.
type MemoryStorage struct {
	users      map[int]*User
	candidates map[int]*Candidate
	votes      map[int]*Vote
	tokens     map[string]*Token
	userIDSeq  int
	voteIDSeq  int
	mu         sync.RWMutex
}

// NewMemoryStorage creates a new in-memory storage.
func NewMemoryStorage() *MemoryStorage {
	db := &MemoryStorage{
		users:      make(map[int]*User),
		candidates: make(map[int]*Candidate),
		votes:      make(map[int]*Vote),
		tokens:     make(map[string]*Token),
		userIDSeq:  0,
		voteIDSeq:  0,
	}

	return db
}

// Connect does nothing for in-memory storage.
func (db *MemoryStorage) Connect() error {
	return nil
}

// Migrate seeds the database with initial data.
func (db *MemoryStorage) Migrate() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Seed candidates
	db.candidates[1] = &Candidate{ID: 1, Name: "Alice Johnson"}
	db.candidates[2] = &Candidate{ID: 2, Name: "Bob Smith"}
	db.candidates[3] = &Candidate{ID: 3, Name: "Charlie Brown"}

	return nil
}

// CreateUser creates a new user.
func (db *MemoryStorage) CreateUser(username, email, password string) (*User, error) {
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

// AuthenticateUser authenticates a user.
func (db *MemoryStorage) AuthenticateUser(username, password string) (*User, error) {
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

// CreateToken creates a new token for a user.
func (db *MemoryStorage) CreateToken(userID int) (string, error) {
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

// ValidateToken validates a token.
func (db *MemoryStorage) ValidateToken(tokenStr string) (*User, error) {
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

// DeleteToken deletes a token.
func (db *MemoryStorage) DeleteToken(tokenStr string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.tokens, tokenStr)
	return nil
}

// GetCandidates returns all candidates.
func (db *MemoryStorage) GetCandidates() ([]*Candidate, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	candidates := make([]*Candidate, 0, len(db.candidates))
	for _, c := range db.candidates {
		candidates = append(candidates, c)
	}
	return candidates, nil
}

// CastVote casts a vote for a candidate.
func (db *MemoryStorage) CastVote(userID, candidateID int) (*Vote, error) {
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

// GetResults returns all votes.
func (db *MemoryStorage) GetResults() ([]*Vote, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	votes := make([]*Vote, 0, len(db.votes))
	for _, v := range db.votes {
		votes = append(votes, v)
	}
	return votes, nil
}
