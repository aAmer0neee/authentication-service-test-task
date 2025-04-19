package api

import (
	"fmt"
	"net"
	"net/http"

	"github.com/aAmer0neee/authentication-service-test-task/internal/auth"
	"github.com/aAmer0neee/authentication-service-test-task/internal/domain"
	"github.com/gin-gonic/gin"
)

type GinApi struct {
	router  *gin.Engine
	service auth.AuthService
}

func configureGin(s auth.AuthService) GinApi {

	api := GinApi{
		router:  gin.Default(),
		service: s,
	}

	api.configureHandlers()
	api.configureMiddleware()

	return api
}

func (r *GinApi) configureHandlers() {

	r.router.POST("/login", func(ctx *gin.Context) {
		r.handleLogin(ctx)
	})
	r.router.POST("/refresh", func(ctx *gin.Context) {
		r.handleRefresh(ctx)
	})
}

func (r *GinApi) configureMiddleware() {
	/* panic("not implemented") */
}

func (r *GinApi) handleLogin(ctx *gin.Context) {
	request := &domain.LoginRequest{}

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("missing required params"))
		return
	}

	ip := net.ParseIP(ctx.RemoteIP())
	if ip == nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	user := &domain.User{
		Id:        request.Id,
		Email:     request.Email,
		IpAddress: ip,
	}

	 tokens,err := r.service.LoginUser(user)
	 if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx.JSON(http.StatusOK,tokens)
}

func (r *GinApi) handleRefresh(ctx *gin.Context) {
	request := &domain.RefreshRequest{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("missing required params"))
		return
	}

	ip := net.ParseIP(ctx.RemoteIP())
	if ip == nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	tokens,err := r.service.RefreshToken(&domain.User{
		AccessToken:  request.AccessToken,
		RefreshToken: request.RefreshToken,
		IpAddress:    ip,
	})
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx.JSON(http.StatusOK, tokens)
}
