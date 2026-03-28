package domain

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"size:150;not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`

	CreatedBy *uint `json:"created_by,omitempty"`
	Creator   *User `gorm:"foreignKey:CreatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"creator,omitempty"`

	Tasks []Task `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"tasks,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
