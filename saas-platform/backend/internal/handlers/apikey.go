package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type APIKeyHandler struct {
	DB *pgxpool.Pool
}

func NewAPIKeyHandler(db *pgxpool.Pool) *APIKeyHandler {
	return &APIKeyHandler{DB: db}
}

func (h *APIKeyHandler) CreateKey(c *gin.Context) {
	tenantID := c.Param("tenant_id")

	key := uuid.New().String()
	id := uuid.New()

	_, err := h.DB.Exec(context.Background(),
		`INSERT INTO api_keys(id, tenant_id, key) VALUES($1,$2,$3)`,
		id, tenantID, key,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"api_key": key})
}
