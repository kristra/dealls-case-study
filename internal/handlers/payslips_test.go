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

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func AuthStubMiddlewareForPayslip() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Set("role", "Employee")
		c.Next()
	}
}

func setupTestRouterForPayslip() *gin.Engine {
	r := gin.Default()
	r.GET("/payslips/:year/:month", AuthStubMiddlewareForPayslip(), handlers.GetPayslip)
	return r
}

func setupTestDBForPayslip() (func(), error) {
	d, cleanup, err := db.InitTestDB()
	if err != nil {
		return nil, err
	}

	user := models.User{
		ID:       1,
		Username: "johndoe",
		Password: "password",
		RoleID:   2,
	}

	if err := d.Create(&user).Error; err != nil {
		return nil, err
	}

	payroll := models.Payroll{
		ID:     1,
		Month:  5,
		Year:   2024,
		Status: "Processed",
	}
	if err := d.Create(&payroll).Error; err != nil {
		return nil, err
	}

	payslip := models.Payslip{
		PayrollID:              1,
		UserID:                 1,
		Month:                  5,
		Year:                   2024,
		BaseSalary:             5000,
		OvertimePay:            200,
		Reimbursement:          100,
		TotalSalary:            5300,
		TotalHoursWorked:       160,
		TotalOvertimeHours:     10,
		AttendanceBreakdown:    "Full attendance",
		OvertimeBreakdown:      "10 hours in total",
		ReimbursementBreakdown: "Travel expenses",
	}
	if err := d.Create(&payslip).Error; err != nil {
		return nil, err
	}

	return cleanup, nil
}

func TestGetPayslip_Success(t *testing.T) {
	r := setupTestRouterForPayslip()

	cleanup, err := setupTestDBForPayslip()
	if err != nil {
		t.Fatalf("Failed to set up test DB: %v", err)
	}
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/payslips/2024/5", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.SuccessResponse[dto.PayslipResponse]
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), response.Data.UserID)
	assert.Equal(t, 2024, response.Data.Year)
	assert.Equal(t, 5, response.Data.Month)
}

func TestGetPayslip_NotFound(t *testing.T) {
	r := setupTestRouterForPayslip()

	// No payslip seeded
	_, cleanup, err := db.InitTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test DB: %v", err)
	}
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/payslips/2024/12", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Payslip not found")
}

func TestGetPayslip_InvalidParams(t *testing.T) {
	r := setupTestRouterForPayslip()

	// No payslip seeded
	_, cleanup, err := db.InitTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test DB: %v", err)
	}
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/payslips/abc/xyz", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid year or month")
}
