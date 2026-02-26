package main

import (
    "log"

    "saas-platform/internal/config"
    "saas-platform/internal/db"
    "saas-platform/internal/handlers"
    "saas-platform/internal/middleware"

    "github.com/gin-gonic/gin"
)

func main() {
    cfg := config.Load()

    database, err := db.New(cfg)
    if err != nil {
        log.Fatal(err)
    }

    r := gin.Default()

    tenantHandler := handlers.NewTenantHandler(database)

    r.POST("/tenants", tenantHandler.CreateTenant)

    auth := r.Group("/")
    auth.Use(middleware.APIKeyAuth(database))
    auth.GET("/protected", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "ok"})
    })

    r.Run(":8080")
}