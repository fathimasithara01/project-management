package repositories

import (
	"errors"
	"project-management/internal/domain"

	"gorm.io/gorm"
)

type TaskRepository interface {
	Create(task *domain.Task) error
	List(limit, offset int, status string, projectID, assignedTo uint) ([]domain.Task, error)
	Count(status string, projectID uint, assignedTo uint) (int64, error)
	GetByID(id uint) (*domain.Task, error)
	Update(task *domain.Task) error
	Delete(id uint) error

	BeginTx() *gorm.DB
	WithTx(tx *gorm.DB) TaskRepository
}

type taskRepo struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepo{db}
}

func (r *taskRepo) BeginTx() *gorm.DB {
	return r.db.Begin()
}

func (r *taskRepo) WithTx(tx *gorm.DB) TaskRepository {
	return &taskRepo{db: tx}
}

func (r *taskRepo) Create(task *domain.Task) error {
	if err := r.db.Create(task).Error; err != nil {
		return err
	}

	return r.db.
		Preload("Project").
		Preload("Assignee").
		First(task, task.ID).Error
}

func (r *taskRepo) List(limit, offset int, status string, projectID, assignedTo uint) ([]domain.Task, error) {
	var tasks []domain.Task

	query := r.db.Preload("Project").Preload("Assignee")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if projectID != 0 {
		query = query.Where("project_id = ?", projectID)
	}

	if assignedTo != 0 {
		query = query.Where("assigned_to = ?", assignedTo)
	}

	err := query.
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&tasks).Error

	return tasks, err
}

func (r *taskRepo) Count(status string, projectID uint, assignedTo uint) (int64, error) {
	var count int64

	query := r.db.Model(&domain.Task{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if projectID != 0 {
		query = query.Where("project_id = ?", projectID)
	}

	if assignedTo != 0 {
		query = query.Where("assigned_to = ?", assignedTo)
	}

	err := query.Count(&count).Error

	return count, err
}

func (r *taskRepo) GetByID(id uint) (*domain.Task, error) {
	var task domain.Task

	err := r.db.Preload("Project").Preload("Assignee").First(&task, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *taskRepo) Update(task *domain.Task) error {
	return r.db.Model(&domain.Task{}).
		Where("id = ?", task.ID).
		Updates(map[string]interface{}{
			"title":       task.Title,
			"description": task.Description,
			"status":      task.Status,
			"assigned_to": task.AssignedTo,
		}).Error
}

func (r *taskRepo) Delete(id uint) error {
	result := r.db.Delete(&domain.Task{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
