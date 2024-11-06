package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type CreateUserRequest struct {
	Email string `json:"email"`
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
	user, err := cfg.dbQueries.CreateUser(r.Context(), req.Email)
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
