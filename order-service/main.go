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
	// Подключение к MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://mongo:27017"))
	if err != nil {
		log.Fatal(err)
	}

	// Коллекция заказов
	orderCollection = client.Database("ecommerce").Collection("orders")

	// Настройка маршрутов
	router := gin.Default()
	router.POST("/orders", createOrder)
	router.GET("/orders/:id", getOrderByID)
	router.PATCH("/orders/:id", updateOrder)
	router.GET("/orders", listOrders)

	// Запуск Order Service на порту 8082
	log.Println("Order Service запущен на порту 8082")
	router.Run(":8082")
}

// Структура заказа
type Order struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Status    string `json:"status"`
}

// Создание нового заказа
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

// Получение заказа по ID
func getOrderByID(c *gin.Context) {
	id := c.Param("id")
	var order Order

	err := orderCollection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&order)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Заказ не найден"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// Обновление статуса заказа
func updateOrder(c *gin.Context) {
	id := c.Param("id")
	var order Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := orderCollection.UpdateOne(context.TODO(), bson.M{"id": id}, bson.M{"$set": order})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления заказа"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Заказ обновлен"})
}

// Список всех заказов
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
