package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var productCollection *mongo.Collection

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://mongo:27017"))
	if err != nil {
		log.Fatal(err)
	}

	productCollection = client.Database("ecommerce").Collection("products")

	router := gin.Default()
	router.POST("/products", createProduct)
	router.GET("/products", listProducts)

	log.Println("Inventory Service запущен на порту 8081")
	router.Run(":8081")
}

type Product struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

func createProduct(c *gin.Context) {
	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := productCollection.InsertOne(context.TODO(), product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания продукта"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Продукт добавлен"})
}

func listProducts(c *gin.Context) {
	cursor, err := productCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения списка продуктов"})
		return
	}
	defer cursor.Close(context.TODO())

	var products []Product
	if err := cursor.All(context.TODO(), &products); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка парсинга данных"})
		return
	}

	c.JSON(http.StatusOK, products)
}
