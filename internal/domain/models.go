package domain

import (
	"net"

	"github.com/google/uuid"
)

type RefreshRequest struct {
	AccessToken  string `json:"access-token" binding:"required"`
	RefreshToken string `json:"refresh-token" binding:"required"`
}

type LoginRequest struct {
	Id    uuid.UUID `json:"id" binding:"required"`
	Email string    `json:"email"`
}

type User struct {
	Id           uuid.UUID
	Email        string
	IpAddress    net.IP
	AccessToken  string
	RefreshToken string
}
