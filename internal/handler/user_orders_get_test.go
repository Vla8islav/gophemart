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

func TestUserOrdersGetHandler(t *testing.T) {
	t.Parallel()

	const userID int64 = 42
	uploadedAt := time.Date(2026, 5, 16, 10, 30, 0, 0,
		time.FixedZone("UTC+3", 3*60*60))
	accrual := int64(50000)
	accrualResponse := 500.0

	tests := []struct {
		name            string
		method          string
		withUserID      bool
		orders          []domain.UserOrder
		serviceErr      error
		wantStatusCode  int
		wantServiceCall bool
		wantResponse    domain.UserOrderGetResponse
	}{
		{
			name:       "success with accrual",
			method:     http.MethodGet,
			withUserID: true,
			orders: []domain.UserOrder{
				{
					Number:     "9278923470",
					Status:     "PROCESSED",
					Accrual:    &accrual,
					UploadedAt: uploadedAt,
				},
			},
			wantStatusCode:  http.StatusOK,
			wantServiceCall: true,
			wantResponse: domain.UserOrderGetResponse{
				{
					Number:     "9278923470",
					Status:     "PROCESSED",
					Accrual:    &accrualResponse,
					UploadedAt: uploadedAt,
				},
			},
		},
		{
			name:       "success without accrual",
			method:     http.MethodGet,
			withUserID: true,
			orders: []domain.UserOrder{
				{
					Number:     "12345678903",
					Status:     "PROCESSING",
					Accrual:    nil,
					UploadedAt: uploadedAt,
				},
			},
			wantStatusCode:  http.StatusOK,
			wantServiceCall: true,
			wantResponse: domain.UserOrderGetResponse{
				{
					Number:     "12345678903",
					Status:     "PROCESSING",
					Accrual:    nil,
					UploadedAt: uploadedAt,
				},
			},
		},
		{
			name:            "no orders",
			method:          http.MethodGet,
			withUserID:      true,
			orders:          []domain.UserOrder{},
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
					GetUserOrders(gomock.Any(), userID).
					Return(tt.orders, tt.serviceErr)
			}

			req := httptest.NewRequest(tt.method, "/api/user/orders", nil)
			if tt.withUserID {
				req = req.WithContext(middlewares.ContextWithUserID(req.Context(), userID))
			}
			rec := httptest.NewRecorder()

			h.UserOrdersGetHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			require.Equal(t, tt.wantStatusCode, res.StatusCode)

			if tt.wantStatusCode != http.StatusOK {
				return
			}

			require.Equal(t, "application/json", res.Header.Get("Content-Type"))

			var response domain.UserOrderGetResponse
			err := json.NewDecoder(res.Body).Decode(&response)
			require.NoError(t, err)

			require.Len(t, response, len(tt.wantResponse))
			for i := range tt.wantResponse {
				require.Equal(t, tt.wantResponse[i].Number, response[i].Number)
				require.Equal(t, tt.wantResponse[i].Status, response[i].Status)

				if tt.wantResponse[i].Accrual == nil {
					require.Nil(t, response[i].Accrual)
				} else {
					require.NotNil(t, response[i].Accrual)
					require.Equal(t, *tt.wantResponse[i].Accrual, *response[i].Accrual)
				}

				require.True(t, tt.wantResponse[i].UploadedAt.Equal(response[i].UploadedAt))
			}
		})
	}
}
