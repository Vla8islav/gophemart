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

func TestUserOrdersPostHandler(t *testing.T) {
	t.Parallel()

	const userID int64 = 42

	tests := []struct {
		name            string
		method          string
		contentType     string
		body            string
		withUserID      bool
		serviceErr      error
		wantStatusCode  int
		wantServiceCall bool
		wantOrderNumber string
	}{
		{
			name:            "new order accepted",
			method:          http.MethodPost,
			contentType:     "text/plain",
			body:            "79927398713",
			withUserID:      true,
			wantStatusCode:  http.StatusAccepted,
			wantServiceCall: true,
			wantOrderNumber: "79927398713",
		},
		{
			name:            "new order accepted with charset content type",
			method:          http.MethodPost,
			contentType:     "text/plain; charset=utf-8",
			body:            "79927398713",
			withUserID:      true,
			wantStatusCode:  http.StatusAccepted,
			wantServiceCall: true,
			wantOrderNumber: "79927398713",
		},
		{
			name:            "trims order number before service call",
			method:          http.MethodPost,
			contentType:     "text/plain",
			body:            "\n 79927398713 \t",
			withUserID:      true,
			wantStatusCode:  http.StatusAccepted,
			wantServiceCall: true,
			wantOrderNumber: "79927398713",
		},
		{
			name:            "order already uploaded by same user",
			method:          http.MethodPost,
			contentType:     "text/plain",
			body:            "79927398713",
			withUserID:      true,
			serviceErr:      domain.ErrOrderAlreadyUploadedByUser,
			wantStatusCode:  http.StatusOK,
			wantServiceCall: true,
			wantOrderNumber: "79927398713",
		},
		{
			name:            "order already uploaded by another user",
			method:          http.MethodPost,
			contentType:     "text/plain",
			body:            "79927398713",
			withUserID:      true,
			serviceErr:      domain.ErrOrderAlreadyUploadedByAnotherUser,
			wantStatusCode:  http.StatusConflict,
			wantServiceCall: true,
			wantOrderNumber: "79927398713",
		},
		{
			name:            "invalid order number",
			method:          http.MethodPost,
			contentType:     "text/plain",
			body:            "79927398710",
			withUserID:      true,
			serviceErr:      domain.ErrInvalidOrderNumber,
			wantStatusCode:  http.StatusUnprocessableEntity,
			wantServiceCall: true,
			wantOrderNumber: "79927398710",
		},
		{
			name:            "empty order number",
			method:          http.MethodPost,
			contentType:     "text/plain",
			body:            "   ",
			withUserID:      true,
			serviceErr:      domain.ErrOrderEmpty,
			wantStatusCode:  http.StatusBadRequest,
			wantServiceCall: true,
			wantOrderNumber: "",
		},
		{
			name:            "service internal error",
			method:          http.MethodPost,
			contentType:     "text/plain",
			body:            "79927398713",
			withUserID:      true,
			serviceErr:      errors.New("database is unavailable"),
			wantStatusCode:  http.StatusInternalServerError,
			wantServiceCall: true,
			wantOrderNumber: "79927398713",
		},
		{
			name:           "method not allowed",
			method:         http.MethodGet,
			contentType:    "text/plain",
			body:           "79927398713",
			withUserID:     true,
			wantStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:           "bad content type",
			method:         http.MethodPost,
			contentType:    "application/json",
			body:           "79927398713",
			withUserID:     true,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "missing content type",
			method:         http.MethodPost,
			body:           "79927398713",
			withUserID:     true,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "unauthorized without user id in context",
			method:         http.MethodPost,
			contentType:    "text/plain",
			body:           "79927398713",
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
					CreateOrder(gomock.Any(), userID, tt.wantOrderNumber).
					Return(tt.serviceErr)
			}

			req := httptest.NewRequest(tt.method, "/api/user/orders", strings.NewReader(tt.body))
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}
			if tt.withUserID {
				req = req.WithContext(middlewares.ContextWithUserID(req.Context(), userID))
			}
			rec := httptest.NewRecorder()

			h.UserOrdersPostHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			require.Equal(t, tt.wantStatusCode, res.StatusCode)
		})
	}
}
