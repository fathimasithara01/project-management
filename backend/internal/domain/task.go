package domain

import (
	"time"

	"gorm.io/gorm"
)

type Status string

const (
	Todo       Status = "todo"
	InProgress Status = "in_progress"
	Done       Status = "done"
)

type Task struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Title       string `gorm:"size:150;not null" json:"title"`
	Description string `gorm:"type:text" json:"description"`
	Status      Status `gorm:"type:varchar(20);default:'todo'" json:"status"`

	ProjectID uint
	Project   Project `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"project,omitempty"`

	AssignedTo *uint `json:"assigned_to,omitempty"`
	Assignee   User  `gorm:"foreignKey:AssignedTo;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"assignee,omitempty"`

	DueDate *time.Time `json:"due_date,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
