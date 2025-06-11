package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"dealls-case-study/internal/db"
	"dealls-case-study/internal/dto"
	"dealls-case-study/internal/models"
	"dealls-case-study/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// @BasePath /api/v1

// UpsertPayroll godoc
// @Summary      Upsert payroll
// @Description  Creates or updates a payroll record for the given month and year.
// @Description  Only updates payrolls with 'draft' status.
// @Tags         Payroll
// @Accept       json
// @Produce      json
// @Param        year   path      int  true  "Year"
// @Param        month  path      int  true  "Month (1-12)"
// @Param        request body     dto.UpsertPayrollRequest false "Payroll optional fields"
// @Success      200    {object}  dto.SuccessResponse[dto.PayrollResponse]
// @Failure      400    {object}  dto.ErrorResponse
// @Failure      401    {object}  dto.ErrorResponse
// @Failure      500    {object}  dto.ErrorResponse
// @Security     BearerAuth
// @Router       /payrolls/{year}/{month} [post]
func UpsertPayroll(c *gin.Context) {
	yearStr := c.Param("year")
	monthStr := c.Param("month")

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year"})
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month (must be 1–12)"})
		return
	}

	var req dto.UpsertPayrollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var payroll models.Payroll
	result := db.DB.Where("month = ? AND year = ?", month, year).First(&payroll)

	if result.Error != nil && result.RowsAffected == 0 {
		payroll = models.Payroll{
			Month: month,
			Year:  year,
		}
	}

	if req.Name != nil {
		payroll.Name = *req.Name
	}
	if req.PeriodStart != nil {
		payroll.PeriodStart = *req.PeriodStart
	}
	if req.PeriodEnd != nil {
		payroll.PeriodEnd = *req.PeriodEnd
	}

	if err := db.DB.Save(&payroll).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upsert payroll"})
		return
	}

	c.JSON(http.StatusOK, utils.WrapSuccessResponse(dto.PayrollResponse{
		ID:          payroll.ID,
		Name:        payroll.Name,
		PeriodStart: &payroll.PeriodStart,
		PeriodEnd:   &payroll.PeriodEnd,
		Status:      payroll.Status,
	}))
}

// RunPayroll godoc
// @Summary      Run payroll
// @Description  Processes the payroll for the given month and year.
// @Description  Generates payslips for all employees.
// @Description  Can only be run once per period. Once run, the payroll status changes to 'pending' or 'completed'.
// @Tags         Payroll
// @Accept       json
// @Produce      json
// @Param        year   path      int  true  "Year"
// @Param        month  path      int  true  "Month (1-12)"
// @Success      200    {object}  dto.SuccessResponse[dto.PayrollResponse]
// @Failure      400    {object}  dto.ErrorResponse
// @Failure      401    {object}  dto.ErrorResponse
// @Failure      500    {object}  dto.ErrorResponse
// @Security     BearerAuth
// @Router       /payrolls/{year}/{month}/run [post]
func RunPayroll(c *gin.Context) {
	year, err1 := strconv.Atoi(c.Param("year"))
	month, err2 := strconv.Atoi(c.Param("month"))
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year or month"})
		return
	}

	userID := c.GetUint("user_id")
	var payroll models.Payroll

	if err := db.DB.Where("year = ? AND month = ?", year, month).First(&payroll).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payroll record not found"})
		return
	}

	switch payroll.Status {
	case models.PayrollStatusProcessed:
		c.JSON(http.StatusBadRequest, gin.H{"error": "payroll has already been processed"})
		return
	case models.PayrollStatusPending:
		c.JSON(http.StatusBadRequest, gin.H{"error": "payroll is currently being processed"})
		return
	}

	payroll.Status = models.PayrollStatusPending
	payroll.ProcessedAt = time.Now()
	payroll.UpdatedBy = userID

	if err := db.DB.Save(&payroll).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to run payroll"})
		return
	}

	// we use goroutine here for simplicity sake, just to demonstrate the async logic
	// in real world application this should be processed using workers or some job queue solutions
	// for example the simplest implementation would be a separate worker that processes any `pending` payroll
	go ProcessPayroll(payroll.ID, userID)

	c.JSON(http.StatusOK, utils.WrapSuccessResponse(dto.PayrollResponse{
		ID:          payroll.ID,
		Name:        payroll.Name,
		PeriodStart: &payroll.PeriodStart,
		PeriodEnd:   &payroll.PeriodEnd,
		Status:      payroll.Status,
	}))
}

func ProcessPayroll(payrollID uint, adminID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var payroll models.Payroll

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", payrollID).
			First(&payroll).Error; err != nil {
			return err
		}

		if payroll.Status != models.PayrollStatusPending {
			return errors.New("payroll is not in a pending state")
		}

		var payslips []models.Payslip
		var users []models.User
		if err := tx.Find(&users).Error; err != nil {
			return err
		}
		for _, user := range users {
			payslip, err := GeneratePayslip(tx, adminID, user, &payroll)

			if err != nil {
				return err
			}

			payslips = append(payslips, payslip)
		}
		if len(payslips) > 0 {
			if err := tx.Create(payslips).Error; err != nil {
				return err
			}
		}

		payroll.Status = models.PayrollStatusProcessed
		if err := tx.Save(&payroll).Error; err != nil {
			return err
		}

		log.Printf("Payroll %d processed successfully", payroll.ID)
		return nil
	})
}

func GeneratePayslip(tx *gorm.DB, adminID uint, user models.User, payroll *models.Payroll) (models.Payslip, error) {
	var attendances []models.Attendance
	tx.Where("user_id = ? AND date BETWEEN ? AND ?", user.ID, payroll.PeriodStart, payroll.PeriodEnd).Find(&attendances)

	var overtimes []models.Overtime
	tx.Where("user_id = ? AND date BETWEEN ? AND ?", user.ID, payroll.PeriodStart, payroll.PeriodEnd).Find(&overtimes)

	var reimbursements []models.Reimbursement
	tx.Where("user_id = ? AND date BETWEEN ? AND ?", user.ID, payroll.PeriodStart, payroll.PeriodEnd).Find(&reimbursements)

	daysWorked := len(attendances)
	// flat 8 hours per days worked instead of using HoursWorked field
	// based on this requirements:
	// No rules for late or early check-ins or check-outs; check-in at any time that day counts.
	totalHours := float64(daysWorked * 8)
	totalOvertime := 0.0
	for _, o := range overtimes {
		totalOvertime += o.HoursWorked
	}

	totalReimbursement := 0.0
	for _, r := range reimbursements {
		totalReimbursement += r.Amount
	}

	expectedWorkingDays := utils.CountWeekdays(payroll.PeriodStart, payroll.PeriodEnd)
	// flat 8 hours per days worked
	hourlyRate := user.Salary / float64(expectedWorkingDays) / 8
	overtimeRatePerHour := hourlyRate * 2
	proratedDaysWorked := float64(daysWorked) / float64(expectedWorkingDays)

	basePay := proratedDaysWorked * user.Salary
	overtimePay := totalOvertime * hourlyRate * 2
	totalPay := basePay + overtimePay + totalReimbursement

	payslip := models.Payslip{
		Month:     payroll.Month,
		Year:      payroll.Year,
		UserID:    user.ID,
		PayrollID: payroll.ID,

		// summary totals
		BaseSalary:    basePay,
		OvertimePay:   overtimePay,
		Reimbursement: totalReimbursement,
		TotalSalary:   totalPay,

		// calculation context
		MonthlySalary:       user.Salary,
		ExpectedWorkingDays: expectedWorkingDays,
		DaysAttended:        daysWorked,
		HourlyRate:          hourlyRate,
		OvertimeRatePerHour: overtimeRatePerHour,

		// breakdowns
		TotalHoursWorked:       totalHours,
		TotalOvertimeHours:     totalOvertime,
		AttendanceBreakdown:    toJSON(attendances),
		OvertimeBreakdown:      toJSON(overtimes),
		ReimbursementBreakdown: toJSON(reimbursements),

		CreatedBy: adminID,
	}

	return payslip, nil
}

// GetPayrollSummary godoc
// @Summary      Get payroll summary
// @Description  Generates a summary of all employee payslips for a given month and year.
// @Tags         Payroll
// @Security     BearerAuth
// @Produce      json
// @Param        year   path      int  true  "Year"
// @Param        month  path      int  true  "Month"
// @Success      200    {object}  dto.SuccessResponse[dto.PayrollSummaryResponse]
// @Failure      400    {object}  dto.ErrorResponse
// @Failure      401    {object}  dto.ErrorResponse
// @Failure      403    {object}  dto.ErrorResponse
// @Failure      500    {object}  dto.ErrorResponse
// @Router       /payrolls/{year}/{month}/summary [get]
func GeneratePayrollSummary(c *gin.Context) {
	year, err1 := strconv.Atoi(c.Param("year"))
	month, err2 := strconv.Atoi(c.Param("month"))
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year or month"})
		return
	}

	var payroll models.Payroll

	if err := db.DB.Where("year = ? AND month = ?", year, month).First(&payroll).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payroll record not found"})
		return
	}

	if payroll.Status != models.PayrollStatusProcessed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payroll has not been processed"})
		return
	}

	var payslips []models.Payslip
	if err := db.DB.
		Preload("User").
		Where("payroll_id = ?", payroll.ID).
		Find(&payslips).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch payslips"})
		return
	}

	var summary dto.PayrollSummaryResponse
	summary.PayrollID = payroll.ID
	summary.Year = payroll.Year
	summary.Month = payroll.Month
	summary.Payslips = make([]dto.EmployeePayslipBrief, 0)

	for _, p := range payslips {
		summary.TotalSalaries += p.TotalSalary
		summary.Payslips = append(summary.Payslips, dto.EmployeePayslipBrief{
			UserID:        p.UserID,
			Username:      p.User.Username,
			BaseSalary:    p.BaseSalary,
			OvertimePay:   p.OvertimePay,
			Reimbursement: p.Reimbursement,
			TotalPay:      p.TotalSalary,
		})
	}

	c.JSON(http.StatusOK, utils.WrapSuccessResponse(summary))
}

func toJSON[T any](v T) string {
	b, err := json.Marshal(v)
	if err != nil {
		log.Printf("error marshaling: %v", err)
		return "[]"
	}
	return string(b)
}
