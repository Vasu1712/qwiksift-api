package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Vasu1712/qwiksift/server/internal/scraper"
	"github.com/gin-gonic/gin"
)

// GetProducts handles the /api/products endpoint
func GetProducts(c *gin.Context) {
    products, err := scraper.ScrapeAll()
    if err != nil {
        fmt.Printf("Scraping error: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to fetch products",
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "count":    len(products),
        "products": products,
    })
}

// GetBlinkitProducts handles the /blinkit-api/products/listing endpoint
func GetBlinkitProducts(c *gin.Context) {
    // Get search query from request
    searchQuery := c.Query("query")
    if searchQuery == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
        return
    }
    
    // Get category IDs based on search query
    l0Cat, l1Cat := scraper.FindCategoryForSearchTerm(searchQuery)
    
    // Get location from query params or use default
    lat := c.DefaultQuery("lat", "12.9337741")
    lon := c.DefaultQuery("lon", "77.7006304")
    
    // Create Blinkit scraper instance
    blinkit := scraper.NewBlinkitScraper()
    
    // Fetch products data
    jsonData, err := blinkit.FetchCategoryProducts(l0Cat, l1Cat, lat, lon)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }
    
    // Return the raw JSON response
    c.Data(http.StatusOK, "application/json", jsonData)
}

// HealthCheck handles the /health endpoint
func HealthCheck(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status":  "ok",
        "version": "1.0.0",
        "time":    time.Now().Unix(),
    })
}
