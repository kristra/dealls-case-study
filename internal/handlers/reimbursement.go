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

// SubmitReimbursement godoc
// @Summary      Submit reimbursement for current user
// @Description  Allows an employee to submit a reimbursement request.
// @Tags         Reimbursements
// @Accept       json
// @Produce      json
// @Param        request body     dto.SubmitReimbursementRequest true "Reimbursement data"
// @Success      201    {object}  dto.SuccessResponse[dto.SubmitReimbursementResponse]
// @Failure      400    {object}  dto.ErrorResponse
// @Failure      401    {object}  dto.ErrorResponse
// @Failure      500    {object}  dto.ErrorResponse
// @Security     BearerAuth
// @Router       /reimbursements [post]
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

	c.JSON(http.StatusOK, utils.WrapSuccessResponse(dto.SubmitReimbursementResponse{
		ID:          reimbursement.ID,
		Amount:      reimbursement.Amount,
		Description: &reimbursement.Description,
	}))
}
