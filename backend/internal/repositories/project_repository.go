package repositories

import (
	"errors"
	"project-management/internal/domain"

	"gorm.io/gorm"
)

var ErrProjectNotFound = errors.New("project not found")

type ProjectRepository interface {
	Create(project *domain.Project) error
	List(limit, offset int) ([]domain.Project, error)
	Count() (int64, error)
	GetByID(id uint) (*domain.Project, error)
	Update(project *domain.Project) error
	Delete(id uint) error
}

type projectRepo struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepo{db: db}
}

func (r *projectRepo) Create(project *domain.Project) error {
	if err := r.db.Create(project).Error; err != nil {
		return err
	}
	return nil
}

func (r *projectRepo) List(limit, offset int) ([]domain.Project, error) {
	var projects []domain.Project

	err := r.db.
		Preload("Creator").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&projects).Error

	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *projectRepo) Count() (int64, error) {
	var count int64

	err := r.db.
		Model(&domain.Project{}).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *projectRepo) GetByID(id uint) (*domain.Project, error) {
	var project domain.Project

	err := r.db.
		Preload("Creator").
		First(&project, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrProjectNotFound
	}
	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (r *projectRepo) Update(project *domain.Project) error {
	result := r.db.
		Model(&domain.Project{}).
		Where("id = ?", project.ID).
		Updates(domain.Project{
			Name:        project.Name,
			Description: project.Description,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrProjectNotFound
	}

	return nil
}

func (r *projectRepo) Delete(id uint) error {
	result := r.db.
		Where("id = ?", id).
		Delete(&domain.Project{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrProjectNotFound
	}

	return nil
}
