package models

import "gorm.io/gorm"

// User represents a user in the system
// @Description User represents a user in the system
type User struct {
	gorm.Model `swaggerignore:"true"`
	Email      string `gorm:"unique" json:"email" binding:"required,email" example:"user@example.com"`
	Password   string `json:"password" binding:"required" example:"password123"`
}
