package handlers

import (
	"net/http"
	"time"

	"dealls-case-study/internal/db"
	"dealls-case-study/internal/dto"
	"dealls-case-study/internal/models"

	"github.com/gin-gonic/gin"
)

func SubmitReimbursement(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.SubmitReimbursementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reimbursement := models.Reimbursement{
		UserID:    userID,
		Amount:    req.Amount,
		Date:      time.Now(),
		CreatedBy: userID,
	}

	if req.Description != nil {
		reimbursement.Description = *req.Description
	}

	if err := db.DB.Create(&reimbursement).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to submit reimbursement"})
		return
	}

	c.JSON(http.StatusOK, reimbursement)
}
