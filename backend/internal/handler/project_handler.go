package handler

import (
	"errors"
	"net/http"
	"project-management/internal/usecase"
	"project-management/internal/utils"

	"strconv"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	ProjectService usecase.ProjectService
}

func NewProjectHandler(p usecase.ProjectService) *ProjectHandler {
	return &ProjectHandler{p}
}

type CreateProjectInput struct {
	Name        string `json:"name" binding:"required,min=3"`
	Description string `json:"description" binding:"max=500"`
}

type UpdateProjectInput struct {
	Name        string `json:"name" binding:"required,min=3"`
	Description string `json:"description" binding:"max=500"`
}

func (p *ProjectHandler) CreateProject(c *gin.Context) {
	var input CreateProjectInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.FormatValidationError(err))
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, ok := userIDVal.(uint)
	if !ok {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Invalid user context")
		return
	}

	project, err := p.ProjectService.CreateProject(input.Name, input.Description, userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create project")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, project)
}

func (p *ProjectHandler) ListProjects(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	projects, total, err := p.ProjectService.ListProjects(page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch projects")
		return
	}

	if limit == 0 {
		limit = 10
	}

	totalPages := (int(total) + limit - 1) / limit

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"data": projects,
		"meta": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

func (p *ProjectHandler) GetProjectByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid project ID")
		return
	}

	project, err := p.ProjectService.GetProjectByID(uint(id))

	if errors.Is(err, usecase.ErrProjectNotFound) {
		utils.ErrorResponse(c, http.StatusNotFound, "Project not found")
		return
	}
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch project")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, project)
}

func (p *ProjectHandler) UpdateProject(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid project ID")
		return
	}

	var input UpdateProjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.FormatValidationError(err))
		return
	}

	project, err := p.ProjectService.UpdateProject(uint(id), input.Name, input.Description)
	if errors.Is(err, usecase.ErrProjectNotFound) {
		utils.ErrorResponse(c, http.StatusNotFound, "Project not found")
		return
	}
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update project")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, project)
}

func (p *ProjectHandler) DeleteProject(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid project ID")
		return
	}

	err = p.ProjectService.DeleteProject(uint(id))

	if errors.Is(err, usecase.ErrProjectNotFound) {
		utils.ErrorResponse(c, http.StatusNotFound, "Project not found")
		return
	}
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete project")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Project deleted successfully")
}
