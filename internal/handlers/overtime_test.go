package handlers_test

import (
	"bytes"
	"dealls-case-study/internal/db"
	"dealls-case-study/internal/dto"
	"dealls-case-study/internal/handlers"
	"dealls-case-study/internal/models"
	"dealls-case-study/internal/utils"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func AuthStubMiddlewareForOvertime() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Set("role", "Employee")
		c.Next()
	}
}

func setupTestRouterForOvertime() *gin.Engine {
	r := gin.Default()
	r.POST("/attendances/overtime", AuthStubMiddlewareForOvertime(), handlers.SubmitOvertime)
	return r
}

func setupDBWithAttendanceAndCheckOut() (*models.User, func(), error) {
	d, cleanup, err := db.InitTestDB()
	if err != nil {
		return nil, nil, err
	}

	password, _ := utils.HashPassword("password")
	user := models.User{
		ID:       1,
		Username: "employee1",
		Password: password,
		RoleID:   2,
	}
	d.Create(&user)

	checkIn := time.Now().Add(-10 * time.Hour)
	checkOut := time.Now().Add(-2 * time.Hour)
	attendance := models.Attendance{
		UserID:     user.ID,
		Date:       time.Now(),
		CheckInAt:  &checkIn,
		CheckOutAt: &checkOut,
		CreatedBy:  user.ID,
	}
	d.Create(&attendance)

	return &user, cleanup, nil
}

func setupDBWithCheckInOnly() (*models.User, func(), error) {
	d, cleanup, err := db.InitTestDB()
	if err != nil {
		return nil, nil, err
	}

	password, _ := utils.HashPassword("password")
	user := models.User{
		ID:       2,
		Username: "nocheckoutuser",
		Password: password,
		RoleID:   2,
	}
	d.Create(&user)

	checkIn := time.Now().Add(-6 * time.Hour)
	attendance := models.Attendance{
		UserID:    user.ID,
		Date:      time.Now(),
		CheckInAt: &checkIn,
		CreatedBy: user.ID,
	}
	d.Create(&attendance)

	return &user, cleanup, nil
}

func TestSubmitOvertime_Success(t *testing.T) {
	r := setupTestRouterForOvertime()
	_, cleanup, err := setupDBWithAttendanceAndCheckOut()
	if err != nil {
		t.Fatalf("failed setup: %v", err)
	}
	defer cleanup()

	payload := dto.SubmitOvertimeRequest{HoursWorked: 2}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/attendances/overtime", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer dummy")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.SuccessResponse[dto.SubmitOvertimeResponse]
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Nil(t, err)
	assert.Equal(t, 2.0, resp.Data.HoursWorked)
}

func TestSubmitOvertime_AttendanceNotFound(t *testing.T) {
	r := setupTestRouterForOvertime()
	_, cleanup, err := db.InitTestDB()
	if err != nil {
		t.Fatalf("failed setup: %v", err)
	}
	defer cleanup()

	payload := dto.SubmitOvertimeRequest{HoursWorked: 2}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/attendances/overtime", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer dummy")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "attendance record not found")
}

func TestSubmitOvertime_MissingCheckOut(t *testing.T) {
	r := gin.Default()
	r.POST("/attendances/overtime", func(c *gin.Context) {
		c.Set("user_id", uint(2))
		c.Set("role", "Employee")
		c.Next()
	}, handlers.SubmitOvertime)

	_, cleanup, err := setupDBWithCheckInOnly()
	if err != nil {
		t.Fatalf("failed setup: %v", err)
	}
	defer cleanup()

	payload := dto.SubmitOvertimeRequest{HoursWorked: 1.5}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/attendances/overtime", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer dummy")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "you must check out before submitting overtime")
}
