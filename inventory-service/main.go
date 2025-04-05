package main

import (
	"log"
	"net/http"
	"os"

	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"eccomerce-microservices/inventory-service/config"
	"eccomerce-microservices/inventory-service/routes"
)

func main() {
	router := gin.Default()

	// Middleware and routes setup
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Load configuration
	config.LoadConfig()

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB!")

	// Initialize routes
	routes.InitializeRoutes(router, client)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	router.Run(":8081")
}
