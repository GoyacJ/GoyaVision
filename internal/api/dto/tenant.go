package dto

import (
	"time"

	"github.com/google/uuid"
)

type TenantCreateReq struct {
	Name   string `json:"name" validate:"required"`
	Code   string `json:"code" validate:"required"`
	Status int    `json:"status"`
}

type TenantUpdateReq struct {
	Name   *string `json:"name"`
	Status *int    `json:"status"`
}

type TenantResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
