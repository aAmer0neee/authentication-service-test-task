package api

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	auth_mock "github.com/aAmer0neee/authentication-service-test-task/internal/auth/mocks"
	"github.com/aAmer0neee/authentication-service-test-task/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestApi_login(t *testing.T) {
	gin.DefaultWriter = io.Discard
	tests := []struct {
		name        string
		requestBody string
		Ip          net.IP
		statusCode  int
		expectCall  func(m *auth_mock.MockAuthService, user *domain.User)
		callReturn  domain.Tokens
		callUser    domain.User
	}{
		{
			name: "OK",
			requestBody: `{
				"id": "0f5ae05b-4d6e-0c0e-43f6-71deb028c0a3",
				"email": "mail@mail.com"
			}`,
			statusCode: http.StatusOK,
			Ip:         net.ParseIP("127.0.0.1"),
			expectCall: func(m *auth_mock.MockAuthService, user *domain.User) {
				m.EXPECT().LoginUser(user).Return(&domain.Tokens{
					Access:  "access",
					Refresh: "refresh",
				}, nil)
			},
			callReturn: domain.Tokens{
				Access:  "access",
				Refresh: "refresh",
			},
			callUser: domain.User{
				Id:        uuid.MustParse("0f5ae05b-4d6e-0c0e-43f6-71deb028c0a3"),
				Email:     "mail@mail.com",
				IpAddress: net.ParseIP("127.0.0.1"),
			},
		},
		{
			name: "MISSING BODY",
			requestBody: `{
			}`,
			statusCode: http.StatusBadRequest,
			Ip:         net.ParseIP("127.0.0.1"),
			expectCall: func(m *auth_mock.MockAuthService, user *domain.User) {
				m.EXPECT().
					LoginUser(gomock.Any()).
					Times(0)
			},
		},
		{
			name: "MISSING IP",
			requestBody: `{
				"id": "0f5ae05b-4d6e-0c0e-43f6-71deb028c0a3",
				"email": "mail@mail.com"
			}`,
			statusCode: http.StatusBadRequest,
			expectCall: func(m *auth_mock.MockAuthService, user *domain.User) {
				m.EXPECT().
					LoginUser(gomock.Any()).
					Times(0)
			},
		},
		{
			name: "BAD UUID",
			requestBody: `{
				"id": "0f5ae05b-4d6",
				"email": "mail@mail.com"
			}`,
			statusCode: http.StatusBadRequest,
			expectCall: func(m *auth_mock.MockAuthService, user *domain.User) {
				m.EXPECT().
					LoginUser(gomock.Any()).
					Times(0)
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			control := gomock.NewController(t)
			defer control.Finish()

			auth := auth_mock.NewMockAuthService(control)
			api := NewGin(auth)

			testCase.expectCall(auth, &testCase.callUser)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/login",
				bytes.NewBufferString(testCase.requestBody))

			req.RemoteAddr = testCase.Ip.String() + ":8000"
			req.Header.Set("Content-Type", "application/json")

			api.router.ServeHTTP(w, req)

			assert.Equal(t, testCase.statusCode, w.Result().StatusCode)
		})
	}
}

func TestApi_refresh(t *testing.T) {
	gin.DefaultWriter = io.Discard
	tests := []struct {
		name        string
		requestBody string
		ip          net.IP
		statusCode  int
		expectCall  func(m *auth_mock.MockAuthService, user *domain.User)
		callReturn  *domain.Tokens
		callUser    *domain.User
	}{
		{
			name: "OK",
			ip:   net.ParseIP("127.0.0.1"),
			requestBody: `{
    "Access": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDUxNzg1MDYsImlwX2FkZHJlc3MiOiIxMjcuMC4wLjEiLCJwYWlyX2lkIjoiMzdiZDEzZGQtNDUwNS00MzQ0LTlkYmMtZTE4MGNlNmVkNGJlIiwidXNlcl9pZCI6IjBmNWFlMDViLTRkN2UtMGMwZS00M2Y2LTcxZGViMDI4YzBhMyJ9.suvwHpkyaAGixJe7lCpQt-O6fzFlBvB8eMJ1XfVveq2frOr8v3QuYFBnTHVDI6KcueTEcbQxNn4x_M839UZbyg",
    "Refresh": "1M8QNRcJYYkiPMHruPlqC8N+gISvxGZs7yE78xBVndc="
	}`,
			statusCode: http.StatusOK,
			expectCall: func(m *auth_mock.MockAuthService, user *domain.User) {
				m.EXPECT().RefreshToken(gomock.Any()).Return(&domain.Tokens{
					Access:  "access",
					Refresh: "refresh",
				}, nil)
			},
			callReturn: &domain.Tokens{
				Access:  "access",
				Refresh: "refresh",
			},
			callUser: &domain.User{
				AccessToken:  "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDUxNzg1MDYsImlwX2FkZHJlc3MiOiIxMjcuMC4wLjEiLCJwYWlyX2lkIjoiMzdiZDEzZGQtNDUwNS00MzQ0LTlkYmMtZTE4MGNlNmVkNGJlIiwidXNlcl9pZCI6IjBmNWFlMDViLTRkN2UtMGMwZS00M2Y2LTcxZGViMDI4YzBhMyJ9.suvwHpkyaAGixJe7lCpQt-O6fzFlBvB8eMJ1XfVveq2frOr8v3QuYFBnTHVDI6KcueTEcbQxNn4x_M839UZbyg",
				RefreshToken: "1M8QNRcJYYkiPMHruPlqC8N+gISvxGZs7yE78xBVndc=",
				IpAddress:    net.ParseIP("127.0.0.1"),
			},
		},
		{
			name: "MISSING BODY",
			requestBody: `{
			}`,
			statusCode: http.StatusBadRequest,
			ip:         net.ParseIP("127.0.0.1"),
			expectCall: func(m *auth_mock.MockAuthService, user *domain.User) {
				m.EXPECT().
					LoginUser(gomock.Any()).
					Times(0)
			},
		},
		{
			name: "MISSING IP",
			requestBody: `{
				"id": "0f5ae05b-4d6e-0c0e-43f6-71deb028c0a3",
				"email": "mail@mail.com"
			}`,
			statusCode: http.StatusBadRequest,
			expectCall: func(m *auth_mock.MockAuthService, user *domain.User) {
				m.EXPECT().
					LoginUser(gomock.Any()).
					Times(0)
			},
		},
		{
			name: "BAD UUID",
			requestBody: `{
				"id": "0f5ae05b-4d6",
				"email": "mail@mail.com"
			}`,
			statusCode: http.StatusBadRequest,
			expectCall: func(m *auth_mock.MockAuthService, user *domain.User) {
				m.EXPECT().
					LoginUser(gomock.Any()).
					Times(0)
			},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			control := gomock.NewController(t)
			defer control.Finish()

			auth := auth_mock.NewMockAuthService(control)

			api := NewGin(auth)
			testCase.expectCall(auth, testCase.callUser)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/refresh",
				bytes.NewBufferString(testCase.requestBody))

			req.RemoteAddr = testCase.ip.String() + ":8000"
			req.Header.Set("Content-Type", "application/json")

			api.router.ServeHTTP(w, req)

			assert.Equal(t, testCase.statusCode, w.Result().StatusCode)
		})
	}
}
