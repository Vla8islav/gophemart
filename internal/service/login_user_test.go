package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Vla8islav/gophemart/internal/domain"
	"github.com/Vla8islav/gophemart/internal/helpers"
	"github.com/stretchr/testify/require"
)

type fakeLoginUserRepository struct {
	getUserByLoginFunc func(ctx context.Context, login string) (*domain.User, error)
}

func (r fakeLoginUserRepository) Ping(ctx context.Context) error {
	return nil
}

func (r fakeLoginUserRepository) CreateUser(ctx context.Context, params domain.CreateUserParams) (int64, error) {
	return 0, nil
}

func (r fakeLoginUserRepository) GetUserByLogin(ctx context.Context, login string) (*domain.User, error) {
	return r.getUserByLoginFunc(ctx, login)
}

func TestMetricsService_LoginUser(t *testing.T) {
	ctx := context.Background()

	passwordHash, err := helpers.HashPassword("test-password")
	require.NoError(t, err)

	repository := fakeLoginUserRepository{
		getUserByLoginFunc: func(ctx context.Context, login string) (*domain.User, error) {
			require.Equal(t, "test-login", login)

			return &domain.User{
				ID:           123,
				Login:        "test-login",
				PasswordHash: passwordHash,
			}, nil
		},
	}

	service := metricsService{
		repository: repository,
	}

	authResult, err := service.LoginUser(ctx, domain.UserLoginRequest{
		Login:    "test-login",
		Password: "test-password",
	})
	require.NoError(t, err)
	require.NotNil(t, authResult)
	require.Equal(t, int64(123), authResult.UserID)
	require.NotEmpty(t, authResult.Token)
}

func TestMetricsService_LoginUser_InvalidPassword(t *testing.T) {
	ctx := context.Background()

	passwordHash, err := helpers.HashPassword("test-password")
	require.NoError(t, err)

	repository := fakeLoginUserRepository{
		getUserByLoginFunc: func(ctx context.Context, login string) (*domain.User, error) {
			return &domain.User{
				ID:           123,
				Login:        "test-login",
				PasswordHash: passwordHash,
			}, nil
		},
	}

	service := metricsService{
		repository: repository,
	}

	authResult, err := service.LoginUser(ctx, domain.UserLoginRequest{
		Login:    "test-login",
		Password: "wrong-password",
	})
	require.ErrorIs(t, err, ErrInvalidUserCredentials)
	require.Nil(t, authResult)
}

func TestMetricsService_LoginUser_UserNotFound(t *testing.T) {
	ctx := context.Background()

	repository := fakeLoginUserRepository{
		getUserByLoginFunc: func(ctx context.Context, login string) (*domain.User, error) {
			return nil, ErrInvalidUserCredentials
		},
	}

	service := metricsService{
		repository: repository,
	}

	authResult, err := service.LoginUser(ctx, domain.UserLoginRequest{
		Login:    "missing-login",
		Password: "test-password",
	})
	require.ErrorIs(t, err, ErrInvalidUserCredentials)
	require.Nil(t, authResult)
}

func TestMetricsService_LoginUser_RepositoryError(t *testing.T) {
	ctx := context.Background()
	repositoryErr := errors.New("repository error")

	repository := fakeLoginUserRepository{
		getUserByLoginFunc: func(ctx context.Context, login string) (*domain.User, error) {
			return nil, repositoryErr
		},
	}

	service := metricsService{
		repository: repository,
	}

	authResult, err := service.LoginUser(ctx, domain.UserLoginRequest{
		Login:    "test-login",
		Password: "test-password",
	})
	require.ErrorIs(t, err, repositoryErr)
	require.Nil(t, authResult)
}
