package auth

import (
	"errors"
	"io"
	"log/slog"
	"net"
	"testing"
	"time"

	"github.com/aAmer0neee/authentication-service-test-task/internal/domain"
	notify_mock "github.com/aAmer0neee/authentication-service-test-task/internal/notify/mocks"
	repository_mock "github.com/aAmer0neee/authentication-service-test-task/internal/repository/mocks"
	"github.com/aAmer0neee/authentication-service-test-task/internal/token"
	token_mock "github.com/aAmer0neee/authentication-service-test-task/internal/token/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuth_login(t *testing.T) {
	tests := []struct {
		name       string
		input      *domain.User
		returnVal  *domain.Tokens
		expectErr  error
		expectCall func(m *token_mock.MockTokenService, m2 *repository_mock.MockRepository)
	}{
		{
			name: "OK",
			input: &domain.User{
				Id: uuid.MustParse("0f5ae05b-4d6e-0c0e-43f6-71deb028c0a3"),
			},
			expectErr: nil,
			returnVal: &domain.Tokens{
				Access:  "token",
				Refresh: "token",
			},

			expectCall: func(m *token_mock.MockTokenService, m2 *repository_mock.MockRepository) {
				m.EXPECT().
					GenerateAccess(gomock.Any(), gomock.Any()).
					Return("token", nil)

				m.EXPECT().
					GenerateRefresh().
					Return("token")

				m.EXPECT().
					GenerateBcryptToken(gomock.Any()).
					Return("token", nil)

				m2.EXPECT().AddRecord(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)

			},
		},
		{
			name: "BAD REPOSITORY ADDING",
			input: &domain.User{
				Id: uuid.MustParse("0f5ae05b-4d6e-0c0e-43f6-71deb028c0a3"),
			},
			expectErr: ErrorInvalidFormat,
			returnVal: nil,

			expectCall: func(m *token_mock.MockTokenService, m2 *repository_mock.MockRepository) {
				m.EXPECT().
					GenerateAccess(gomock.Any(), gomock.Any()).
					Return("token", nil)

				m.EXPECT().
					GenerateRefresh().
					Return("token")

				m.EXPECT().
					GenerateBcryptToken(gomock.Any()).
					Return("token", nil)

				m2.EXPECT().AddRecord(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New(""))

			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			control := gomock.NewController(t)
			defer control.Finish()

			repo := repository_mock.NewMockRepository(control)
			token := token_mock.NewMockTokenService(control)
			notify := notify_mock.NewMockNotifyer(control)
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))

			auth := New(repo, token, notify, logger)

			testCase.expectCall(token, repo)
			tokens, err := auth.LoginUser(testCase.input)

			assert.Equal(t, testCase.returnVal, tokens)
			assert.Equal(t, testCase.expectErr, err)
			assert.Equal(t, testCase.returnVal, tokens)
		})
	}
}

func TestAuth_refresh(t *testing.T) {
	tests := []struct {
		name       string
		input      *domain.User
		expectErr  error
		returnVal  *domain.Tokens
		expectCall func(m1 *token_mock.MockTokenService, m2 *repository_mock.MockRepository)
	}{
		{
			name: "OK",
			input: &domain.User{
				AccessToken:  "token",
				RefreshToken: "token",
			},
			expectErr: nil,
			returnVal: &domain.Tokens{
				Access:  "token",
				Refresh: "token",
			},
			expectCall: func(m *token_mock.MockTokenService, m2 *repository_mock.MockRepository) {

				m.EXPECT().
					Validate(gomock.Any()).
					Return(&token.AccessClaims{
						UserId:    "0f5ae05b-4d6e-0c0e-43f6-71deb028c0a3",
						PairId:    "0f5ae05b-4d6e-0c0e-43f6-71deb028c0a1",
						ExpiredAt: float64(time.Now().Add(time.Minute).Unix()),
						IpAddress: "127.0.0.1",
					}, nil)
				m2.EXPECT().
					GetRecord(gomock.Any()).Return(&domain.User{
					Id:           uuid.MustParse("0f5ae05b-4d6e-0c0e-43f6-71deb028c0a3"),
					Email:        "mail@mail.com",
					IpAddress:    net.ParseIP("127.0.0.1"),
					AccessToken:  "token",
					RefreshToken: "token",
					TokenPairId:  uuid.MustParse("0f5ae05b-4d6e-0c0e-43f6-71deb028c0a1"),
				}, nil)
				m.EXPECT().
					CompareBcryptTokens(gomock.Any(), gomock.Any()).
					Return(nil)
				m.EXPECT().
					GenerateAccess(gomock.Any(), gomock.Any()).
					Return("token", nil)
				m.EXPECT().
					GenerateRefresh().
					Return("token")

				m.EXPECT().
					GenerateBcryptToken(gomock.Any()).
					Return("token", nil)

				m2.EXPECT().UpdateRecord(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
		},
		{
			name: "TOKEN EXPIRED",
			input: &domain.User{
				AccessToken:  "token",
				RefreshToken: "token",
			},
			expectErr: ErrorTokenExpired,
			returnVal: nil,
			expectCall: func(m *token_mock.MockTokenService, m2 *repository_mock.MockRepository) {

				m.EXPECT().
					Validate(gomock.Any()).
					Return(&token.AccessClaims{
						UserId:    "0f5ae05b-4d6e-0c0e-43f6-71deb028c0a3",
						PairId:    "0f5ae05b-4d6e-0c0e-43f6-71deb028c0a1",
						ExpiredAt: float64(time.Now().Unix()),
						IpAddress: "127.0.0.1",
					}, nil)
				m2.EXPECT().
					GetRecord(gomock.Any()).
					Times(0)
				m.EXPECT().
					CompareBcryptTokens(gomock.Any(), gomock.Any()).
					Times(0)
				m.EXPECT().
					GenerateAccess(gomock.Any(), gomock.Any()).
					Times(0)
				m.EXPECT().
					GenerateRefresh().
					Times(0)

				m.EXPECT().
					GenerateBcryptToken(gomock.Any()).
					Times(0)

				m2.EXPECT().UpdateRecord(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			control := gomock.NewController(t)
			defer control.Finish()

			repo := repository_mock.NewMockRepository(control)
			token := token_mock.NewMockTokenService(control)
			notify := notify_mock.NewMockNotifyer(control)
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))

			auth := New(repo, token, notify, logger)

			testCase.expectCall(token, repo)
			tokens, err := auth.RefreshToken(testCase.input)

			assert.Equal(t, testCase.returnVal, tokens)
			assert.Equal(t, testCase.expectErr, err)
			assert.Equal(t, testCase.returnVal, tokens)
		})
	}
}
