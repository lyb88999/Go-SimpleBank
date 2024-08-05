package util

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password using the bcrypt algorithm
func HashPassword(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashPassword), nil
}

// CheckPassword checks if the provided password is correct or not
func CheckPassword(password, hashPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
}
