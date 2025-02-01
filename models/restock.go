package models

import (
	"time"

	"gorm.io/gorm"
)

// Restock represents a restock event for an inventory item
// @Description Restock represents a restock event for an inventory item
type Restock struct {
	gorm.Model `swaggerignore:"true"`
	ItemID     uint      `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"item_id" example:"1"` // Foreign key to InventoryItem
	Amount     int       `json:"amount" binding:"required,min=10,max=1000" example:"100"`
	CreatedAt  time.Time `json:"created_at" example:"2025-01-31T19:30:08Z" swaggerignore:"true"`
}
