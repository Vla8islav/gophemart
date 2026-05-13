package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Vla8islav/gophemart/internal/domain"
	"github.com/Vla8islav/gophemart/internal/repository"
	"github.com/stretchr/testify/require"
)

func TestUserRegisterHandler_Success(t *testing.T) {
	service := fakeRegisterService{
		createUserFunc: func(ctx context.Context, req domain.UserRegisterRequest) (*domain.AuthResult, error) {
			require.Equal(t, "test-login", req.Login)
			require.Equal(t, "test-password", req.Password)

			return &domain.AuthResult{
				Token:  "test-token",
				UserID: 1,
			}, nil
		},
	}
	h := newTestRegisterHandler(service)

	body := bytes.NewBufferString(`{"login":"test-login","password":"test-password"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.UserRegisterHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)
}

func TestUserRegisterHandler_AllowsJSONContentTypeWithCharset(t *testing.T) {
	service := fakeRegisterService{
		createUserFunc: func(ctx context.Context, req domain.UserRegisterRequest) (*domain.AuthResult, error) {
			return &domain.AuthResult{
				Token:  "test-token",
				UserID: 1,
			}, nil
		},
	}
	h := newTestRegisterHandler(service)

	body := bytes.NewBufferString(`{"login":"test-login","password":"test-password"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", body)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	w := httptest.NewRecorder()

	h.UserRegisterHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)
}

func TestUserRegisterHandler_MethodNotAllowed(t *testing.T) {
	service := fakeRegisterService{
		createUserFunc: func(ctx context.Context, req domain.UserRegisterRequest) (*domain.AuthResult, error) {
			t.Fatal("CreateUser should not be called")
			return nil, nil
		},
	}
	h := newTestRegisterHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/api/user/register", nil)
	w := httptest.NewRecorder()

	h.UserRegisterHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
}

func TestUserRegisterHandler_BadContentType(t *testing.T) {
	service := fakeRegisterService{
		createUserFunc: func(ctx context.Context, req domain.UserRegisterRequest) (*domain.AuthResult, error) {
			t.Fatal("CreateUser should not be called")
			return nil, nil
		},
	}
	h := newTestRegisterHandler(service)

	body := bytes.NewBufferString(`{"login":"test-login","password":"test-password"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", body)
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()

	h.UserRegisterHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestUserRegisterHandler_InvalidJSON(t *testing.T) {
	service := fakeRegisterService{
		createUserFunc: func(ctx context.Context, req domain.UserRegisterRequest) (*domain.AuthResult, error) {
			t.Fatal("CreateUser should not be called")
			return nil, nil
		},
	}
	h := newTestRegisterHandler(service)

	body := bytes.NewBufferString(`{"login":`)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.UserRegisterHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestUserRegisterHandler_EmptyPassword(t *testing.T) {
	service := fakeRegisterService{
		createUserFunc: func(ctx context.Context, req domain.UserRegisterRequest) (*domain.AuthResult, error) {
			t.Fatal("CreateUser should not be called")
			return nil, nil
		},
	}
	h := newTestRegisterHandler(service)

	body := bytes.NewBufferString(`{"login":"test-login","password":""}`)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.UserRegisterHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestUserRegisterHandler_ServiceError(t *testing.T) {
	service := fakeRegisterService{
		createUserFunc: func(ctx context.Context, req domain.UserRegisterRequest) (*domain.AuthResult, error) {
			return nil, errors.New("service error")
		},
	}
	h := newTestRegisterHandler(service)

	body := bytes.NewBufferString(`{"login":"test-login","password":"test-password"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.UserRegisterHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestUserRegisterHandler_UserAlreadyExists(t *testing.T) {
	service := fakeRegisterService{
		createUserFunc: func(ctx context.Context, req domain.UserRegisterRequest) (*domain.AuthResult, error) {
			return nil, repository.ErrUserAlreadyExists
		},
	}
	h := newTestRegisterHandler(service)

	body := bytes.NewBufferString(`{"login":"test-login","password":"test-password"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.UserRegisterHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusConflict, res.StatusCode)
}
