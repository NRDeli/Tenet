package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// Simple fixed-window limiter: allow N requests per minute per api key.
// Key: rl:<api_key>:<YYYYMMDDHHMM>
func RateLimitRedis(rdb *redis.Client, rpm int) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			// Let APIKeyAuth handle missing key (keeps single source of truth)
			c.Next()
			return
		}

		now := time.Now().UTC()
		window := now.Format("200601021504") // minute bucket
		redisKey := fmt.Sprintf("rl:%s:%s", apiKey, window)

		ctx, cancel := context.WithTimeout(c.Request.Context(), 200*time.Millisecond)
		defer cancel()

		// Atomic increment; set expiry on first hit.
		n, err := rdb.Incr(ctx, redisKey).Result()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"error": "rate limiter unavailable"})
			return
		}
		if n == 1 {
			_ = rdb.Expire(ctx, redisKey, 70*time.Second).Err()
		}

		if int(n) > rpm {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
				"rpm":   rpm,
			})
			return
		}

		c.Next()
	}
}