package repository

import (
	"context"
	"testing"

	"github.com/Vla8islav/gophemart/internal/config"
	"github.com/Vla8islav/gophemart/internal/domain"
	"github.com/Vla8islav/gophemart/internal/helpers"
	"github.com/stretchr/testify/require"
)

func TestPostgresStorage_GetActiveOrders(t *testing.T) {
	cfg := config.ReadFlagsServer(nil)
	storage := InitTestPostgresStorage(t, cfg)

	ctx := context.Background()

	userID, err := storage.CreateUser(ctx, domain.CreateUserParams{
		Login:        helpers.UniqueLogin("get-active-orders-user-test"),
		PasswordHash: "hashed-password",
	})
	require.NoError(t, err)
	require.Greater(t, userID, int64(0))

	newOrderNumber := helpers.UniqueLogin("get-active-orders-new")
	processingOrderNumber := helpers.UniqueLogin("get-active-orders-processing")
	processedOrderNumber := helpers.UniqueLogin("get-active-orders-processed")
	invalidOrderNumber := helpers.UniqueLogin("get-active-orders-invalid")

	processingAccrual := int64(25000)
	processedAccrual := int64(50000)

	_, err = storage.db.ExecContext(ctx,
		`INSERT INTO orders (number, user_id, status, accrual, uploaded_at)
		 VALUES
		     ($1, $2, 'NEW', NULL, now() - interval '4 minutes'),
		     ($3, $2, 'PROCESSING', $4, now() - interval '3 minutes'),
		     ($5, $2, 'PROCESSED', $6, now() - interval '2 minutes'),
		     ($7, $2, 'INVALID', NULL, now() - interval '1 minute')`,
		newOrderNumber, userID,
		processingOrderNumber, processingAccrual,
		processedOrderNumber, processedAccrual,
		invalidOrderNumber,
	)
	require.NoError(t, err)

	orders, err := storage.GetActiveOrders(ctx)
	require.NoError(t, err)
	require.Len(t, orders, 2)

	require.Equal(t, newOrderNumber, orders[0].Number)
	require.Equal(t, "NEW", orders[0].Status)
	require.Nil(t, orders[0].Accrual)
	require.False(t, orders[0].UploadedAt.IsZero())

	require.Equal(t, processingOrderNumber, orders[1].Number)
	require.Equal(t, "PROCESSING", orders[1].Status)
	require.NotNil(t, orders[1].Accrual)
	require.Equal(t, processingAccrual, *orders[1].Accrual)
	require.False(t, orders[1].UploadedAt.IsZero())
}

func TestPostgresStorage_GetActiveOrders_Empty(t *testing.T) {
	cfg := config.ReadFlagsServer(nil)
	storage := InitTestPostgresStorage(t, cfg)

	ctx := context.Background()

	userID, err := storage.CreateUser(ctx, domain.CreateUserParams{
		Login:        helpers.UniqueLogin("get-active-orders-empty-user-test"),
		PasswordHash: "hashed-password",
	})
	require.NoError(t, err)
	require.Greater(t, userID, int64(0))

	processedOrderNumber := helpers.UniqueLogin("get-active-orders-empty-processed")
	invalidOrderNumber := helpers.UniqueLogin("get-active-orders-empty-invalid")

	processedAccrual := int64(50000)

	_, err = storage.db.ExecContext(ctx,
		`INSERT INTO orders (number, user_id, status, accrual)
		 VALUES
		     ($1, $2, 'PROCESSED', $3),
		     ($4, $2, 'INVALID', NULL)`,
		processedOrderNumber, userID, processedAccrual,
		invalidOrderNumber,
	)
	require.NoError(t, err)

	orders, err := storage.GetActiveOrders(ctx)
	require.NoError(t, err)
	require.Empty(t, orders)
}
