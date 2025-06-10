package handlers

import (
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
		Where("user_id = ? AND year = ? AND month = ?", userID, year, month).
		First(&payslip).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payslip not found"})
		return
	}

	c.JSON(http.StatusOK, utils.WrapSuccessResponse(dto.PayslipResponse{
		ID:                     payslip.ID,
		Month:                  payslip.Month,
		Year:                   payslip.Year,
		UserID:                 payslip.UserID,
		BaseSalary:             payslip.BaseSalary,
		OvertimePay:            payslip.OvertimePay,
		Reimbursement:          payslip.Reimbursement,
		TotalSalary:            payslip.TotalSalary,
		TotalHoursWorked:       payslip.TotalHoursWorked,
		TotalOvertimeHours:     payslip.TotalOvertimeHours,
		AttendanceBreakdown:    payslip.AttendanceBreakdown,
		OvertimeBreakdown:      payslip.OvertimeBreakdown,
		ReimbursementBreakdown: payslip.ReimbursementBreakdown,
	}))
}
