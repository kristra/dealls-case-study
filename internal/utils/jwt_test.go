package utils

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	testSecret := "testsecret"
	_ = os.Setenv("JWT_SECRET", testSecret)
	jwtSecret = []byte(testSecret)

	userID := uint(1)
	role := "Employee"

	tokenString, err := GenerateToken(userID, role)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Parse the token back
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	assert.NoError(t, err)
	assert.True(t, token.Valid)

	// Assert claims
	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, float64(userID), claims["user_id"])
	assert.Equal(t, role, claims["role"])

	// Check if expiration exists and is in the future
	exp, ok := claims["exp"].(float64)
	assert.True(t, ok)
	assert.Greater(t, int64(exp), time.Now().Unix())
}
