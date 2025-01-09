package main

import (
	"encoding/json"
	"net/http"

	"github.com/bpan0007/chirpy-go/internal/auth"
	"github.com/google/uuid"
)

type polkaWebhookBody struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) polka_key(w http.ResponseWriter, r *http.Request) {

	// Add your code here

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key", err)
		return
	}

	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key", nil)
		return
	}

	decoder := json.NewDecoder(r.Body)

	var webhook polkaWebhookBody
	if err := decoder.Decode(&webhook); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	userUUID, err := uuid.Parse(webhook.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}
	_, err = cfg.db.UpgradeUserToChirpyRed(r.Context(), userUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error upgrading user", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("API key is valid"))
}
