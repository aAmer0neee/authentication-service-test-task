package repository

import (
	"github.com/aAmer0neee/authentication-service-test-task/internal/config"
	"github.com/aAmer0neee/authentication-service-test-task/internal/domain"
	"github.com/aAmer0neee/authentication-service-test-task/internal/repository/postgres"
)

type Repository interface {
	AddRecord(record *domain.User, hash string) error
}

func New(cfg config.Cfg) (Repository, error) {
	return postgres.Connect(cfg)
}
