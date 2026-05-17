package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Vla8islav/gophemart/internal/domain"
	"github.com/Vla8islav/gophemart/internal/middlewares"
	"github.com/Vla8islav/gophemart/internal/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func newTestUserWithdrawalsHandler(service domain.GophemartService) *Handler {
	return &Handler{
		service: service,
		logger:  zap.NewNop(),
	}
}

func TestUserWithdrawalsHandler(t *testing.T) {
	t.Parallel()

	const userID int64 = 42
	processedAt := time.Date(2026, 5, 16, 10, 30, 0, 0, time.FixedZone("UTC+3", 3*60*60))

	tests := []struct {
		name            string
		method          string
		withUserID      bool
		withdrawals     []domain.Withdrawal
		serviceErr      error
		wantStatusCode  int
		wantServiceCall bool
		wantResponse    domain.WithdrawResponse
	}{
		{
			name:       "success",
			method:     http.MethodGet,
			withUserID: true,
			withdrawals: []domain.Withdrawal{
				{
					OrderNumber: "2377225624",
					Amount:      50000,
					ProcessedAt: processedAt,
				},
			},
			wantStatusCode:  http.StatusOK,
			wantServiceCall: true,
			wantResponse: domain.WithdrawResponse{
				{
					Order:       "2377225624",
					Sum:         500,
					ProcessedAt: processedAt,
				},
			},
		},
		{
			name:            "no withdrawals",
			method:          http.MethodGet,
			withUserID:      true,
			withdrawals:     []domain.Withdrawal{},
			wantStatusCode:  http.StatusNoContent,
			wantServiceCall: true,
		},
		{
			name:            "service internal error",
			method:          http.MethodGet,
			withUserID:      true,
			serviceErr:      errors.New("database is unavailable"),
			wantStatusCode:  http.StatusInternalServerError,
			wantServiceCall: true,
		},
		{
			name:           "method not allowed",
			method:         http.MethodPost,
			withUserID:     true,
			wantStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:           "unauthorized without user id in context",
			method:         http.MethodGet,
			wantStatusCode: http.StatusUnauthorized,
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

			if tt.wantServiceCall {
				service.EXPECT().
					GetUserWithdrawals(gomock.Any(), userID).
					Return(tt.withdrawals, tt.serviceErr)
			}

			req := httptest.NewRequest(tt.method, "/api/user/withdrawals", nil)
			if tt.withUserID {
				req = req.WithContext(middlewares.ContextWithUserID(req.Context(), userID))
			}
			rec := httptest.NewRecorder()

			h.UserWithdrawalsHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			require.Equal(t, tt.wantStatusCode, res.StatusCode)

			if tt.wantStatusCode != http.StatusOK {
				return
			}

			require.Equal(t, "application/json", res.Header.Get("Content-Type"))

			var response domain.WithdrawResponse
			err := json.NewDecoder(res.Body).Decode(&response)
			require.NoError(t, err)

			require.Len(t, response, len(tt.wantResponse))

			for i := range tt.wantResponse {
				require.Equal(t, tt.wantResponse[i].Order, response[i].Order)
				require.Equal(t, tt.wantResponse[i].Sum, response[i].Sum)
				require.True(t, tt.wantResponse[i].ProcessedAt.Equal(response[i].ProcessedAt))

			}
		})
	}
}
