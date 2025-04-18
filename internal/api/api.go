package api

import (
	"net/http"

	"github.com/aAmer0neee/authentication-service-test-task/internal/auth"
	"github.com/aAmer0neee/authentication-service-test-task/internal/config"
)

type AuthApi interface {
	LoginUser()
	RefreshToken()
}

func ConfigureApi(cfg config.Cfg, s auth.AuthService) *http.Server {

	return &http.Server{
		Handler: configureGin(s).router,
		Addr:    cfg.Server.Host + ":" + cfg.Server.Port}
}
