package handler

import (
	"net/http"
	"project-management/internal/usecase"
	"project-management/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	UserService usecase.UserService
}

func NewAuthHandler(u usecase.UserService) *AuthHandler {
	return &AuthHandler{
		UserService: u,
	}
}

type RegisterInput struct {
	Name     string `json:"name" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"omitempty,oneof=admin developer"`
}

func (ac *AuthHandler) Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.FormatValidationError(err))
		return
	}

	user, err := ac.UserService.Register(input.Name, input.Email, input.Password, input.Role)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to register")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, gin.H{
		"message": "user registered successfully",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (ac *AuthHandler) Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.FormatValidationError(err))
		return
	}

	token, user, err := ac.UserService.Login(input.Email, input.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "invalid credentials")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"message": "login successful",
		"token":   token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}
