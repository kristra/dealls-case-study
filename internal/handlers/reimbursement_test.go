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

func AuthStubMiddlewareForReimbursement() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Set("role", "Employee")
		c.Next()
	}
}

func setupTestRouterForReimbursement() *gin.Engine {
	r := gin.Default()
	r.POST("/reimbursements", AuthStubMiddlewareForReimbursement(), handlers.SubmitReimbursement)

	return r
}

func setupTestDBForReimbursement() (*gorm.DB, func(), error) {
	d, cleanup, err := db.InitTestDB()
	if err != nil {
		log.Fatalf("failed to initialize test db: %v", err)
	}

	password, _ := utils.HashPassword("password")
	user := models.User{
		ID:       1,
		Username: "johndoe",
		Password: password,
		RoleID:   2,
	}
	d.Create(&user)

	return d, cleanup, err
}

func TestSubmitReimbursement_Success(t *testing.T) {
	r := setupTestRouterForReimbursement()

	_, cleanup, err := setupTestDBForReimbursement()
	if err != nil {
		t.Fatalf("Failed to set up test DB: %v", err)
	}
	defer cleanup()

	description := "Travel expenses"
	reqBody := dto.SubmitReimbursementRequest{
		Amount:      50000,
		Description: &description,
	}

	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/reimbursements", bytes.NewBuffer(body))

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")

	var resp dto.SuccessResponse[dto.SubmitReimbursementResponse]

	err1 := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Nil(t, err1)
	assert.Equal(t, reqBody.Amount, resp.Data.Amount)
	assert.Equal(t, *reqBody.Description, *resp.Data.Description)
}

func TestSubmitReimbursement_InvalidPayload(t *testing.T) {

	r := setupTestRouterForReimbursement()

	_, cleanup, err := setupTestDBForReimbursement()
	if err != nil {
		t.Fatalf("Failed to set up test DB: %v", err)
	}
	defer cleanup()

	body := []byte(`{invalid json}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/reimbursements", bytes.NewBuffer(body))

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}
