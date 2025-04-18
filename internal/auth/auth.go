package auth

import (
	"log/slog"

	"github.com/aAmer0neee/authentication-service-test-task/internal/config"
	"github.com/aAmer0neee/authentication-service-test-task/internal/domain"
	"github.com/aAmer0neee/authentication-service-test-task/internal/logger"
	"github.com/aAmer0neee/authentication-service-test-task/internal/repository"
)

type AuthService interface {
	LoginUser(user *domain.User) error
	RefreshToken(tokens *domain.User) error
}

type authService struct {
	repo repository.Repository
	log  *slog.Logger
	jwt  *JWTService
}

func New(r repository.Repository, cfg config.Cfg) AuthService {
	return &authService{
		repo: r,
		jwt:  configureJWT(cfg.AuthSecret),
		log:  logger.ConfigureLogger(cfg.Server.Env),
	}
}

func (s *authService) LoginUser(user *domain.User) error {
	access, err := s.jwt.generateAccess(*user)
	if err != nil {
		s.log.Info("errror generate access key", "message", err)
		return err
	}
	user.AccessToken = access
	user.RefreshToken = s.jwt.generateRefresh()

	refreshHash, err := s.jwt.generateBcryptToken(user.RefreshToken)
	if err != nil {
		s.log.Info("errror hashing token", "message", err)
		return err
	}

	if err := s.repo.AddRecord(user, refreshHash); err != nil {
		s.log.Info("errror adding record to db", "message", err)
		return err
	}
	return nil
}

func (s *authService) RefreshToken(user *domain.User) error {

	return 	s.jwt.validate(*user)
}
