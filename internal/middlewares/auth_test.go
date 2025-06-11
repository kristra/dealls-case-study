package middlewares

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func generateTestToken(userID uint, role string, secret []byte) string {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(secret)
	return tokenString
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	token := generateTestToken(1, "Employee", jwtSecret)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AuthMiddleware())
	r.GET("/protected", func(c *gin.Context) {
		userID := c.MustGet("user_id").(uint)
		role := c.MustGet("role").(string)
		c.JSON(http.StatusOK, gin.H{"user_id": userID, "role": role})
	})

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), `"user_id":1`)
	assert.Contains(t, resp.Body.String(), `"role":"Employee"`)
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AuthMiddleware())
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "should not reach here"})
	})

	req, _ := http.NewRequest("GET", "/protected", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "Missing or invalid token")
}

func TestAdminOnlyMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Inject role manually for testing AdminOnly
	r.Use(func(c *gin.Context) {
		c.Set("role", "Admin")
		c.Next()
	})
	r.Use(AdminOnly())
	r.GET("/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "admin allowed"})
	})

	req, _ := http.NewRequest("GET", "/admin", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "admin allowed")
}

func TestAdminOnlyMiddleware_Denied(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Inject non-admin role
	r.Use(func(c *gin.Context) {
		c.Set("role", "Employee")
		c.Next()
	})
	r.Use(AdminOnly())
	r.GET("/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "should not be allowed"})
	})

	req, _ := http.NewRequest("GET", "/admin", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusForbidden, resp.Code)
	assert.Contains(t, resp.Body.String(), "Admin access required")
}
