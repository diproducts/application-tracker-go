package password_hasher_test

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/diproducts/application-tracker-go/internal/lib/auth/password_hasher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBcryptPasswordHasher(t *testing.T) {
	h := password_hasher.NewBcryptPasswordHasher()
	password := generatePassword()
	invalidPassword := "invalid_password"

	hashedPassword, err := h.Generate(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	err = h.Compare(hashedPassword, invalidPassword)
	assert.Error(t, err)

	err = h.Compare(hashedPassword, password)
	assert.NoError(t, err)
}

func generatePassword() string {
	return gofakeit.Password(true, true, true, true, false, 8)
}
