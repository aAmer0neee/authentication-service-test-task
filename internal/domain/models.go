package domain

import (
	"net"

	"github.com/google/uuid"
)

type RefreshRequest struct {
	AccessToken  string `json:"Access" binding:"required"`
	RefreshToken string `json:"Refresh" binding:"required"`
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

type AccessClaims struct {
	UserId string
	ExpiredAt float64
	IpAddress string
}

type Tokens struct {
	Access  string
	Refresh string
}