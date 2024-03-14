package main

import (
	"shortener/cache"
	"shortener/handlers"
	"shortener/storage"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := gin.Default()
	store := storage.NewStorage(os.Getenv("DATABASE_URL"))
	cache := cache.NewCache(os.Getenv("CACHE_URL"))
	// analytics := analytics.NewAnalytics()

	// Adjust handlers to use analytics
	router.POST("/create", handlers.CreateShortURL(store, cache))
	router.GET("/:shortUrl", handlers.RedirectShortURL(store, cache))

	router.Run(":8080")
}
