package controllers

import (
	"fmt"
	"inventory-api/config"
	"inventory-api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// RestockItem godoc
// @Summary Restock an inventory item
// @Description Restock an inventory item with the specified amount. Limits restocks to 3 per item in 24 hours.
// @Tags restock
// @Accept  json
// @Produce  json
// @Param restock body models.Restock true "Restock details"
// @Security ApiKeyAuth
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 429 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/restock [post]
func RestockItem(c *gin.Context) {
	var restock models.Restock
	if err := c.ShouldBindJSON(&restock); err != nil {
		fmt.Println("Error Binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var item models.InventoryItem
	config.DB.First(&item, restock.ItemID)
	if item.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	cutoffTime := time.Now().Add(-24 * time.Hour)

	var restockCount int64
	config.DB.Model(&models.Restock{}).
		Where("item_id = ? AND created_at >= ?", restock.ItemID, cutoffTime).
		Count(&restockCount)

	if restockCount >= 3 {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many restocks in 24 hours"})
		return
	}

	item.Quantity += restock.Amount
	config.DB.Save(&item)

	restock.CreatedAt = time.Now()
	config.DB.Create(&restock)

	c.JSON(http.StatusOK, gin.H{"message": "Restocked successfully"})
}

// GetRestockHistory godoc
// @Summary Get restock history for an item
// @Description Retrieve the restock history for a specific inventory item, sorted by most recent.
// @Tags restock
// @Accept  json
// @Produce  json
// @Param item_id path string true "Item ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/items/{item_id}/restock-history [get]
func GetRestockHistory(c *gin.Context) {
	var restocks []models.Restock
	var item models.InventoryItem
	itemID := c.Param("item_id")

	if err := config.DB.Where("id = ?", itemID).First(&item).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	if err := config.DB.Where("item_id = ?", itemID).Order("created_at DESC").Find(&restocks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve restock history"})
		return
	}

	var history []gin.H
	for _, restock := range restocks {
		history = append(history, gin.H{
			"timestamp": restock.CreatedAt,
			"amount":    restock.Amount,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"item_id": itemID,
		"name":    item.Name,
		"history": history,
	})
}

// GetAllRestockHistory godoc
// @Summary Get restock history for all items
// @Description Retrieve the restock history for all inventory items, sorted by most recent.
// @Tags restock
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /admin/items/restock-history [get]
func GetAllRestockHistory(c *gin.Context) {
	var restocks []models.Restock

	if err := config.DB.Order("created_at DESC").Find(&restocks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve restock history"})
		return
	}

	historyMap := make(map[string][]gin.H)
	for _, restock := range restocks {
		var item models.InventoryItem
		if err := config.DB.Where("id = ?", restock.ItemID).First(&item).Error; err != nil {
			continue
		}

		historyMap[item.Name] = append(historyMap[item.Name], gin.H{
			"item_id":   restock.ItemID,
			"timestamp": restock.CreatedAt,
			"amount":    restock.Amount,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"restock_history": historyMap,
	})
}
