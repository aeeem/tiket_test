package helper

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        uint           `json:"id" example:"1" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at" example:"1994-12-15T13:47" gorm:"index"`
	UpdatedAt time.Time      `json:"updated_at" example:"1994-12-15T13:47"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
