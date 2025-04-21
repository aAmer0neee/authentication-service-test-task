package domain

import (
	"net"

	"github.com/google/uuid"
)

type User struct {
	Id           uuid.UUID
	Email        string
	IpAddress    net.IP
	AccessToken  string
	RefreshToken string
	TokenPairId  uuid.UUID
}

type Tokens struct {
	Access  string
	Refresh string
}
