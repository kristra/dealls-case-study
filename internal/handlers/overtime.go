package handlers

import (
	"net/http"
	"time"

	"dealls-case-study/internal/db"
	"dealls-case-study/internal/dto"
	"dealls-case-study/internal/models"
	"dealls-case-study/internal/utils"

	"github.com/gin-gonic/gin"
)

// SubmitOvertime godoc
// @Summary      Submit Overtime for current user
// @Description  Allows an employee to submit overtime for the current day.
// @Description  Overtime must be submitted after check-out, and cannot exceed 3 hours per day.
// @Tags         Attendance
// @Accept       json
// @Produce      json
// @Param        request body     dto.SubmitOvertimeRequest true "Overtime payloads"
// @Success      200    {object}  dto.SuccessResponse[dto.SubmitOvertimeResponse]
// @Failure      400    {object}  dto.ErrorResponse
// @Failure      401    {object}  dto.ErrorResponse
// @Failure      500    {object}  dto.ErrorResponse
// @Security     BearerAuth
// @Router       /attendances/overtime [post]
func SubmitOvertime(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.SubmitOvertimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	today := time.Now()

	var attendance models.Attendance
	if err := db.DB.
		Where("user_id = ? AND DATE(date) = ?", userID, today).
		First(&attendance).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "attendance record not found"})
		return
	}
	if attendance.CheckOutAt == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you must check out before submitting overtime"})
		return
	}

	var existing models.Overtime
	if err := db.DB.
		Where("user_id = ? AND DATE(date) = ?", userID, today).
		First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "overtime already submitted for today"})
		return
	}

	overtime := models.Overtime{
		UserID:      userID,
		Date:        today,
		HoursWorked: req.HoursWorked,
		CreatedBy:   userID,
	}
	if err := db.DB.Create(&overtime).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to submit overtime"})
		return
	}

	c.JSON(http.StatusOK, utils.WrapSuccessResponse(dto.SubmitOvertimeResponse{
		ID:          overtime.ID,
		HoursWorked: overtime.HoursWorked,
	}))
}
