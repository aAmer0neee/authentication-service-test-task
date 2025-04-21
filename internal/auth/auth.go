package auth

import (
	"errors"
	"github.com/aAmer0neee/authentication-service-test-task/internal/domain"
	"github.com/aAmer0neee/authentication-service-test-task/internal/notify"
	"github.com/aAmer0neee/authentication-service-test-task/internal/repository"
	"github.com/aAmer0neee/authentication-service-test-task/internal/token"
	"github.com/google/uuid"
	"log/slog"
	"net"
	"time"
)

//go:generate mockgen -source=auth.go -destination=mocks/auth_mock.go -package=auth_mock
type AuthService interface {
	LoginUser(user *domain.User) (*domain.Tokens, error)
	RefreshToken(inputUser *domain.User) (*domain.Tokens, error)
}

type authService struct {
	repo   repository.Repository
	log    *slog.Logger
	token  token.TokenService
	notify notify.Notifyer
}

type keys struct {
	pairId      uuid.UUID
	access      string
	refresh     string
	refreshHash string
}

var (
	ErrorInvalidFormat = errors.New("некорректный формат ввода")
	ErrorInvalidTokens = errors.New("невалидные токены")
	ErrorTokenExpired  = errors.New("токен устарел")
)

func New(r repository.Repository, token token.TokenService, notify notify.Notifyer, logger *slog.Logger) AuthService {
	return &authService{
		repo:   r,
		token:  token,
		log:    logger,
		notify: notify,
	}
}

func (s *authService) LoginUser(user *domain.User) (*domain.Tokens, error) {

	keys, err := s.getKeysPair(user)
	if err != nil {
		return nil, ErrorInvalidFormat
	}

	if err := s.repo.AddRecord(user, keys.refreshHash, keys.pairId); err != nil {
		s.log.Info("cant't add record to db", "message", err)
		return nil, ErrorInvalidFormat
	}

	user.AccessToken = keys.access
	user.RefreshToken = keys.refresh
	return &domain.Tokens{Access: keys.access,
		Refresh: keys.refresh,
	}, nil
}

func (s *authService) RefreshToken(inputUser *domain.User) (*domain.Tokens, error) {

	claims, err := s.token.Validate(*inputUser)
	if err != nil {
		s.log.Info("access token invalid", "message", err)
		return nil, ErrorInvalidTokens
	}
	if time.Now().After(time.Unix(int64(claims.ExpiredAt), 0)) {
		return nil, ErrorTokenExpired
	}

	registerUser, err := s.repo.GetRecord(uuid.MustParse(claims.UserId))
	if err != nil {
		s.log.Info("can't find user at data base", "message", err)
		return nil, ErrorInvalidTokens
	}

	if err := s.token.CompareBcryptTokens(
		registerUser.RefreshToken, inputUser.RefreshToken); err != nil {
		s.log.Info("refresh tokens not equal", "message", err)
		return nil, ErrorInvalidTokens
	}

	if registerUser.TokenPairId != uuid.MustParse(claims.PairId) {
		s.log.Info("access token reuse", "old", registerUser.AccessToken, "new", uuid.MustParse(claims.PairId))
		return nil, ErrorInvalidTokens
	}

	if !registerUser.IpAddress.Equal(inputUser.IpAddress) ||
		!inputUser.IpAddress.Equal(net.ParseIP(claims.IpAddress)) {
		//s.notify.SendMail(registerUser.Email,"вход с нвового ip")

		s.log.Info("вход с нового ip***", "mail", registerUser.Email)
	}

	keys, err := s.getKeysPair(registerUser)
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpdateRecord(registerUser, keys.refreshHash, keys.pairId); err != nil {
		s.log.Info("can't update record in data base", "message", err)
		return nil, err
	}

	return &domain.Tokens{Access: keys.access,
		Refresh: keys.refresh,
	}, nil

}

func (s *authService) getKeysPair(user *domain.User) (keys, error) {

	pairId := uuid.New()
	access, err := s.token.GenerateAccess(*user, pairId)
	if err != nil {
		s.log.Info("errror generate access key", "message", err)
		return keys{}, err
	}

	refresh := s.token.GenerateRefresh()

	refreshHash, err := s.token.GenerateBcryptToken(refresh)
	if err != nil {
		s.log.Info("errror hashing token", "message", err)
		return keys{}, err
	}

	return keys{
		access:      access,
		pairId:      pairId,
		refresh:     refresh,
		refreshHash: refreshHash,
	}, nil
}
