package usecase

import (
	"errors"
	"project-management/internal/domain"
	"project-management/internal/repositories"

	"gorm.io/gorm"
)

var ErrProjectNotFound = errors.New("project not found")

type ProjectService interface {
	CreateProject(name, description string, userID uint) (*domain.Project, error)
	ListProjects(page, limit int) ([]domain.Project, int64, error)
	GetProjectByID(id uint) (*domain.Project, error)
	UpdateProject(id uint, name, description string) (*domain.Project, error)
	DeleteProject(id uint) error
}

type projectService struct {
	repo repositories.ProjectRepository
}

func NewProjectService(repo repositories.ProjectRepository) ProjectService {
	return &projectService{repo}
}

func (p *projectService) CreateProject(name, description string, userID uint) (*domain.Project, error) {
	project := &domain.Project{
		Name:        name,
		Description: description,
		CreatedBy:   &userID,
	}

	if err := p.repo.Create(project); err != nil {
		return nil, err
	}

	return project, nil
}

func (p *projectService) ListProjects(page, limit int) ([]domain.Project, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	projects, err := p.repo.List(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := p.repo.Count()
	if err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

func (p *projectService) GetProjectByID(id uint) (*domain.Project, error) {
	project, err := p.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if project == nil {
		return nil, ErrProjectNotFound
	}

	return project, nil
}

func (p *projectService) UpdateProject(id uint, name, description string) (*domain.Project, error) {
	project, err := p.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if project == nil {
		return nil, ErrProjectNotFound
	}

	project.Name = name
	project.Description = description

	if err := p.repo.Update(project); err != nil {
		return nil, err
	}

	return project, nil
}

func (p *projectService) DeleteProject(id uint) error {
	err := p.repo.Delete(id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrProjectNotFound
	}

	return err
}
