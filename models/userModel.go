package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
	Username      string         `json:"username" gorm:"unique" binding:"required"`
	Password      string         `json:"-" binding:"required"`
	User_type     string         `json:"user_type" binding:"required"`
	Refresh_token string         `json:"-"`
}
