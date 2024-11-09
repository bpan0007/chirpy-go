package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	passw := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(passw, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

}
