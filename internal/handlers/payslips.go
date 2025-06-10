package handlers

import (
	"net/http"
	"strconv"

	"dealls-case-study/internal/db"
	"dealls-case-study/internal/models"

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
// @Success      200    {object}  models.Payslip
// @Failure      400    {object}  map[string]string
// @Failure      401    {object}  map[string]string
// @Failure      404    {object}  map[string]string
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

	c.JSON(http.StatusOK, payslip)
}
