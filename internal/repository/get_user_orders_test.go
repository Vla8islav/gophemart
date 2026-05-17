package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Vla8islav/gophemart/internal/config"
	"github.com/Vla8islav/gophemart/internal/domain"
	"github.com/Vla8islav/gophemart/internal/helpers"
	"github.com/stretchr/testify/require"
)

func TestPostgresStorage_GetUserOrders(t *testing.T) {
	cfg := config.ReadFlagsServer(nil)
	storage := InitTestPostgresStorage(t, cfg)

	ctx := context.Background()

	userID, err := storage.CreateUser(ctx, domain.CreateUserParams{
		Login:        helpers.UniqueLogin("get-user-orders-user"),
		PasswordHash: "hash",
	})
	require.NoError(t, err)

	otherUserID, err := storage.CreateUser(ctx, domain.CreateUserParams{
		Login:        helpers.UniqueLogin("get-user-orders-other-user"),
		PasswordHash: "hash",
	})
	require.NoError(t, err)

	firstUploadedAt := time.Date(2020, 12, 10, 15, 12, 1, 0, time.UTC)
	secondUploadedAt := time.Date(2020, 12, 10, 15, 15, 45, 0, time.UTC)
	thirdUploadedAt := time.Date(2020, 12, 9, 16, 9, 53, 0, time.UTC)

	insertTestOrder(t, ctx, storage, userID, "12345678903", "PROCESSING", nil, firstUploadedAt)
	insertTestOrder(t, ctx, storage, userID, "9278923470", "PROCESSED", ptrInt64(500), secondUploadedAt)
	insertTestOrder(t, ctx, storage, userID, "346436439", "INVALID", nil, thirdUploadedAt)
	insertTestOrder(t, ctx, storage, otherUserID, "79927398713", "PROCESSED", ptrInt64(100), time.Date(2021, 1, 1, 12, 0, 0, 0, time.UTC))

	orders, err := storage.GetUserOrders(ctx, userID)
	require.NoError(t, err)

	require.Len(t, orders, 3)
	require.Equal(t, "9278923470", orders[0].Number)
	require.Equal(t, "PROCESSED", orders[0].Status)
	require.NotNil(t, orders[0].Accrual)
	require.Equal(t, int64(500), *orders[0].Accrual)
	require.True(t, secondUploadedAt.Equal(orders[0].UploadedAt))

	require.Equal(t, "12345678903", orders[1].Number)
	require.Equal(t, "PROCESSING", orders[1].Status)
	require.Nil(t, orders[1].Accrual)
	require.True(t, firstUploadedAt.Equal(orders[1].UploadedAt))

	require.Equal(t, "346436439", orders[2].Number)
	require.Equal(t, "INVALID", orders[2].Status)
	require.Nil(t, orders[2].Accrual)
	require.True(t, thirdUploadedAt.Equal(orders[2].UploadedAt))
}

func TestPostgresStorage_GetUserOrders_Empty(t *testing.T) {
	cfg := config.ReadFlagsServer(nil)
	storage := InitTestPostgresStorage(t, cfg)

	ctx := context.Background()

	userID, err := storage.CreateUser(ctx, domain.CreateUserParams{
		Login:        helpers.UniqueLogin("get-user-orders-empty-user"),
		PasswordHash: "hash",
	})
	require.NoError(t, err)

	orders, err := storage.GetUserOrders(ctx, userID)
	require.NoError(t, err)
	require.Empty(t, orders)
}

func insertTestOrder(
	t *testing.T,
	ctx context.Context,
	storage *PostgresStorage,
	userID int64,
	number string,
	status string,
	accrual *int64,
	uploadedAt time.Time,
) {
	t.Helper()

	_, err := storage.db.ExecContext(ctx, `
		INSERT INTO orders (number, user_id, status, accrual, uploaded_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $5)
	`, number, userID, status, accrual, uploadedAt)
	require.NoError(t, err)
}

func ptrInt64(v int64) *int64 {
	return &v
}
