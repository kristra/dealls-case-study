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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func AuthStubMiddlewareForPayroll() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", uint(1)) // Admin
		c.Set("role", "Admin")
		c.Next()
	}
}

func setupTestRouterForPayroll() *gin.Engine {
	r := gin.Default()
	r.POST("/payrolls/:year/:month", AuthStubMiddlewareForPayroll(), handlers.UpsertPayroll)
	r.POST("/payrolls/:year/:month/run", AuthStubMiddlewareForPayroll(), handlers.RunPayroll)
	return r
}

func setupTestDBForPayroll() (*gorm.DB, func(), error) {
	d, cleanup, err := db.InitTestDB()
	if err != nil {
		log.Fatalf("failed to initialize test db: %v", err)
	}

	password, _ := utils.HashPassword("password")
	admin := models.User{
		ID:       1,
		Username: "adminuser",
		Password: password,
		RoleID:   1, // admin
		Salary:   2200000,
	}
	employee := models.User{
		ID:       2,
		Username: "employee",
		Password: password,
		RoleID:   2, // employee
		Salary:   2200000,
	}
	d.Create(&admin)
	d.Create(&employee)

	return d, cleanup, err
}

func TestUpsertPayroll_Success(t *testing.T) {
	r := setupTestRouterForPayroll()
	_, cleanup, err := setupTestDBForPayroll()
	if err != nil {
		t.Fatalf("Failed to set up test DB: %v", err)
	}
	defer cleanup()

	name := "June Payroll"
	start := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC)

	reqBody := dto.UpsertPayrollRequest{
		Name:        &name,
		PeriodStart: &start,
		PeriodEnd:   &end,
	}

	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/payrolls/2025/6", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.SuccessResponse[dto.PayrollResponse]
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Nil(t, err)
	assert.Equal(t, name, resp.Data.Name)
	assert.Equal(t, "draft", resp.Data.Status)
}

func TestRunPayroll_Success(t *testing.T) {
	r := setupTestRouterForPayroll()
	d, cleanup, err := setupTestDBForPayroll()
	if err != nil {
		t.Fatalf("Failed to set up test DB: %v", err)
	}
	defer cleanup()

	start := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC)

	// Insert a payroll record
	payroll := models.Payroll{
		Month:       6,
		Year:        2025,
		Name:        "June Payroll",
		PeriodStart: start,
		PeriodEnd:   end,
	}
	d.Create(&payroll)

	// Insert attendance, overtime, and reimbursement data for employee
	attendance := models.Attendance{
		UserID:     2,
		Date:       time.Date(2025, 6, 5, 0, 0, 0, 0, time.UTC),
		CheckInAt:  timePtr(time.Date(2025, 6, 5, 9, 0, 0, 0, time.UTC)),
		CheckOutAt: timePtr(time.Date(2025, 6, 5, 17, 0, 0, 0, time.UTC)),
	}
	overtime := models.Overtime{
		UserID:      2,
		Date:        time.Date(2025, 6, 5, 0, 0, 0, 0, time.UTC),
		HoursWorked: 2,
		CreatedBy:   1,
	}
	reimbursement := models.Reimbursement{
		UserID:    2,
		Amount:    100000,
		Date:      time.Date(2025, 6, 5, 0, 0, 0, 0, time.UTC),
		CreatedBy: 1,
	}
	d.Create(&attendance)
	d.Create(&overtime)
	d.Create(&reimbursement)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/payrolls/2025/6/run", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.SuccessResponse[dto.PayrollResponse]
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Nil(t, err)
	assert.Equal(t, "pending", resp.Data.Status)
}

func TestUpsertPayroll_InvalidYear(t *testing.T) {
	r := setupTestRouterForPayroll()
	_, cleanup, err := setupTestDBForPayroll()
	if err != nil {
		t.Fatalf("DB setup failed: %v", err)
	}
	defer cleanup()

	body := []byte(`{}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/payrolls/0/6", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid year")
}

func TestUpsertPayroll_InvalidMonth(t *testing.T) {
	r := setupTestRouterForPayroll()
	_, cleanup, err := setupTestDBForPayroll()
	if err != nil {
		t.Fatalf("DB setup failed: %v", err)
	}
	defer cleanup()

	body := []byte(`{}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/payrolls/2025/13", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid month")
}

func TestUpsertPayroll_InvalidPayload(t *testing.T) {
	r := setupTestRouterForPayroll()
	_, cleanup, err := setupTestDBForPayroll()
	if err != nil {
		t.Fatalf("DB setup failed: %v", err)
	}
	defer cleanup()

	body := []byte(`{invalid json}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/payrolls/2025/6", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestRunPayroll_InvalidParams(t *testing.T) {
	r := setupTestRouterForPayroll()
	_, cleanup, err := setupTestDBForPayroll()
	if err != nil {
		t.Fatalf("DB setup failed: %v", err)
	}
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/payrolls/abc/xyz/run", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid year or month")
}

func TestRunPayroll_NotFound(t *testing.T) {
	r := setupTestRouterForPayroll()
	_, cleanup, err := setupTestDBForPayroll()
	if err != nil {
		t.Fatalf("DB setup failed: %v", err)
	}
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/payrolls/2025/12/run", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "payroll record not found")
}

func TestRunPayroll_AlreadyProcessed(t *testing.T) {
	r := setupTestRouterForPayroll()
	d, cleanup, err := setupTestDBForPayroll()
	if err != nil {
		t.Fatalf("DB setup failed: %v", err)
	}
	defer cleanup()

	payroll := models.Payroll{
		Month:  6,
		Year:   2025,
		Status: models.PayrollStatusProcessed,
	}
	d.Create(&payroll)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/payrolls/2025/6/run", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "already been processed")
}

func TestRunPayroll_PendingState(t *testing.T) {
	r := setupTestRouterForPayroll()
	d, cleanup, err := setupTestDBForPayroll()
	if err != nil {
		t.Fatalf("DB setup failed: %v", err)
	}
	defer cleanup()

	payroll := models.Payroll{
		Month:  6,
		Year:   2025,
		Status: models.PayrollStatusPending,
	}
	d.Create(&payroll)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/payrolls/2025/6/run", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "currently being processed")
}

func timePtr(t time.Time) *time.Time {
	return &t
}
