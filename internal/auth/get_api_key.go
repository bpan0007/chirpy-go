package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Add a func GetAPIKey(headers http.Header) (string, error) to  auth package. It should extract the api key from the Authorization header, which is expected to be in this format Authorization: ApiKey THE_KEY_HERE:
// Authorization
func GetAPIKey(headers http.Header) (string, error) {
	// Add your code here

	authHeaderValue := headers.Get("Authorization")
	if authHeaderValue == "" {
		return "", errors.New("No API key provided")
	}

	// Expected format: "ApiKey THE_KEY_HERE"
	const prefix = "ApiKey "
	if !strings.HasPrefix(authHeaderValue, prefix) {
		return "", errors.New("Invalid API key format")
	}

	cleaned_api_Key := strings.TrimPrefix(authHeaderValue, prefix)

	return cleaned_api_Key, nil
}
