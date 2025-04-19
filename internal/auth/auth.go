package auth

import (
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/aAmer0neee/authentication-service-test-task/internal/config"
	"github.com/aAmer0neee/authentication-service-test-task/internal/domain"
	"github.com/aAmer0neee/authentication-service-test-task/internal/logger"
	"github.com/aAmer0neee/authentication-service-test-task/internal/repository"
	"github.com/google/uuid"
)

type AuthService interface {
	LoginUser(user *domain.User)(*domain.Tokens, error)
	RefreshToken(inputUser *domain.User)(*domain.Tokens, error)
}

type authService struct {
	repo repository.Repository
	log  *slog.Logger
	jwt  *JWTService
}

type keys struct {
	access string
	refresh string

	refreshHash string
}

func New(r repository.Repository, cfg config.Cfg) AuthService {
	return &authService{
		repo: r,
		jwt:  configureJWT(cfg.AuthSecret),
		log:  logger.ConfigureLogger(cfg.Server.Env),
	}
}

func (s *authService) LoginUser(user *domain.User) (*domain.Tokens, error) {
	
	keys, err := s.getKeysPair(user)
	if err != nil {
		return nil,err
	}
	
	if err := s.repo.AddRecord(user, keys.refreshHash); err != nil {
		s.log.Info("cant't add record to db", "message", err)
		return nil,err
	}

	user.AccessToken = keys.access
	user.RefreshToken = keys.refresh
	return &domain.Tokens{Access: keys.access,
		Refresh: keys.refresh,
		},nil
}

func (s *authService) RefreshToken(inputUser *domain.User)(*domain.Tokens, error) {

	claims, err := s.jwt.validate(*inputUser)
	if 
	err != nil ||
	time.Now().After(time.Unix(int64(claims.ExpiredAt), 0)) {
		s.log.Info("access token invalid", "message", err)
		return nil,err
	}

	registerUser, err := s.repo.GetRecord(uuid.MustParse(claims.UserId))
	if err != nil {
		s.log.Info("can't find user at data base", "message", err)
		return nil,err
	}

	if err := s.jwt.compareBcryptTokens(
		registerUser.RefreshToken, inputUser.RefreshToken);
		err != nil {
			s.log.Info("refresh tokens not equal", "message", err)
			fmt.Println(registerUser.RefreshToken)
			return nil, err
		}
	
	if registerUser.IpAddress.Equal(inputUser.IpAddress) ||
	inputUser.IpAddress.Equal(net.ParseIP(claims.IpAddress)) {
		fmt.Printf("implement send warning to %s\n",registerUser.Email)
	}

	keys, err := s.getKeysPair(registerUser)
	if err != nil {
		return nil,err
	}
	fmt.Println("Generated refreshHash:", keys.refreshHash) 
	if err := s.repo.UpdateRecord(inputUser, keys.refreshHash); err  != nil {
		s.log.Info("can't update record in data base", "message", err)
		return nil,err
	}

	return &domain.Tokens{Access: keys.access,
		Refresh: keys.refresh,
		},nil

}

func (s *authService)getKeysPair(user *domain.User)(keys, error){

	access, err := s.jwt.generateAccess(*user)
	if err != nil {
		s.log.Info("errror generate access key", "message", err)
		return keys{},err
	}

	refresh := s.jwt.generateRefresh()

	refreshHash, err := s.jwt.generateBcryptToken(refresh)
	if err != nil {
		s.log.Info("errror hashing token", "message", err)
		return keys{},err
	}

	return keys{
		access: access,
		refresh: refresh,
		refreshHash: refreshHash,
	}, nil
}