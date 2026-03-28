package usecase

import (
	"errors"
	"project-management/internal/domain"
	"project-management/internal/repositories"
	"project-management/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(name, email, password, role string) (*domain.User, error)
	Login(email, password string) (string, *domain.User, error)
	ListUsers() ([]domain.User, error)
	GetUserByID(id uint) (*domain.User, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(r repositories.UserRepository) UserService {
	return &userService{r}
}

func (u *userService) Register(name, email, password, role string) (*domain.User, error) {

	existing, err := u.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	if role == "" {
		role = string(domain.Developer)
	}

	switch domain.Role(role) {
	case domain.Admin, domain.Developer:
	default:
		return nil, errors.New("invalid role")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Name:     name,
		Email:    email,
		Password: string(hashed),
		Role:     domain.Role(role),
	}

	if err := u.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userService) Login(email, password string) (string, *domain.User, error) {

	user, err := u.repo.GetByEmail(email)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		return "", nil, err
	}

	user.Password = ""

	return token, user, nil
}

func (u *userService) ListUsers() ([]domain.User, error) {
	return u.repo.List()
}

func (u *userService) GetUserByID(id uint) (*domain.User, error) {
	return u.repo.GetByID(id)
}
