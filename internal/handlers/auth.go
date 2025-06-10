package handlers

import (
	"net/http"

	"dealls-case-study/internal/db"
	"dealls-case-study/internal/dto"
	"dealls-case-study/internal/models"
	"dealls-case-study/internal/utils"

	"github.com/gin-gonic/gin"
)

// Login godoc
// @Summary      User Login
// @Description  Auntheticates a user using username and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body     dto.LoginRequest true "Login credentials"
// @Success      200    {object}  dto.SuccessResponse[dto.LoginResponse]
// @Failure      400    {object}  dto.ErrorResponse
// @Failure      401    {object}  dto.ErrorResponse
// @Failure      500    {object}  dto.ErrorResponse
// @Router       /auth/login [post]
func Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := db.DB.Preload("Role").First(&user, "users.username = ?", req.Username).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username"})
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Role.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, utils.WrapSuccessResponse(dto.LoginResponse{
		Token: token,
		User: dto.LoginUser{
			ID:       user.ID,
			Username: user.Username,
			Role:     user.Role.Name,
		},
	}))
}
