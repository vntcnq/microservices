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

var orderCollection *mongo.Collection

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://mongo:27017"))
	if err != nil {
		log.Fatal(err)
	}

	orderCollection = client.Database("ecommerce").Collection("orders")

	router := gin.Default()
	router.POST("/orders", createOrder)
	router.GET("/orders", listOrders)

	log.Println("Order Service запущен на порту 8082")
	router.Run(":8082")
}

type Order struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

func createOrder(c *gin.Context) {
	var order Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := orderCollection.InsertOne(context.TODO(), order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания заказа"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Заказ создан"})
}

func listOrders(c *gin.Context) {
	cursor, err := orderCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения заказов"})
		return
	}
	defer cursor.Close(context.TODO())

	var orders []Order
	if err := cursor.All(context.TODO(), &orders); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка парсинга данных"})
		return
	}

	c.JSON(http.StatusOK, orders)
}
