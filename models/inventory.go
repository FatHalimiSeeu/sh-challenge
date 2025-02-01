package models

import "gorm.io/gorm"

// InventoryItem represents an item in the inventory
// @Description InventoryItem represents an item in the inventory
type InventoryItem struct {
	gorm.Model  `swaggerignore:"true"`
	Name        string    `gorm:"unique" json:"name" binding:"required" example:"Item A"`
	Description string    `json:"description" example:"This is a sample item"`
	Quantity    int       `json:"quantity" example:"100"`
	Restocks    []Restock `gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" swaggerignore:"true"` // One-to-Many relationship, ignored in Swagger
}
