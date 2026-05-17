package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Vla8islav/gophemart/internal/domain"
	"github.com/Vla8islav/gophemart/internal/middlewares"
	"github.com/Vla8islav/gophemart/internal/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestUserBalanceHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		method         string
		withUserID     bool
		userID         int64
		balance        *domain.UserBalance
		serviceErr     error
		wantStatusCode int
		wantBody       string
	}{
		{
			name:       "success",
			method:     http.MethodGet,
			withUserID: true,
			userID:     42,
			balance: &domain.UserBalance{
				Current:   212452,
				Withdrawn: 1000,
			},
			wantStatusCode: http.StatusOK,
			wantBody:       `{"current":2124.52,"withdrawn":10}`,
		},
		{
			name:           "method not allowed",
			method:         http.MethodPost,
			wantStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:           "unauthorized without user id in context",
			method:         http.MethodGet,
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name:           "service error",
			method:         http.MethodGet,
			withUserID:     true,
			userID:         42,
			serviceErr:     errors.New("db is unavailable"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			service := mocks.NewMockGophemartService(ctrl)
			h := &Handler{
				service: service,
				logger:  zap.NewNop(),
			}

			if tt.withUserID && tt.method == http.MethodGet {
				service.EXPECT().
					GetUserBalance(gomock.Any(), tt.userID).
					Return(tt.balance, tt.serviceErr)
			}

			req := httptest.NewRequest(tt.method, "/api/user/balance", nil)
			if tt.withUserID {
				req = req.WithContext(middlewares.ContextWithUserID(req.Context(), tt.userID))
			}
			rec := httptest.NewRecorder()

			h.UserBalanceHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			require.Equal(t, tt.wantStatusCode, res.StatusCode)

			if tt.wantBody != "" {
				require.JSONEq(t, tt.wantBody, rec.Body.String())
			}
		})
	}
}
