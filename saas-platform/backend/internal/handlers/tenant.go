package handlers

import (
    "context"
    "net/http"

    "saas-platform/internal/models"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
)

type TenantHandler struct {
    DB *pgxpool.Pool
}

func NewTenantHandler(db *pgxpool.Pool) *TenantHandler {
    return &TenantHandler{DB: db}
}

func (h *TenantHandler) CreateTenant(c *gin.Context) {
    var t models.Tenant
    if err := c.BindJSON(&t); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    t.ID = uuid.New()

    _, err := h.DB.Exec(context.Background(),
        `INSERT INTO tenants(id, name) VALUES($1,$2)`,
        t.ID, t.Name,
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, t)
}