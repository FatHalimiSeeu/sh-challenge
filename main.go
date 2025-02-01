package main

import (
	"inventory-api/config"
	"inventory-api/routes"

	"github.com/gin-gonic/gin"
)

// @title Inventory API
// @version 1.0
// @description This is a sample server for an inventory management system. This API requires an API key for authentication. Use the "Authorize" button to enter your API key that you will obtain from the response of the Login method in the header.
// @termsOfService http://swagger.io/terms/

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @host localhost:8080
// @BasePath /
// @description This API uses an API key authentication method. Include your API key in the Authorization header.

func main() {
	config.InitDB()
	config.InitRedis()
	r := gin.Default()
	routes.RegisterRoutes(r)

	r.Run(":8080")
}
