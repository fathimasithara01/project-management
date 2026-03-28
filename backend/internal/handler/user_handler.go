package handler

import (
	"net/http"
	"project-management/internal/usecase"
	"project-management/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService usecase.UserService
}

func NewUserHandler(u usecase.UserService) *UserHandler {
	return &UserHandler{u}
}

func (u *UserHandler) CreateUser(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required,min=3"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Role     string `json:"role" binding:"omitempty,oneof=admin developer"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.FormatValidationError(err))
		return
	}

	user, err := u.UserService.Register(input.Name, input.Email, input.Password, input.Role)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to create user")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func (u *UserHandler) ListUsers(c *gin.Context) {
	users, err := u.UserService.ListUsers()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch users")
		return
	}

	result := []gin.H{}
	for _, user := range users {
		result = append(result, gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		})
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"users": result,
	})
}

func (u *UserHandler) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := u.UserService.GetUserByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to get user")
		return
	}

	if user == nil {
		utils.ErrorResponse(c, http.StatusNotFound, "user not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}
