package main

import (
	"net/http"

	"github.com/bpan0007/chirpy-go/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	dbChirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}
	var chirpFound bool
	for _, dbChirp := range dbChirps {
		if dbChirp.ID == chirpID && dbChirp.UserID == userID {
			chirpFound = true
			break
		}
	}
	if !chirpFound {
		respondWithError(w, http.StatusForbidden, "You can't delete this chirp", err)
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
