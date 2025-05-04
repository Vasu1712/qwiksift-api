package main

import (
    "github.com/gin-gonic/gin"
    "github.com/Vasu1712/qwiksift-api/internal/api"
)

func main() {
    router := gin.Default()
    
    // Register all routes
    api.SetupRoutes(router)
    
    // Start the server
    router.Run(":8080")
}
