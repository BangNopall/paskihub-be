package dto

import "github.com/google/uuid"

type JuryRequest struct {
	Name string `json:"name" binding:"required"`
}

type JuryResponse struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
