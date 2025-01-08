package main

import (
	"encoding/json"
	"net/http"

	"github.com/bpan0007/chirpy-go/internal/auth"
	"github.com/bpan0007/chirpy-go/internal/database"
)

type UpdateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	// Validate the access token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	// Decode the request body
	var updateReq UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	// Validate input
	if updateReq.Email == "" || updateReq.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Email and password are required", nil)
		return
	}

	// Hash the new password
	hashedPassword, err := auth.HashPassword(updateReq.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
		return
	}

	// Update the user in the database
	updatedUser, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userID,
		Email:          updateReq.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error updating user", err)
		return
	}

	// Respond with the updated user (excluding the password)
	responseUser := User{
		ID:        updatedUser.ID,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
		Email:     updatedUser.Email,
	}
	respondWithJSON(w, http.StatusOK, responseUser)
}
