package api

import (
	"errors"
	"net"
	"net/http"

	"github.com/aAmer0neee/authentication-service-test-task/internal/auth"
	"github.com/aAmer0neee/authentication-service-test-task/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GinApi struct {
	router  *gin.Engine
	service auth.AuthService
}

func NewGin(s auth.AuthService) *GinApi {

	api := &GinApi{
		router:  gin.Default(),
		service: s,
	}

	api.configureHandlers()

	return api
}

func (r *GinApi) configureHandlers() {

	r.router.POST("/login", r.LoginUser)
	r.router.POST("/refresh", r.RefreshToken)
}

type RefreshRequest struct {
	AccessToken  string `json:"access" binding:"required"`
	RefreshToken string `json:"refresh" binding:"required"`
}

type LoginRequest struct {
	Id    uuid.UUID `json:"id" binding:"required"`
	Email string    `json:"email"`
}

func (r *GinApi) LoginUser(ctx *gin.Context) {

	request := &LoginRequest{}

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ip := net.ParseIP(ctx.RemoteIP())
	if ip == nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	tokens, err := r.service.LoginUser(&domain.User{
		Id:        request.Id,
		Email:     request.Email,
		IpAddress: ip,
	})
	if err != nil {
		unwrapErr(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, tokens)
}

func (r *GinApi) RefreshToken(ctx *gin.Context) {
	request := &RefreshRequest{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ip := net.ParseIP(ctx.RemoteIP())
	if ip == nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	tokens, err := r.service.RefreshToken(&domain.User{
		AccessToken:  request.AccessToken,
		RefreshToken: request.RefreshToken,
		IpAddress:    ip,
	})
	if err != nil {
		unwrapErr(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, tokens)
}

func unwrapErr(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, auth.ErrorInvalidFormat):
		ctx.AbortWithError(http.StatusBadRequest, err)
	case errors.Is(err, auth.ErrorInvalidTokens):
		ctx.AbortWithError(http.StatusUnauthorized, err)
	case errors.Is(err, auth.ErrorTokenExpired):
		ctx.AbortWithError(http.StatusUnauthorized, err)
	default:
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}

}
