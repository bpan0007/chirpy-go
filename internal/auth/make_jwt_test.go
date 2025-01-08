package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "mysecret"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	parsedToken, err := jwt.ParseWithClaims(token, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	claims, ok := parsedToken.Claims.(*MyCustomClaims)
	if !ok || !parsedToken.Valid {
		t.Fatalf("expected valid token, got invalid token")
	}

	if claims.Subject != userID.String() {
		t.Errorf("expected subject %v, got %v", userID.String(), claims.Subject)
	}

	if claims.Issuer != "chirpy" {
		t.Errorf("expected issuer chirpy, got %v", claims.Issuer)
	}

}
