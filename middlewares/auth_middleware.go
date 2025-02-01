package middlewares

import (
	"fmt"
	"inventory-api/config"
	"inventory-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func Authenticate(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	fmt.Println("Authorization Header:", tokenString) // Log received token

	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		c.Abort()
		return
	}

	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	storedToken, err := config.RedisClient.Get(config.Ctx, claims.Email).Result()
	if err == redis.Nil || storedToken != tokenString {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired or invalid"})
		c.Abort()
		return
	}

	c.Set("email", claims.Email)
	c.Next()
}
