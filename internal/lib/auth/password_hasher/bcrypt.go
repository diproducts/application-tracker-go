package password_hasher

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type BcryptPasswordHasher struct{}

func NewBcryptPasswordHasher() *BcryptPasswordHasher {
	return &BcryptPasswordHasher{}
}

func (h *BcryptPasswordHasher) Generate(password string) (string, error) {
	const op = "lib.auth.password_hasher.Generate"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return string(hashedPassword), nil
}

func (h *BcryptPasswordHasher) Compare(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
