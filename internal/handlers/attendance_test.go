package handlers_test

import (
	"dealls-case-study/internal/db"
	"dealls-case-study/internal/dto"
	"dealls-case-study/internal/handlers"
	"dealls-case-study/internal/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestDBforAtt() (*gorm.DB, func(), error) {
	d, cleanup, err := db.InitTestDB()

	d.Exec("DELETE FROM attendances")

	user := models.User{
		ID:       1,
		Username: "johndoe",
		Password: "password",
		RoleID:   2,
	}
	d.Create(&user)

	return d, cleanup, err
}

func AuthStubMiddlewareForAtt() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Set("role", "Employee")
		c.Next()
	}
}

func setupTestRouterforAtt() *gin.Engine {
	r := gin.Default()
	r.POST("/attendances/check-in", AuthStubMiddlewareForAtt(), handlers.CheckInAttendance)
	r.POST("/attendances/check-out", AuthStubMiddlewareForAtt(), handlers.CheckOutAttendance)
	return r
}

func TestCheckInAttendance(t *testing.T) {
	r := setupTestRouterforAtt()

	_, cleanup, err := setupTestDBforAtt()
	if err != nil {
		t.Fatalf("failed to set up test DB: %v", err)
	}
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/attendances/check-in", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")

	var att dto.SuccessResponse[dto.AttendanceResponse]

	err1 := json.Unmarshal(w.Body.Bytes(), &att)
	assert.Nil(t, err1)
	assert.Equal(t, att.Data.ID, uint(1))
	assert.Equal(t, att.Data.Date.Day(), time.Now().Day())
	assert.Equal(t, att.Data.Date.Month(), time.Now().Month())
	assert.Equal(t, att.Data.Date.Year(), time.Now().Year())
	assert.NotNil(t, att.Data.CheckInAt)
}

func TestCheckOutAttendance(t *testing.T) {
	r := setupTestRouterforAtt()
	now := time.Now()

	d, cleanup, err := setupTestDBforAtt()
	if err != nil {
		t.Fatalf("failed to set up test DB: %v", err)
	}
	checkInTime := now.Add(-9 * time.Hour)
	att := models.Attendance{
		UserID:    1,
		Date:      now,
		CheckInAt: &checkInTime,
		CreatedBy: 1,
	}
	d.Create(&att)
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/attendances/check-out", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")

	var updated dto.SuccessResponse[dto.AttendanceResponse]

	err1 := json.Unmarshal(w.Body.Bytes(), &updated)
	assert.Nil(t, err1)
	assert.NotNil(t, updated.Data.CheckOutAt)
}
