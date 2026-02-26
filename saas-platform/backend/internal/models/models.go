package models

import "github.com/google/uuid"

type Tenant struct {
    ID   uuid.UUID `json:"id"`
    Name string    `json:"name"`
}

type APIKey struct {
    ID       uuid.UUID `json:"id"`
    TenantID uuid.UUID `json:"tenant_id"`
    Key      string    `json:"key"`
}