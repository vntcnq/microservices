package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"eccomerce-microservices/inventory-service/controllers"
)

func InitializeRoutes(router *gin.Engine, client *mongo.Client) {
	productController := controllers.NewProductController(client)

	router.POST("/products", productController.CreateProduct)
	router.GET("/products/:id", productController.GetProduct)
	router.PATCH("/products/:id", productController.UpdateProduct)
	router.DELETE("/products/:id", productController.DeleteProduct)
	router.GET("/products", productController.ListProducts)
}
