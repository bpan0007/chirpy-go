package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type returnVals struct {
	CleanedBody string `json:"cleaned_body"`
}

type Request struct {
	Body string `json:"body"`
}

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	var req Request

	// Decode the JSON request body into the req struct
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// Use responseText to store the modified body text
	responseText := req.Body

	// Check and replace each restricted word and its capitalized version
	restrictedWords := []string{"kerfuffle", "sharbert", "fornax", "Kerfuffle", "Sharbert", "Fornax"}

	for _, word := range restrictedWords {
		// Replace lowercase version
		if strings.Contains(responseText, word) {
			responseText = strings.ReplaceAll(responseText, word, "****")
		}
		// Replace capitalized version
		// capitalizedWord := strings.ToLower(word) // Capitalizes the first letter
		// if strings.Contains(responseText, capitalizedWord) {
		// 	responseText = strings.ReplaceAll(responseText, capitalizedWord, "****")
		// }
	}

	// Check if the body length exceeds the maximum allowed length
	const maxChirpLength = 140
	if len(responseText) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	// Construct the response with the modified text
	returnVals := returnVals{
		CleanedBody: responseText,
	}

	respondWithJSON(w, http.StatusOK, returnVals)
}

// Mocked helper function for sending error responses
// func respondWithError(w http.ResponseWriter, code int, message string, err error) {
// 	http.Error(w, message, code)
// }

// // Mocked helper function for sending JSON responses
// func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(code)
// 	json.NewEncoder(w).Encode(payload)
// }
