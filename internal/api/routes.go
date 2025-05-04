package api

import (
    "github.com/gin-gonic/gin"
)

// SetupRoutes registers all the API routes
func SetupRoutes(router *gin.Engine) {
    router.GET("/api/products", GetProducts)
    router.GET("/blinkit-api/products/listing", GetBlinkitProducts)
    router.GET("/health", HealthCheck)
}
