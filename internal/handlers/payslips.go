package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"dealls-case-study/internal/db"
	"dealls-case-study/internal/dto"
	"dealls-case-study/internal/models"
	"dealls-case-study/internal/utils"

	"github.com/gin-gonic/gin"
)

// GetPayslip godoc
// @Summary      Get payslip for current user
// @Description  Fetches payslip for a specific month and year
// @Tags         Payslip
// @Security     BearerAuth
// @Produce      json
// @Param        year   path      int  true  "Year"
// @Param        month  path      int  true  "Month"
// @Success      200    {object}  dto.SuccessResponse[dto.PayslipResponse]
// @Failure      400    {object}  dto.ErrorResponse
// @Failure      401    {object}  dto.ErrorResponse
// @Failure      404    {object}  dto.ErrorResponse
// @Router       /payslips/{year}/{month} [get]
func GetPayslip(c *gin.Context) {
	userID := c.GetUint("user_id")

	yearStr := c.Param("year")
	monthStr := c.Param("month")
	year, err1 := strconv.Atoi(yearStr)
	month, err2 := strconv.Atoi(monthStr)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year or month"})
		return
	}

	var payslip models.Payslip
	err := db.DB.
		Preload("User").
		Where("user_id = ? AND year = ? AND month = ?", userID, year, month).
		First(&payslip).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payslip not found"})
		return
	}

	var aB []dto.AttendanceBreakdownItem
	err = json.Unmarshal([]byte(payslip.AttendanceBreakdown), &aB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var oB []dto.OvertimeBreakdownItem
	err = json.Unmarshal([]byte(payslip.OvertimeBreakdown), &oB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var rB []dto.ReimbursementBreakdownItem
	err = json.Unmarshal([]byte(payslip.ReimbursementBreakdown), &rB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, utils.WrapSuccessResponse(dto.PayslipResponse{
		ID:     payslip.ID,
		Month:  payslip.Month,
		Year:   payslip.Year,
		UserID: payslip.UserID,

		// summary totals
		BaseSalary:    payslip.BaseSalary,
		OvertimePay:   payslip.OvertimePay,
		Reimbursement: payslip.Reimbursement,
		TotalSalary:   payslip.TotalSalary,

		// calculation context
		MonthlySalary:       payslip.User.Salary,
		ExpectedWorkingDays: payslip.ExpectedWorkingDays,
		DaysAttended:        payslip.DaysAttended,
		HourlyRate:          payslip.HourlyRate,
		OvertimeRatePerHour: payslip.OvertimeRatePerHour,

		// breakdowns
		TotalHoursWorked:       payslip.TotalHoursWorked,
		TotalOvertimeHours:     payslip.TotalOvertimeHours,
		AttendanceBreakdown:    aB,
		OvertimeBreakdown:      oB,
		ReimbursementBreakdown: rB,
	}))
}
