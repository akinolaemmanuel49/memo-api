package helpers

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

//HashPassword is used to encrypt the password before it is stored in the repository.
func HashPassword(password *string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(*password), 14)
	if err != nil {
		return err
	}

	*password = string(hash)
	return nil
}

//VerifyPassword checks the inputted plaintext password against the password hash in the repository.
func VerifyPassword(hashedPassword, plainPassword string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
