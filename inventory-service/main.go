package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	var err error

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Product{})

	router := gin.Default()

	router.POST("/products", createProduct)
	router.GET("/products/:id", getProductByID)
	router.PATCH("/products/:id", updateProduct)
	router.DELETE("/products/:id", deleteProduct)
	router.GET("/products", listProducts)

	log.Println("Inventory Service запущен на порту 8081")
	router.Run(":8081")
}

type Product struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

func createProduct(c *gin.Context) {
	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Create(&product)
	c.JSON(http.StatusCreated, gin.H{"message": "Продукт добавлен"})
}

func getProductByID(c *gin.Context) {
	var product Product
	id := c.Param("id")
	result := db.First(&product, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Продукт не найден"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func updateProduct(c *gin.Context) {
	var product Product
	id := c.Param("id")
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Model(&product).Where("id = ?", id).Updates(product)
	c.JSON(http.StatusOK, gin.H{"message": "Продукт обновлен"})
}

func deleteProduct(c *gin.Context) {
	id := c.Param("id")
	db.Delete(&Product{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Продукт удален"})
}

func listProducts(c *gin.Context) {
	var products []Product
	db.Find(&products)
	c.JSON(http.StatusOK, products)
}
