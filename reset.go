package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Only for prod", nil)
	}
	log.Printf("Attempting to delete all users")

	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		log.Printf("Error deleting users: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete users", err)
		return
	}

	log.Printf("Successfully deleted all users")

}
