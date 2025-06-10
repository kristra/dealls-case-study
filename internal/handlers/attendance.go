package handlers

import (
	"dealls-case-study/internal/db"
	"dealls-case-study/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CheckInAttendance godoc
// @Summary      Create attendance
// @Description  Create attendance record and set check-in based on now value
// @Tags         Attendance
// @Accept       json
// @Produce      json
// @Success      200    {object}  models.Attendance
// @Failure      400    {object}  map[string]string
// @Failure      401    {object}  map[string]string
// @Failure      500    {object}  map[string]string
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

	c.JSON(http.StatusOK, newAttendance)
}

// CheckOutAttendance godoc
// @Summary      Update attendance
// @Description  Update attendance record and set check-out based on now value
// @Tags         Attendance
// @Accept       json
// @Produce      json
// @Success      200    {object}  models.Attendance
// @Failure      400    {object}  map[string]string
// @Failure      401    {object}  map[string]string
// @Failure      500    {object}  map[string]string
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

	c.JSON(http.StatusOK, attendance)
}
