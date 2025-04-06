package main

import (
    "net/http"
    "time"
    "fmt"
    
    "github.com/gin-gonic/gin"
    "github.com/Vasu1712/qwiksift/server/internal/scraper"
)

func main() {
    router := gin.Default()
    
    // API Endpoint
    router.GET("/api/products", func(c *gin.Context) {
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
    })

    router.GET("/api/products/listing", func(c *gin.Context) {
        // Get category parameters from query string
        l0Cat := c.DefaultQuery("l0_cat", "14")
        l1Cat := c.DefaultQuery("l1_cat", "922")
        
        // Use the specific lat/lon that works with Blinkit's API
        lat := "12.9337741"
        lon := "77.7006304"
        
        // Create Blinkit scraper instance
        blinkit := scraper.NewBlinkitScraper()
        
        // Use existing FetchCategoryProducts function
        jsonData, err := blinkit.FetchCategoryProducts(l0Cat, l1Cat, lat, lon)
        
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": err.Error(),
            })
            return
        }
        
        // Return the raw JSON response exactly as received from Blinkit
        c.Data(http.StatusOK, "application/json", jsonData)
    })

    
    // Health Check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status":  "ok",
            "version": "1.0.0",
            "time":    time.Now().Unix(),
        })
    })
    
    router.Run(":8080")
}
