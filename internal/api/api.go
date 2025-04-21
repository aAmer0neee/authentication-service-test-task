package api

import (
	"net/http"

	"github.com/aAmer0neee/authentication-service-test-task/internal/auth"
	"github.com/aAmer0neee/authentication-service-test-task/internal/config"
	"github.com/gin-gonic/gin"
)

//go:generate mockgen -source=api.go -destination=mocks/api_mock.go -package=api_mock
type AuthApi interface {
	LoginUser(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
}

func ConfigureApi(cfg config.Cfg, s auth.AuthService) *http.Server {
	return &http.Server{
		Handler: NewGin(s).router,
		Addr:    cfg.Server.Host + ":" + cfg.Server.Port}
}
