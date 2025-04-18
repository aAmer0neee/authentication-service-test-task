package postgres

import (
	"fmt"
	"log"

	"github.com/aAmer0neee/authentication-service-test-task/internal/config"
	"github.com/aAmer0neee/authentication-service-test-task/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func Connect(cfg config.Cfg) (*PostgresRepository, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Postgres.Host,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Name,
		cfg.Postgres.Port,
		cfg.Postgres.Sslmode)
	postgres, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, err
	}
	fmt.Printf("[Repository] [INFO] Open Data Base %s\n", postgres.Name())

	if cfg.Postgres.Migrate {
		if err := postgres.AutoMigrate(&Users{}); err != nil {
			return nil, err
		}
		log.Printf("[Repository][INFO] AutoMigrate")
	}

	return &PostgresRepository{db: postgres}, nil
}

func (r *PostgresRepository) AddRecord(record *domain.User, hash string) error {
	return r.db.Create(&Users{
		Id:        record.Id,
		IpAddress: record.IpAddress.String(),
		TokenHash: hash,
		Email:     record.Email,
	}).Error
}
