package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/Vla8islav/gophemart/internal/config"
	"github.com/Vla8islav/gophemart/internal/domain"
	"github.com/Vla8islav/gophemart/internal/helpers"
	"github.com/stretchr/testify/require"
)

func TestPostgresStorage_CreateOrder(t *testing.T) {
	cfg := config.ReadFlagsServer(nil)
	storage := InitTestPostgresStorage(t, cfg)

	ctx := context.Background()

	userID, err := storage.CreateUser(ctx, domain.CreateUserParams{
		Login:        helpers.UniqueLogin("create-order-user"),
		PasswordHash: "hashed-password",
	})
	require.NoError(t, err)
	require.Greater(t, userID, int64(0))

	orderNumber := helpers.UniqueLogin("create-order-number")

	err = storage.CreateOrder(ctx, userID, orderNumber)
	require.NoError(t, err)

	var storedUserID int64
	var storedStatus string

	err = storage.db.QueryRowContext(ctx,
		`SELECT user_id, status
		 FROM orders
		 WHERE number = $1`,
		orderNumber,
	).Scan(&storedUserID, &storedStatus)
	require.NoError(t, err)

	require.Equal(t, userID, storedUserID)
	require.Equal(t, "NEW", storedStatus)
}

func TestPostgresStorage_CreateOrder_AlreadyUploadedBySameUser(t *testing.T) {
	cfg := config.ReadFlagsServer(nil)
	storage := InitTestPostgresStorage(t, cfg)

	ctx := context.Background()

	userID, err := storage.CreateUser(ctx, domain.CreateUserParams{
		Login:        helpers.UniqueLogin("create-order-same-user"),
		PasswordHash: "hashed-password",
	})
	require.NoError(t, err)
	require.Greater(t, userID, int64(0))

	orderNumber := helpers.UniqueLogin("same-user-order-number")

	err = storage.CreateOrder(ctx, userID, orderNumber)
	require.NoError(t, err)

	err = storage.CreateOrder(ctx, userID, orderNumber)
	require.Error(t, err)
	require.True(t, errors.Is(err, domain.ErrOrderAlreadyUploadedByUser))
}

func TestPostgresStorage_CreateOrder_AlreadyUploadedByAnotherUser(t *testing.T) {
	cfg := config.ReadFlagsServer(nil)
	storage := InitTestPostgresStorage(t, cfg)

	ctx := context.Background()

	firstUserID, err := storage.CreateUser(ctx, domain.CreateUserParams{
		Login:        helpers.UniqueLogin("create-order-first-user"),
		PasswordHash: "hashed-password",
	})
	require.NoError(t, err)
	require.Greater(t, firstUserID, int64(0))

	secondUserID, err := storage.CreateUser(ctx, domain.CreateUserParams{
		Login:        helpers.UniqueLogin("create-order-second-user"),
		PasswordHash: "hashed-password",
	})
	require.NoError(t, err)
	require.Greater(t, secondUserID, int64(0))
	require.NotEqual(t, firstUserID, secondUserID)

	orderNumber := helpers.UniqueLogin("another-user-order-number")

	err = storage.CreateOrder(ctx, firstUserID, orderNumber)
	require.NoError(t, err)

	err = storage.CreateOrder(ctx, secondUserID, orderNumber)
	require.Error(t, err)
	require.True(t, errors.Is(err, domain.ErrOrderAlreadyUploadedByAnotherUser))
}
