package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Vla8islav/gophemart/internal/config"
	"github.com/Vla8islav/gophemart/internal/domain"
	"github.com/Vla8islav/gophemart/internal/helpers"
	"github.com/stretchr/testify/require"
)

func TestPostgresStorage_GetUserWithdrawals(t *testing.T) {
	cfg := config.ReadFlagsServer(nil)
	storage := InitTestPostgresStorage(t, cfg)

	ctx := context.Background()
	suffix := time.Now().UnixNano()

	userID, err := storage.CreateUser(ctx, domain.CreateUserParams{
		Login:        helpers.UniqueLogin(fmt.Sprintf("get-user-withdrawals-user-test-%d", suffix)),
		PasswordHash: "hashed-password",
	})
	require.NoError(t, err)
	require.Greater(t, userID, int64(0))

	otherUserID, err := storage.CreateUser(ctx, domain.CreateUserParams{
		Login:        helpers.UniqueLogin(fmt.Sprintf("get-user-withdrawals-other-user-test-%d", suffix)),
		PasswordHash: "hashed-password",
	})
	require.NoError(t, err)
	require.Greater(t, otherUserID, int64(0))

	oldOrderNumber := fmt.Sprintf("get-user-withdrawals-old-%d", suffix)
	newOrderNumber := fmt.Sprintf("get-user-withdrawals-new-%d", suffix)
	otherUserOrderNumber := fmt.Sprintf("get-user-withdrawals-other-%d", suffix)

	oldAmount := int64(10000)
	newAmount := int64(25000)
	otherUserAmount := int64(50000)

	_, err = storage.db.ExecContext(ctx,
		`INSERT INTO withdrawals (user_id, order_number, amount, processed_at)
		 VALUES
		     ($1, $2, $3, now() - interval '2 minutes'),
		     ($1, $4, $5, now() - interval '1 minute'),
		     ($6, $7, $8, now())`,
		userID, oldOrderNumber, oldAmount,
		newOrderNumber, newAmount,
		otherUserID, otherUserOrderNumber, otherUserAmount,
	)
	require.NoError(t, err)

	withdrawals, err := storage.GetUserWithdrawals(ctx, userID)
	require.NoError(t, err)
	require.Len(t, withdrawals, 2)

	require.Equal(t, userID, withdrawals[0].UserID)
	require.Equal(t, newOrderNumber, withdrawals[0].OrderNumber)
	require.Equal(t, newAmount, withdrawals[0].Amount)
	require.False(t, withdrawals[0].ProcessedAt.IsZero())

	require.Equal(t, userID, withdrawals[1].UserID)
	require.Equal(t, oldOrderNumber, withdrawals[1].OrderNumber)
	require.Equal(t, oldAmount, withdrawals[1].Amount)
	require.False(t, withdrawals[1].ProcessedAt.IsZero())

	require.True(t, withdrawals[0].ProcessedAt.After(withdrawals[1].ProcessedAt))

	for _, withdrawal := range withdrawals {
		require.NotEqual(t, otherUserID, withdrawal.UserID)
		require.NotEqual(t, otherUserOrderNumber, withdrawal.OrderNumber)
	}
}

func TestPostgresStorage_GetUserWithdrawals_Empty(t *testing.T) {
	cfg := config.ReadFlagsServer(nil)
	storage := InitTestPostgresStorage(t, cfg)

	ctx := context.Background()
	suffix := time.Now().UnixNano()

	userID, err := storage.CreateUser(ctx, domain.CreateUserParams{
		Login:        helpers.UniqueLogin(fmt.Sprintf("get-user-withdrawals-empty-user-test-%d", suffix)),
		PasswordHash: "hashed-password",
	})
	require.NoError(t, err)
	require.Greater(t, userID, int64(0))

	otherUserID, err := storage.CreateUser(ctx, domain.CreateUserParams{
		Login:        helpers.UniqueLogin(fmt.Sprintf("get-user-withdrawals-empty-other-user-test-%d", suffix)),
		PasswordHash: "hashed-password",
	})
	require.NoError(t, err)
	require.Greater(t, otherUserID, int64(0))

	otherUserOrderNumber := fmt.Sprintf("get-user-withdrawals-empty-other-%d", suffix)

	_, err = storage.db.ExecContext(ctx,
		`INSERT INTO withdrawals (user_id, order_number, amount, processed_at)
		 VALUES ($1, $2, $3, now())`,
		otherUserID, otherUserOrderNumber, int64(50000),
	)
	require.NoError(t, err)

	withdrawals, err := storage.GetUserWithdrawals(ctx, userID)
	require.NoError(t, err)
	require.Empty(t, withdrawals)
}
