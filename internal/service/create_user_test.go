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

func TestMetricsService_CreateUser_HashesSamePasswordDifferently(t *testing.T) {
	ctx := context.Background()

	firstReq := domain.UserRegisterRequest{
		Login:    "first-user",
		Password: "same-password",
	}
	secondReq := domain.UserRegisterRequest{
		Login:    "second-user",
		Password: "same-password",
	}

	var passwordHashes []string

	repository := fakeCreateUserRepository{
		createUserFunc: func(ctx context.Context, params domain.CreateUserParams) (int64, error) {
			passwordHashes = append(passwordHashes, params.PasswordHash)
			return int64(len(passwordHashes)), nil
		},
	}

	service := metricsService{
		repository: repository,
	}

	firstUserID, err := service.CreateUser(ctx, firstReq)
	require.NoError(t, err)
	require.Equal(t, int64(1), firstUserID)

	secondUserID, err := service.CreateUser(ctx, secondReq)
	require.NoError(t, err)
	require.Equal(t, int64(2), secondUserID)

	require.Len(t, passwordHashes, 2)
	require.NotEqual(t, firstReq.Password, passwordHashes[0])
	require.NotEqual(t, secondReq.Password, passwordHashes[1])
	require.NotEqual(t, passwordHashes[0], passwordHashes[1])
	require.NoError(t, helpers.CompareHashAndPassword(passwordHashes[0], firstReq.Password))
	require.NoError(t, helpers.CompareHashAndPassword(passwordHashes[1], secondReq.Password))
}
