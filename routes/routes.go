package routes

import (
	"inventory-api/controllers"
	"inventory-api/middlewares"

	_ "inventory-api/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(r *gin.Engine) {
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Admin routes (protected by authentication middleware)
	admin := r.Group("/admin")
	admin.GET("/items", controllers.ListItems)
	admin.GET("/items/:item_id/restock-history", controllers.GetRestockHistory)
	admin.GET("/items/restock-history", controllers.GetAllRestockHistory)
	admin.Use(middlewares.Authenticate)
	{
		admin.POST("/items", controllers.CreateItem)
		admin.POST("/restock", controllers.RestockItem)
		// admin.GET("/items", controllers.ListItems)
		// admin.GET("/items/:item_id/restock-history", controllers.GetRestockHistory)
	}
}
