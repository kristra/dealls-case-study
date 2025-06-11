package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword_ReturnsHashedValue(t *testing.T) {
	password := "securepassword123"
	hash, err := HashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)
}

func TestCheckPasswordHash_MatchesCorrectPassword(t *testing.T) {
	password := "securepassword123"
	hash, _ := HashPassword(password)

	match := CheckPasswordHash(password, hash)
	assert.True(t, match, "expected password to match hash")
}

func TestCheckPasswordHash_FailsWithIncorrectPassword(t *testing.T) {
	password := "securepassword123"
	wrongPassword := "wrongpassword"
	hash, _ := HashPassword(password)

	match := CheckPasswordHash(wrongPassword, hash)
	assert.False(t, match, "expected password mismatch to fail")
}
