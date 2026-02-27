package main

import (
	"log"

	"saas-platform/internal/cache"
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

	redisClient, err := cache.New(cfg.RedisAddr)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	tenantHandler := handlers.NewTenantHandler(database)
	apiKeyHandler := handlers.NewAPIKeyHandler(database)

	r.POST("/tenants", tenantHandler.CreateTenant)
	r.POST("/tenants/:tenant_id/keys", apiKeyHandler.CreateKey)

	auth := r.Group("/")
	// Order matters: Auth first (validate key), then rate limit (enforce quota)
	auth.Use(middleware.RateLimitRedis(redisClient, cfg.RateLimitRPM))
	auth.Use(middleware.APIKeyAuth(database))

	auth.GET("/protected", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	r.Run(":8080")
}