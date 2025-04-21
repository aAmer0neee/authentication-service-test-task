package repository

import (
	"github.com/aAmer0neee/authentication-service-test-task/internal/config"
	"github.com/aAmer0neee/authentication-service-test-task/internal/domain"
	"github.com/aAmer0neee/authentication-service-test-task/internal/repository/postgres"
	"github.com/google/uuid"
)

//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go -package=repository_mock
type Repository interface {
	AddRecord(record *domain.User, hash string, pairId uuid.UUID) error
	GetRecord(id uuid.UUID) (*domain.User, error)
	UpdateRecord(record *domain.User, hash string, pairId uuid.UUID) error
}

func New(cfg config.Cfg) (Repository, error) {
	return postgres.Connect(cfg)
}
