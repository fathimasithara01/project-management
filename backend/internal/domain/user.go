package domain

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	Admin     Role = "admin"
	Developer Role = "developer"
)

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `gorm:"size:100;not null" json:"name"`
	Email    string `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`

	Role Role `gorm:"type:varchar(20);default:'developer'" json:"role"`

	Projects []Project `gorm:"foreignKey:CreatedBy" json:"projects,omitempty"`
	Tasks    []Task    `gorm:"foreignKey:AssignedTo" json:"tasks,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
