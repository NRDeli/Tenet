package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func APIKeyAuth(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-API-Key")
		if key == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing api key"})
			return
		}

		var exists bool
		err := db.QueryRow(context.Background(),
			`SELECT EXISTS(SELECT 1 FROM api_keys WHERE key=$1)`,
			key,
		).Scan(&exists)

		if err != nil || !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid api key"})
			return
		}

		_, _ = db.Exec(context.Background(),
			`INSERT INTO usage_logs(id, api_key, path) VALUES($1,$2,$3)`,
			uuid.New(), key, c.Request.URL.Path,
		)

		c.Next()
	}
}
