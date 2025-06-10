package handlers

import (
	"net/http"
	"strconv"
	"time"

	"dealls-case-study/internal/db"
	"dealls-case-study/internal/dto"
	"dealls-case-study/internal/models"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// UpsertPayroll godoc
// @Summary      Upsert payroll
// @Description  Create or update payroll record for a given month and year
// @Tags         Payroll
// @Accept       json
// @Produce      json
// @Param        year   path      int  true  "Year"
// @Param        month  path      int  true  "Month"
// @Param        request body     dto.UpsertPayrollRequest false "Payroll fields"
// @Success      200    {object}  models.Payroll
// @Failure      400    {object}  map[string]string
// @Failure      401    {object}  map[string]string
// @Failure      500    {object}  map[string]string
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
			Month:       month,
			Year:        year,
			GeneratedAt: time.Now(),
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

	status := "updated"
	if result.RowsAffected == 0 {
		status = "created"
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payroll " + status, "data": payroll})
}
