package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode) // Установка режима релиза
	router := gin.Default()

	// Установка доверенных прокси
	router.SetTrustedProxies([]string{"127.0.0.1", "0.0.0.0"})

	// Прокси для Inventory Service
	router.Any("/inventory/*any", proxyHandler("http://inventory-service:8081"))

	// Прокси для Order Service
	router.Any("/orders/*any", proxyHandler("http://order-service:8082"))

	// Прокси для продуктов
	router.Any("/products/*any", proxyHandler("http://inventory-service:8081/products"))

	// Прокси для корневого пути
	router.Any("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to API Gateway"})
	})

	log.Println("API Gateway запущен на порту 8080")
	router.Run(":8080")
}

func proxyHandler(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		proxyURL := target + c.Param("any")
		resp, err := http.Get(proxyURL)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "Ошибка проксирования"})
			return
		}
		defer resp.Body.Close()

		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	}
}
