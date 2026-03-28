package usecase

import (
	"errors"
	"project-management/internal/domain"
	"project-management/internal/repositories"

	"time"

	"gorm.io/gorm"
)

var ErrTaskNotFound = errors.New("task not found")

type TaskService interface {
	CreateTask(title, description string, projectID, assigneeID uint, status string, dueDate *time.Time) (*domain.Task, error)
	ListTasks(page, limit int, status string, projectID, assignedTo uint) ([]domain.Task, int64, error)
	GetTaskByID(id uint) (*domain.Task, error)
	UpdateTask(id uint, title, description, status string, assigneeID uint) (*domain.Task, error)
	DeleteTask(id uint) error
	UpdateTaskStatus(taskID uint, status string, userID uint, role string) error
}

type taskService struct {
	repo repositories.TaskRepository
}

func NewTaskService(r repositories.TaskRepository) TaskService {
	return &taskService{r}
}

func (t *taskService) CreateTask(title, description string, projectID, assigneeID uint, status string, dueDate *time.Time) (*domain.Task, error) {
	if status == "" {
		status = "todo"
	}

	tx := t.repo.BeginTx()
	repo := t.repo.WithTx(tx)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var project domain.Project
	if err := tx.First(&project, projectID).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("invalid project")
	}

	var assignedTo *uint
	if assigneeID != 0 {
		var user domain.User
		if err := tx.First(&user, assigneeID).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("invalid assignee")
		}
		assignedTo = &assigneeID
	}

	task := &domain.Task{
		Title:       title,
		Description: description,
		ProjectID:   projectID,
		AssignedTo:  assignedTo,
		Status:      domain.Status(status),
		DueDate:     dueDate,
	}

	if err := repo.Create(task); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return task, nil
}

func (t *taskService) ListTasks(page, limit int, status string, projectID, assignedTo uint) ([]domain.Task, int64, error) {

	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	tasks, err := t.repo.List(limit, offset, status, projectID, assignedTo)
	if err != nil {
		return nil, 0, err
	}

	total, err := t.repo.Count(status, projectID, assignedTo)
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (t *taskService) GetTaskByID(id uint) (*domain.Task, error) {
	task, err := t.repo.GetByID(id)

	if task == nil {
		return nil, ErrTaskNotFound
	}
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (t *taskService) UpdateTask(id uint, title, description, status string, assigneeID uint) (*domain.Task, error) {

	tx := t.repo.BeginTx()
	repo := t.repo.WithTx(tx)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	task, err := repo.GetByID(id)
	if task == nil {
		tx.Rollback()
		return nil, ErrTaskNotFound
	}
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if title != "" {
		task.Title = title
	}
	if description != "" {
		task.Description = description
	}

	if status != "" {
		task.Status = domain.Status(status)
	}
	if assigneeID != 0 {
		task.AssignedTo = &assigneeID
	}

	if err := repo.Update(task); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return task, nil
}

func (t *taskService) DeleteTask(id uint) error {

	tx := t.repo.BeginTx()
	repo := t.repo.WithTx(tx)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := repo.Delete(id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return ErrTaskNotFound
	}

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (t *taskService) UpdateTaskStatus(taskID uint, status string, userID uint, role string) error {
	task, err := t.repo.GetByID(taskID)
	if err != nil {
		return err
	}

	if task == nil {
		return errors.New("task not found")
	}

	if role == "developer" && task.AssignedTo != &userID {
		return errors.New("not allowed to update this task")
	}

	task.Status = domain.Status(status)

	return t.repo.Update(task)
}
