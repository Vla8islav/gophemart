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

func TestPostgresStorage_WithdrawFromUserBalance(t *testing.T) {
	cfg := config.ReadFlagsServer(nil)
	storage := InitTestPostgresStorage(t, cfg)

	ctx := context.Background()

	userID, err := storage.CreateUser(ctx, domain.CreateUserParams{
		Login:        helpers.UniqueLogin("withdraw-user-test"),
		PasswordHash: "hashed-password",
	})
	require.NoError(t, err)
	require.Greater(t, userID, int64(0))

	orderNumber := helpers.UniqueLogin("withdraw-source-order")
	_, err = storage.db.ExecContext(ctx,
		`INSERT INTO orders (number, user_id, status, accrual)
		 VALUES ($1, $2, 'PROCESSED', $3)`,
		orderNumber, userID, int64(50000),
	)
	require.NoError(t, err)

	withdrawOrderNumber := helpers.UniqueLogin("withdraw-payment-order")
	request := domain.UserBalanceWithdraw{
		Order: withdrawOrderNumber,
		Sum:   30000,
	}

	err = storage.WithdrawFromUserBalance(ctx, userID, request)
	require.NoError(t, err)

	var amount int64
	err = storage.db.QueryRowContext(ctx,
		`SELECT amount FROM withdrawals WHERE user_id = $1 AND order_number = $2`,
		userID, withdrawOrderNumber,
	).Scan(&amount)
	require.NoError(t, err)
	require.Equal(t, request.Sum, amount)
}

func TestPostgresStorage_WithdrawFromUserBalanceNotEnoughMoney(t *testing.T) {
	cfg := config.ReadFlagsServer(nil)
	storage := InitTestPostgresStorage(t, cfg)

	ctx := context.Background()

	userID, err := storage.CreateUser(ctx, domain.CreateUserParams{
		Login:        helpers.UniqueLogin("withdraw-not-enough-money-user-test"),
		PasswordHash: "hashed-password",
	})
	require.NoError(t, err)
	require.Greater(t, userID, int64(0))

	orderNumber := helpers.UniqueLogin("withdraw-not-enough-money-source-order")
	_, err = storage.db.ExecContext(ctx,
		`INSERT INTO orders (number, user_id, status, accrual)
		 VALUES ($1, $2, 'PROCESSED', $3)`,
		orderNumber, userID, int64(10000),
	)
	require.NoError(t, err)

	withdrawOrderNumber := helpers.UniqueLogin("withdraw-not-enough-money-payment-order")
	request := domain.UserBalanceWithdraw{
		Order: withdrawOrderNumber,
		Sum:   30000,
	}

	err = storage.WithdrawFromUserBalance(ctx, userID, request)
	require.Error(t, err)
	require.True(t, errors.Is(err, domain.ErrNotEnoughMoney))

	var withdrawalsCount int
	err = storage.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM withdrawals WHERE user_id = $1 AND order_number = $2`,
		userID, withdrawOrderNumber,
	).Scan(&withdrawalsCount)
	require.NoError(t, err)
	require.Zero(t, withdrawalsCount)
}

func TestPostgresStorage_WithdrawFromUserBalanceAccountsForPreviousWithdrawals(t *testing.T) {
	cfg := config.ReadFlagsServer(nil)
	storage := InitTestPostgresStorage(t, cfg)

	ctx := context.Background()

	userID, err := storage.CreateUser(ctx, domain.CreateUserParams{
		Login:        helpers.UniqueLogin("withdraw-previous-withdrawals-user-test"),
		PasswordHash: "hashed-password",
	})
	require.NoError(t, err)
	require.Greater(t, userID, int64(0))

	orderNumber := helpers.UniqueLogin("withdraw-previous-withdrawals-source-order")
	_, err = storage.db.ExecContext(ctx,
		`INSERT INTO orders (number, user_id, status, accrual)
		 VALUES ($1, $2, 'PROCESSED', $3)`,
		orderNumber, userID, int64(50000),
	)
	require.NoError(t, err)

	firstWithdrawal := domain.UserBalanceWithdraw{
		Order: helpers.UniqueLogin("withdraw-previous-withdrawals-first"),
		Sum:   40000,
	}
	err = storage.WithdrawFromUserBalance(ctx, userID, firstWithdrawal)
	require.NoError(t, err)

	secondWithdrawal := domain.UserBalanceWithdraw{
		Order: helpers.UniqueLogin("withdraw-previous-withdrawals-second"),
		Sum:   20000,
	}
	err = storage.WithdrawFromUserBalance(ctx, userID, secondWithdrawal)
	require.Error(t, err)
	require.True(t, errors.Is(err, domain.ErrNotEnoughMoney))

	var withdrawalsCount int
	err = storage.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM withdrawals WHERE user_id = $1`,
		userID,
	).Scan(&withdrawalsCount)
	require.NoError(t, err)
	require.Equal(t, 1, withdrawalsCount)
}
