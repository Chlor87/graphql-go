package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        int            `json:"id"`
	Name      string         `json:"name" gorm:"uniqueIndex;not null;check:name <> ''"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null;check:email <> ''"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
	UpdatedBy string         `json:"updatedBy" gorm:"not null"`
}
