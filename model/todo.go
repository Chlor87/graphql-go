package model

import (
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	ID          int            `json:"id" gorm:"primarykey"`
	Text        string         `json:"text" gorm:"uniqueIndex;not null;check:text <> ''"`
	Done        bool           `json:"done"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt" gorm:"index"`
	UpdatedByID int            `json:"-" gorm:"not null"`
	UpdatedBy   *User          `json:"updatedBy"`
}
