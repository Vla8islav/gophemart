package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Vla8islav/gophemart/internal/domain"
	"github.com/Vla8islav/gophemart/internal/middlewares"
	"github.com/Vla8islav/gophemart/internal/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestUserBalanceWithdrawHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		method          string
		contentType     string
		body            string
		withUserID      bool
		userID          int64
		wantServiceCall bool
		wantRequest     domain.UserBalanceWithdraw
		serviceErr      error
		wantStatusCode  int
	}{
		{
			name:            "success",
			method:          http.MethodPost,
			contentType:     "application/json",
			body:            `{"order":"2377225624","sum":751}`,
			withUserID:      true,
			userID:          42,
			wantServiceCall: true,
			wantRequest: domain.UserBalanceWithdraw{
				Order: "2377225624",
				Sum:   75100,
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:            "success with fractional sum",
			method:          http.MethodPost,
			contentType:     "application/json",
			body:            `{"order":"2377225624","sum":751.55}`,
			withUserID:      true,
			userID:          42,
			wantServiceCall: true,
			wantRequest: domain.UserBalanceWithdraw{
				Order: "2377225624",
				Sum:   75155,
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "method not allowed",
			method:         http.MethodGet,
			contentType:    "application/json",
			body:           `{"order":"2377225624","sum":751}`,
			withUserID:     true,
			userID:         42,
			wantStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:           "unauthorized without user id in context",
			method:         http.MethodPost,
			contentType:    "application/json",
			body:           `{"order":"2377225624","sum":751}`,
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name:           "missing content type",
			method:         http.MethodPost,
			body:           `{"order":"2377225624","sum":751}`,
			withUserID:     true,
			userID:         42,
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:           "unsupported content type",
			method:         http.MethodPost,
			contentType:    "text/plain",
			body:           `{"order":"2377225624","sum":751}`,
			withUserID:     true,
			userID:         42,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "invalid json",
			method:         http.MethodPost,
			contentType:    "application/json",
			body:           `{invalid-json`,
			withUserID:     true,
			userID:         42,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "empty order",
			method:         http.MethodPost,
			contentType:    "application/json",
			body:           `{"order":"","sum":751}`,
			withUserID:     true,
			userID:         42,
			wantStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:           "zero sum",
			method:         http.MethodPost,
			contentType:    "application/json",
			body:           `{"order":"2377225624","sum":0}`,
			withUserID:     true,
			userID:         42,
			wantStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:           "tiny positive sum rounds to zero cents",
			method:         http.MethodPost,
			contentType:    "application/json",
			body:           `{"order":"2377225624","sum":0.001}`,
			withUserID:     true,
			userID:         42,
			wantStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:            "not enough money",
			method:          http.MethodPost,
			contentType:     "application/json",
			body:            `{"order":"2377225624","sum":751}`,
			withUserID:      true,
			userID:          42,
			wantServiceCall: true,
			wantRequest: domain.UserBalanceWithdraw{
				Order: "2377225624",
				Sum:   75100,
			},
			serviceErr:     domain.ErrNotEnoughMoney,
			wantStatusCode: http.StatusPaymentRequired,
		},
		{
			name:            "invalid order number from service",
			method:          http.MethodPost,
			contentType:     "application/json",
			body:            `{"order":"1234567890","sum":751}`,
			withUserID:      true,
			userID:          42,
			wantServiceCall: true,
			wantRequest: domain.UserBalanceWithdraw{
				Order: "1234567890",
				Sum:   75100,
			},
			serviceErr:     domain.ErrInvalidOrderNumber,
			wantStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:            "service error",
			method:          http.MethodPost,
			contentType:     "application/json",
			body:            `{"order":"2377225624","sum":751}`,
			withUserID:      true,
			userID:          42,
			wantServiceCall: true,
			wantRequest: domain.UserBalanceWithdraw{
				Order: "2377225624",
				Sum:   75100,
			},
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

			if tt.wantServiceCall {
				service.EXPECT().
					WithdrawFromUserBalance(gomock.Any(), tt.userID, tt.wantRequest).
					Return(tt.serviceErr)
			}

			req := httptest.NewRequest(tt.method, "/api/user/balance/withdraw", strings.NewReader(tt.body))
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}
			if tt.withUserID {
				req = req.WithContext(middlewares.ContextWithUserID(req.Context(), tt.userID))
			}

			rec := httptest.NewRecorder()

			h.UserBalanceWithdrawHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			require.Equal(t, tt.wantStatusCode, res.StatusCode)
		})
	}
}
