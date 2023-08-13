package models

import (
	"time"

	"github.com/google/uuid"
)

type Note struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	UserID    uuid.UUID `gorm:"type:uuid" json:"userId,omitempty"`
	Title     string    `gorm:"varchar(255);uniqueIndex;not null" json:"title,omitempty"`
	Content   string    `gorm:"not null" json:"content,omitempty"`
	Category  string    `gorm:"varchar(100)" json:"category,omitempty"`
	Published bool      `gorm:"default:false;not null" json:"published"`
	CreatedAt time.Time `gorm:"not null" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"not null" json:"updatedAt,omitempty"`
}

type CreateNoteSchema struct {
	Title     string `json:"title" validate:"required"`
	Content   string `json:"content" validate:"required"`
	Category  string `json:"category,omitempty"`
	Published bool   `json:"published,omitempty"`
}

type UpdateNoteSchema struct {
	Title     string `json:"title,omitempty"`
	Content   string `json:"content,omitempty"`
	Category  string `json:"category,omitempty"`
	Published *bool  `json:"published,omitempty"`
}
