package repository

import (
	"context"
	"testing"

	"github.com/Vla8islav/gophemart/internal/config"
	"github.com/Vla8islav/gophemart/internal/domain"
	"github.com/Vla8islav/gophemart/internal/helpers"
	"github.com/stretchr/testify/require"
)

func TestPostgresStorage_GetUserBalance(t *testing.T) {
	cfg := config.ReadFlagsServer(nil)
	storage := InitTestPostgresStorage(t, cfg)

	ctx := context.Background()

	userID, err := storage.CreateUser(ctx, domain.CreateUserParams{
		Login:        helpers.UniqueLogin("balance-user-test"),
		PasswordHash: "hashed-password",
	})
	require.NoError(t, err)
	require.Greater(t, userID, int64(0))

	_, err = storage.db.ExecContext(ctx,
		`INSERT INTO orders (number, user_id, status, accrual)
		 VALUES
		     ($1, $2, 'PROCESSED', $3),
		     ($4, $2, 'PROCESSED', $5),
		     ($6, $2, 'NEW', NULL),
		     ($7, $2, 'INVALID', NULL)`,
		helpers.UniqueLogin("balance-order-processed-1"), userID, int64(212452),
		helpers.UniqueLogin("balance-order-processed-2"), int64(1000),
		helpers.UniqueLogin("balance-order-new"),
		helpers.UniqueLogin("balance-order-invalid"),
	)
	require.NoError(t, err)

	_, err = storage.db.ExecContext(ctx,
		`INSERT INTO withdrawals (user_id, order_number, amount)
		 VALUES
		     ($1, $2, $3),
		     ($1, $4, $5)`,
		userID,
		helpers.UniqueLogin("balance-withdrawal-1"), int64(1000),
		helpers.UniqueLogin("balance-withdrawal-2"), int64(500),
	)
	require.NoError(t, err)

	balance, err := storage.GetUserBalance(ctx, userID)
	require.NoError(t, err)
	require.NotNil(t, balance)

	require.Equal(t, int64(211952), balance.Current)
	require.Equal(t, int64(1500), balance.Withdrawn)
}

func TestPostgresStorage_GetUserBalanceEmptyBalance(t *testing.T) {
	cfg := config.ReadFlagsServer(nil)
	storage := InitTestPostgresStorage(t, cfg)

	ctx := context.Background()

	userID, err := storage.CreateUser(ctx, domain.CreateUserParams{
		Login:        helpers.UniqueLogin("empty-balance-user-test"),
		PasswordHash: "hashed-password",
	})
	require.NoError(t, err)
	require.Greater(t, userID, int64(0))

	balance, err := storage.GetUserBalance(ctx, userID)
	require.NoError(t, err)
	require.NotNil(t, balance)

	require.Equal(t, int64(0), balance.Current)
	require.Equal(t, int64(0), balance.Withdrawn)
}
