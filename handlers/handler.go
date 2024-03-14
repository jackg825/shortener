package handlers

import (
	"net/http"
	"shortener/cache"
	"shortener/storage"
	"shortener/utils"

	"github.com/gin-gonic/gin"

	"log"
)

type URLCreateRequest struct {
	LongURL string `json:"longUrl" binding:"required"`
	UserId  string `json:"userId" binding:"required"`
}

func CreateShortURL(store *storage.Storage, cache *cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req URLCreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("Error binding JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hashedURL := utils.GenerateShortURL(req.LongURL)
		shortURL, err := store.SaveShortURL(hashedURL, req.LongURL, req.UserId)
		if err != nil {
			log.Printf("Failed to save URL: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save URL"})
			return
		}
		cache.Set(shortURL, req.LongURL, 3600) // Cache for 1 hour
		store.SaveShortURL(shortURL, req.LongURL, req.UserId)

		c.JSON(http.StatusOK, gin.H{"shortUrl": shortURL})
	}
}

func RedirectShortURL(store *storage.Storage, cache *cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortURL := c.Param("shortUrl")
		longURL, found := cache.Get(shortURL)
		if !found {
			var err error
			longURL, err = store.GetLongURL(shortURL)
			if err != nil {
				log.Printf("URL not found: %v", err)
				c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
				return
			}
			cache.Set(shortURL, longURL, 3600) // Cache for 1 hour
		}

		c.Redirect(http.StatusMovedPermanently, longURL)
	}
}
