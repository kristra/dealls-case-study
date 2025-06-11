package handlers_test

import (
	"bytes"
	"dealls-case-study/internal/db"
	"dealls-case-study/internal/dto"
	"dealls-case-study/internal/handlers"
	"dealls-case-study/internal/models"
	"dealls-case-study/internal/utils"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestDBforAuth() (*gorm.DB, func(), error) {

	d, cleanup, err := db.InitTestDB()

	password, err := utils.HashPassword("password")
	if err != nil {
		log.Fatalf("Failed creating admin password: %v", err)
	}
	// use this if hash fn fails for whatever reason
	// password := "$2a$14$TkqaFJJGv1jP1r.GByPbE.0smTliuY53ccS96ggXkroZbt0WsKG/S"
	user := models.User{
		ID:       1,
		Username: "johndoe",
		Password: password,
		RoleID:   2,
	}
	d.Create(&user)

	return d, cleanup, err
}

func setupTestRouterforAuth() *gin.Engine {
	r := gin.Default()
	r.POST("/auth/login", func(c *gin.Context) {
		handlers.Login(c)
	})

	return r
}

func TestLogin_CorrectCredentials(t *testing.T) {
	r := setupTestRouterforAuth()

	_, cleanup, err := setupTestDBforAuth()
	if err != nil {
		t.Fatalf("failed to set up test DB: %v", err)
	}
	defer cleanup()

	body := dto.LoginRequest{Username: "johndoe", Password: "password"}
	out, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(out))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")

	var usr dto.SuccessResponse[dto.LoginResponse]
	err1 := json.Unmarshal(w.Body.Bytes(), &usr)
	assert.Nil(t, err1)
	token, _ := utils.GenerateToken(uint(1), "Employee")
	assert.Equal(t, usr.Data, dto.LoginResponse{
		Token: token,
		User: dto.LoginUser{
			ID:       1,
			Username: "johndoe",
			Role:     "Employee",
		},
	})
}

func TestLogin_WrongPassword(t *testing.T) {
	r := setupTestRouterforAuth()

	_, cleanup, err := setupTestDBforAuth()
	if err != nil {
		t.Fatalf("failed to set up test DB: %v", err)
	}
	defer cleanup()

	body := dto.LoginRequest{Username: "johndoe", Password: "password1"}
	out, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(out))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid password")
}

func TestLogin_WrongUsername(t *testing.T) {
	r := setupTestRouterforAuth()

	_, cleanup, err := setupTestDBforAuth()
	if err != nil {
		t.Fatalf("failed to set up test DB: %v", err)
	}
	defer cleanup()

	body := dto.LoginRequest{Username: "johndoe1", Password: "password"}
	out, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(out))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid username")
}
