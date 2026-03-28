package handler

import (
	"errors"
	"net/http"
	"project-management/internal/usecase"
	"project-management/internal/utils"

	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	TaskService usecase.TaskService
}

func NewTaskHandler(t usecase.TaskService) *TaskHandler {
	return &TaskHandler{t}
}

type CreateTaskInput struct {
	Title       string     `json:"title" binding:"required,min=3,max=100"`
	Description string     `json:"description" binding:"max=500"`
	ProjectID   uint       `json:"project_id" binding:"required"`
	AssigneeID  uint       `json:"assigned_to"`
	Status      string     `json:"status" binding:"omitempty,oneof=todo in_progress done"`
	DueDate     *time.Time `json:"due_date"`
}

func (t *TaskHandler) CreateTask(c *gin.Context) {
	var input CreateTaskInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.FormatValidationError(err))
		return
	}

	task, err := t.TaskService.CreateTask(input.Title, input.Description, input.ProjectID, input.AssigneeID, input.Status, input.DueDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create task")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, task)
}

func (t *TaskHandler) ListTasks(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid page")
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid limit")
		return
	}

	status := c.Query("status")
	projectIDStr := c.Query("project_id")

	var projectID uint
	if projectIDStr != "" {
		id, err := strconv.Atoi(projectIDStr)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid project_id")
			return
		}
		projectID = uint(id)
	}

	assignedToStr := c.Query("assigned_to")
	var assignedTo uint

	if assignedToStr != "" {
		id, err := strconv.Atoi(assignedToStr)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid assigned_to")
			return
		}
		assignedTo = uint(id)
	}
	tasks, total, err := t.TaskService.ListTasks(page, limit, status, projectID, assignedTo)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch tasks")
		return
	}

	totalPages := (int(total) + limit - 1) / limit

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"data": tasks,
		"meta": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

func (t *TaskHandler) GetTaskByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid task ID")
		return
	}

	task, err := t.TaskService.GetTaskByID(uint(id))
	if errors.Is(err, usecase.ErrTaskNotFound) {
		utils.ErrorResponse(c, http.StatusNotFound, "Task not found")
		return
	}
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch task")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, task)

}

type UpdateTaskInput struct {
	Title       string `json:"title" binding:"omitempty,min=3,max=100"`
	Description string `json:"description" binding:"omitempty,max=500"`
	Status      string `json:"status" binding:"omitempty,oneof=pending in-progress done"`
	AssigneeID  uint   `json:"assignee_id"`
}

func (t *TaskHandler) UpdateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid task ID")
		return
	}

	var input UpdateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.FormatValidationError(err))
		return
	}

	task, err := t.TaskService.UpdateTask(
		uint(id),
		input.Title,
		input.Description,
		input.Status,
		input.AssigneeID,
	)

	if errors.Is(err, usecase.ErrTaskNotFound) {
		utils.ErrorResponse(c, http.StatusNotFound, "Task not found")
		return
	}
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update task")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, task)
}

func (t *TaskHandler) DeleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid task ID")
		return
	}

	err = t.TaskService.DeleteTask(uint(id))

	if errors.Is(err, usecase.ErrTaskNotFound) {
		utils.ErrorResponse(c, http.StatusNotFound, "Task not found")
		return
	}
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete task")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Task deleted successfully")
}

func (t *TaskHandler) UpdateTaskStatus(c *gin.Context) {
	idParam := c.Param("id")
	taskID64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	taskID := uint(taskID64)

	var req struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists || userIDVal == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	userID, ok := userIDVal.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id"})
		return
	}

	roleVal, exists := c.Get("role")
	if !exists || roleVal == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "role not found"})
		return
	}

	role, ok := roleVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid role"})
		return
	}

	err = t.TaskService.UpdateTaskStatus(taskID, req.Status, userID, role)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task status updated successfully"})
}
