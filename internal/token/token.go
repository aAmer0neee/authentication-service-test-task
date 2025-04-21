package token

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/aAmer0neee/authentication-service-test-task/internal/domain"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

//go:generate mockgen -source=token.go -destination=mocks/token_mock.go -package=token_mock
type TokenService interface {
	GenerateAccess(user domain.User, pairId uuid.UUID) (string, error)
	Validate(user domain.User) (*AccessClaims, error)
	GenerateRefresh() string
	GenerateBcryptToken(token string) (string, error)
	CompareBcryptTokens(hash, token string) error
}

func New(secret string) TokenService {
	return configureJWT(secret)
}

type JWTService struct {
	secretKey []byte
}

type AccessClaims struct {
	UserId    string
	ExpiredAt float64
	IpAddress string
	PairId    string
}

func configureJWT(secretKey string) *JWTService {
	return &JWTService{secretKey: []byte(secretKey)}
}

func (j *JWTService) GenerateAccess(user domain.User, pairId uuid.UUID) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"pair_id":    pairId.String(),
		"user_id":    user.Id,
		"ip_address": user.IpAddress.String(),
		"exp":        time.Now().Add(time.Minute * 10).Unix(),
	}).SignedString(j.secretKey)
}

func (j *JWTService) Validate(user domain.User) (*AccessClaims, error) {
	token, err := jwt.Parse(user.AccessToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); ok {
			if t.Method.Alg() == jwt.SigningMethodHS512.Alg() {
				return j.secretKey, nil
			}
		}
		return nil, fmt.Errorf("non valid token")
	})

	if err != nil {
		return nil, err
	}

	return &AccessClaims{
		ExpiredAt: j.getExpTime(token),
		UserId:    j.getUserID(token),
		IpAddress: j.getIpAddress(token),
		PairId:    j.getPairID(token),
	}, nil
}

func (j *JWTService) GenerateRefresh() string {
	seed := make([]byte, 32)
	rand.Read(seed)

	return base64.StdEncoding.EncodeToString(seed)
}

func (j *JWTService) GenerateBcryptToken(token string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.MinCost)
	return string(hash), err
}

func (j *JWTService) CompareBcryptTokens(hash, token string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(token))
}

func (j *JWTService) getExpTime(token *jwt.Token) float64 {
	return token.Claims.(jwt.MapClaims)["exp"].(float64)
}

func (j *JWTService) getPairID(token *jwt.Token) string {
	return token.Claims.(jwt.MapClaims)["pair_id"].(string)
}

func (j *JWTService) getUserID(token *jwt.Token) string {
	return token.Claims.(jwt.MapClaims)["user_id"].(string)
	/* value, ok :=  token.Claims.(jwt.MapClaims)[key]
	if !ok {
		return ""
	}
	return value.(string) */
}

func (j *JWTService) getIpAddress(token *jwt.Token) string {
	return token.Claims.(jwt.MapClaims)["ip_address"].(string)
}
