package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/aAmer0neee/authentication-service-test-task/internal/domain"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type JWTService struct {
	secretKey []byte
}

func configureJWT(secretKey string) *JWTService {
	return &JWTService{secretKey: []byte(secretKey)}
}

func (j *JWTService) generateAccess(user domain.User) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"user_id":    user.Id,
		"ip_address": user.IpAddress,
		"exp":        time.Now().Add(time.Minute * 10).Unix(),
	}).SignedString(j.secretKey)
}

func (j *JWTService) validate(user domain.User) error {
	token, err := jwt.Parse(user.AccessToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); ok {
			if t.Method.Alg() == jwt.SigningMethodHS512.Alg() {
				return j.secretKey, nil
			}
		}
		return nil, fmt.Errorf("non valid token")
	})

	if err != nil {
		return err
	}
	if j.getExpTime(token) < time.Now().Unix() ||
		j.getUserID(token) != user.Id.String() {
		return fmt.Errorf("non valid token")
	}
	return nil
}

func (j *JWTService) getExpTime(token *jwt.Token) int64 {
	return token.Claims.(jwt.MapClaims)["exp"].(int64)
}

func (j *JWTService) getUserID(token *jwt.Token) string {
	return token.Claims.(jwt.MapClaims)["user_id"].(string)
	/* value, ok :=  token.Claims.(jwt.MapClaims)[key]
	if !ok {
		return ""
	}
	return value.(string) */
}

func (j *JWTService) generateRefresh() string {
	seed := make([]byte, 32)
	rand.Read(seed)

	return base64.StdEncoding.EncodeToString(seed)
}

func (j *JWTService) generateBcryptToken(token string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.MinCost)
	return string(hash), err
}

func (j *JWTService) compareBcryptTokens(hash, token string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(token))
}
