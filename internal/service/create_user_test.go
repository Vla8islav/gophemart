package service

import (
	"context"
	"testing"

	"github.com/Vla8islav/gophemart/internal/domain"
	"github.com/Vla8islav/gophemart/internal/helpers"
	"github.com/stretchr/testify/require"
)

func TestMetricsService_CreateUser(t *testing.T) {
	ctx := context.Background()

	userRegReq := domain.UserRegisterRequest{
		Login:    "test-login",
		Password: "test-password",
	}

	repository := fakeCreateUserRepository{
		createUserFunc: func(ctx context.Context, params domain.CreateUserParams) (int64, error) {
			require.Equal(t, userRegReq.Login, params.Login)
			require.NotEqual(t, userRegReq.Password, params.PasswordHash)
			require.NoError(t, helpers.CompareHashAndPassword(params.PasswordHash, userRegReq.Password))

			return 123, nil
		},
	}

	service := metricsService{
		repository: repository,
	}

	userID, err := service.CreateUser(ctx, userRegReq)
	require.NoError(t, err)
	require.Equal(t, int64(123), userID)
}
