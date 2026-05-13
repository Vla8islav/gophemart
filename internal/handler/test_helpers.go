package handler

import (
	"context"

	"github.com/Vla8islav/gophemart/internal/domain"
	"go.uber.org/zap"
)

type fakeRegisterService struct {
	createUserFunc func(ctx context.Context, req domain.UserRegisterRequest) (*domain.AuthResult, error)
	loginUserFunc  func(ctx context.Context, req domain.UserLoginRequest) (*domain.AuthResult, error)
}

func (s fakeRegisterService) Ping(ctx context.Context) error {
	return nil
}

func (s fakeRegisterService) CreateUser(ctx context.Context, req domain.UserRegisterRequest) (*domain.AuthResult, error) {
	return s.createUserFunc(ctx, req)
}

func (s fakeRegisterService) LoginUser(ctx context.Context, req domain.UserLoginRequest) (*domain.AuthResult, error) {
	return s.loginUserFunc(ctx, req)
}

func newTestRegisterHandler(service fakeRegisterService) *Handler {
	return &Handler{
		service: service,
		logger:  zap.NewNop(),
	}
}
