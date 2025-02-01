package controllers

import (
	"inventory-api/config"
	"inventory-api/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CreateItem godoc
// @Summary Create a new inventory item
// @Description Create a new inventory item with the input payload
// @Tags inventory
// @Accept  json
// @Produce  json
// @Param item body models.InventoryItem true "Create inventory item"
// @Security ApiKeyAuth
// @Success 200 {object} models.InventoryItem
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/items [post]
func CreateItem(c *gin.Context) {
	var item models.InventoryItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	// config.DB.Create(&item)

	// database level validation - i chose this because it will not query
	// the database without need to check before inserting
	if err := config.DB.Create(&item).Error; err != nil {
		// Check for unique constraint violation (Postgres & SQLite return this error)
		if strings.Contains(err.Error(), "duplicate key value") || strings.Contains(err.Error(), "UNIQUE constraint failed") {
			c.JSON(http.StatusConflict, gin.H{"error": "Item with this name already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item"})
		}
		return
	}
	c.JSON(http.StatusOK, item)
}

// ListItems godoc
// @Summary List all inventory items
// @Description Get a list of all inventory items, optionally filtered by low stock
// @Tags inventory
// @Accept  json
// @Produce  json
// @Param lowStock query string false "Filter low stock items (true/false)"
// @Success 200 {array} models.InventoryItem
// @Failure 500 {object} map[string]string
// @Router /admin/items [get]
func ListItems(c *gin.Context) {
	var items []models.InventoryItem
	lowStock := c.Query("lowStock")
	if lowStock == "true" {
		config.DB.Where("quantity <= ?", 20).Find(&items)
	} else {
		config.DB.Find(&items)
	}
	c.JSON(http.StatusOK, items)
}
