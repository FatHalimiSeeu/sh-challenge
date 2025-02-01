package controllers

import (
	"inventory-api/config"
	"inventory-api/models"
	"inventory-api/utils"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary Register a new user
// @Description Register a new user with the input payload
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body models.User true "Register user"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /register [post]
func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	hashedPassword, _ := utils.HashPassword(user.Password)
	user.Password = hashedPassword
	// config.DB.Create(&user)
	if err := config.DB.Create(&user).Error; err != nil {
		// Handle unique constraint error
		if strings.Contains(err.Error(), "duplicate key value") || strings.Contains(err.Error(), "Duplicate entry") {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
			return
		}
		// Handle unexpected database errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not register user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User registered"})
}

// Login godoc
// @Summary Login a user
// @Description Login with email and password
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body models.User true "Login user"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /login [post]
func Login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var dbUser models.User
	config.DB.First(&dbUser, "email = ?", user.Email)
	if dbUser.ID == 0 || !utils.CheckPassword(dbUser.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := utils.GenerateToken(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	err = config.RedisClient.Set(config.Ctx, user.Email, token, 24*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
