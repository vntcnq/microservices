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
	// Подключение к MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://mongo:27017"))
	if err != nil {
		log.Fatal(err)
	}

	// Коллекция продуктов
	productCollection = client.Database("ecommerce").Collection("products")

	// Установка режима релиза
	gin.SetMode(gin.ReleaseMode)

	// Настройка маршрутов
	router := gin.Default()

	// Установка доверенных прокси
	router.SetTrustedProxies([]string{"127.0.0.1", "0.0.0.0"})

	router.POST("/products", createProduct)
	router.GET("/products/:id", getProductByID)
	router.PATCH("/products/:id", updateProduct)
	router.DELETE("/products/:id", deleteProduct)
	router.GET("/products", listProducts)

	// Запуск Inventory Service на порту 8081
	log.Println("Inventory Service запущен на порту 8081")
	router.Run(":8081")
}

// Структура продукта
type Product struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

// Создание нового продукта
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

// Получение продукта по ID
func getProductByID(c *gin.Context) {
	id := c.Param("id")
	var product Product

	err := productCollection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&product)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Продукт не найден"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// Обновление продукта
func updateProduct(c *gin.Context) {
	id := c.Param("id")
	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := productCollection.UpdateOne(context.TODO(), bson.M{"id": id}, bson.M{"$set": product})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления продукта"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Продукт обновлен"})
}

// Удаление продукта
func deleteProduct(c *gin.Context) {
	id := c.Param("id")

	_, err := productCollection.DeleteOne(context.TODO(), bson.M{"id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления продукта"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Продукт удален"})
}

// Список всех продуктов
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
