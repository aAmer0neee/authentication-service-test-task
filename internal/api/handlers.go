package api

import "github.com/google/uuid"

type RefreshRequest struct {
	AccessToken  string `json:"access" binding:"required"`
	RefreshToken string `json:"refresh" binding:"required"`
}

type LoginRequest struct {
	Id    uuid.UUID `json:"id" binding:"required"`
	Email string    `json:"email"`
}
