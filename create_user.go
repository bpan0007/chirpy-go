package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/bpan0007/chirpy-go/internal/auth"
	"github.com/bpan0007/chirpy-go/internal/database"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Password  string    `json:"hashed_password"`
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) createUsers(w http.ResponseWriter, r *http.Request) {

	var req CreateUserRequest

	// Decode the JSON request body into the req struct
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	log.Printf("Decoded request body with email: %s", req.Email)

	// Try to create the user
	log.Printf("Attempting to create user in database")
	// Create user in database
	hashedPassword, _ := auth.HashPassword(req.Password)
	params := database.CreateUserParams{
		Email:          req.Email,
		HashedPassword: hashedPassword,
	}
	user, err := cfg.db.CreateUser(r.Context(), params)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	log.Printf("Successfully created user with ID: %v", user.ID)
	const maxChirpLength = 140

	if len(req.Email) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	returnUser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJSON(w, http.StatusCreated, returnUser)
}

func (cfg *apiConfig) login(w http.ResponseWriter, r *http.Request) {

	var loginReq CreateUserRequest

	// Decode the JSON request body into the req struct
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	log.Printf("Decoded request body with email: %s", loginReq.Email)

	// Hash the passw
	log.Printf("Attempting to Hash the passw")

	user, err := cfg.db.GetUserByEmail(r.Context(), loginReq.Email)
	if err != nil {
		// If user not found or other database error
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	log.Printf("loginReq.Password: %s", loginReq.Password)
	log.Printf(" user.HashedPassword: %s", user.HashedPassword)
	if err := auth.CheckPasswordHash(loginReq.Password, user.HashedPassword); err != nil {
		// Passwords don't match
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	// Passwords match! Return the user
	returnUser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	respondWithJSON(w, http.StatusOK, returnUser)
}
