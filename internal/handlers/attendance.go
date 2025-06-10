package handlers

import (
	"dealls-case-study/internal/db"
	"dealls-case-study/internal/dto"
	"dealls-case-study/internal/models"
	"dealls-case-study/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CheckInAttendance godoc
// @Summary      Submit check-in for current user
// @Description  Allows an employee to check in for the current day.
// @Description  Only one check-in is allowed per day. Check-ins on weekends are not allowed.
// @Tags         Attendance
// @Accept       json
// @Produce      json
// @Success      200    {object}  dto.SuccessResponse[dto.AttendanceResponse]
// @Failure      400    {object}  dto.ErrorResponse
// @Failure      401    {object}  dto.ErrorResponse
// @Failure      500    {object}  dto.ErrorResponse
// @Security     BearerAuth
// @Router       /attendances/check-in [post]
func CheckInAttendance(c *gin.Context) {
	userID := c.GetUint("user_id")
	now := time.Now()

	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot check in on weekends"})
		return
	}

	var attendance models.Attendance
	dateStr := now.Format("2006-01-02")
	tx := db.DB.Where("user_id = ? AND DATE(date) = ?", userID, dateStr).First(&attendance)

	if tx.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Already checked in today"})
		return
	}

	newAttendance := models.Attendance{
		UserID:    userID,
		Date:      now,
		CheckInAt: &now,
		CreatedBy: userID,
	}

	if err := db.DB.Create(&newAttendance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check in"})
		return
	}

	c.JSON(http.StatusOK, utils.WrapSuccessResponse(dto.AttendanceResponse{
		ID:         newAttendance.ID,
		Date:       newAttendance.Date,
		CheckInAt:  newAttendance.CheckInAt,
		CheckOutAt: newAttendance.CheckOutAt,
	}))
}

// CheckOutAttendance godoc
// @Summary      Submit check-out for current user
// @Description  Allows an employee to check out for the current day.
// @Description  Must have checked in first. Only one check-out is allowed per day. Check-outs on weekends are not allowed.
// @Tags         Attendance
// @Accept       json
// @Produce      json
// @Success      200    {object}  dto.SuccessResponse[dto.AttendanceResponse]
// @Failure      400    {object}  dto.ErrorResponse
// @Failure      401    {object}  dto.ErrorResponse
// @Failure      500    {object}  dto.ErrorResponse
// @Security     BearerAuth
// @Router       /attendances/check-out [post]
func CheckOutAttendance(c *gin.Context) {
	userID := c.GetUint("user_id")
	now := time.Now()

	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot check out on weekends"})
		return
	}

	var attendance models.Attendance
	dateStr := now.Format("2006-01-02")
	tx := db.DB.Where("user_id = ? AND DATE(date) = ?", userID, dateStr).First(&attendance)

	if tx.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You have not checked in today"})
		return
	}

	if attendance.CheckOutAt != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Already checked out today"})
		return
	}

	attendance.CheckOutAt = &now
	attendance.UpdatedBy = userID

	if err := db.DB.Save(&attendance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check out"})
		return
	}

	c.JSON(http.StatusOK, utils.WrapSuccessResponse(dto.AttendanceResponse{
		ID:         attendance.ID,
		Date:       attendance.Date,
		CheckInAt:  attendance.CheckInAt,
		CheckOutAt: attendance.CheckOutAt,
	}))
}
