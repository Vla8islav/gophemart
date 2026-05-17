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

func TestPostgresStorage_GetActiveOrders(t *testing.T) {
	cfg := config.ReadFlagsServer(nil)
	storage := InitTestPostgresStorage(t, cfg)

	ctx := context.Background()
	suffix := time.Now().UnixNano()

	userID, err := storage.CreateUser(ctx, domain.CreateUserParams{
		Login:        helpers.UniqueLogin(fmt.Sprintf("get-active-orders-user-test-%d", suffix)),
		PasswordHash: "hashed-password",
	})
	require.NoError(t, err)
	require.Greater(t, userID, int64(0))

	newOrderNumber := fmt.Sprintf("get-active-orders-new-%d", suffix)
	processingOrderNumber := fmt.Sprintf("get-active-orders-processing-%d", suffix)
	processedOrderNumber := fmt.Sprintf("get-active-orders-processed-%d", suffix)
	invalidOrderNumber := fmt.Sprintf("get-active-orders-invalid-%d", suffix)

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

	ordersByNumber := make(map[string]domain.UserOrder, len(orders))
	for _, order := range orders {
		ordersByNumber[order.Number] = order
	}

	newOrder, ok := ordersByNumber[newOrderNumber]
	require.True(t, ok)
	require.Equal(t, "NEW", newOrder.Status)
	require.Nil(t, newOrder.Accrual)
	require.False(t, newOrder.UploadedAt.IsZero())

	processingOrder, ok := ordersByNumber[processingOrderNumber]
	require.True(t, ok)
	require.Equal(t, "PROCESSING", processingOrder.Status)
	require.NotNil(t, processingOrder.Accrual)
	require.Equal(t, processingAccrual, *processingOrder.Accrual)
	require.False(t, processingOrder.UploadedAt.IsZero())

	_, ok = ordersByNumber[processedOrderNumber]
	require.False(t, ok)

	_, ok = ordersByNumber[invalidOrderNumber]
	require.False(t, ok)
}

func TestPostgresStorage_GetActiveOrders_Empty(t *testing.T) {
	cfg := config.ReadFlagsServer(nil)
	storage := InitTestPostgresStorage(t, cfg)

	ctx := context.Background()
	suffix := time.Now().UnixNano()

	userID, err := storage.CreateUser(ctx, domain.CreateUserParams{
		Login:        helpers.UniqueLogin(fmt.Sprintf("get-active-orders-empty-user-test-%d", suffix)),
		PasswordHash: "hashed-password",
	})
	require.NoError(t, err)
	require.Greater(t, userID, int64(0))

	processedOrderNumber := fmt.Sprintf("get-active-orders-empty-processed-%d", suffix)
	invalidOrderNumber := fmt.Sprintf("get-active-orders-empty-invalid-%d", suffix)

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

	ordersByNumber := make(map[string]domain.UserOrder, len(orders))
	for _, order := range orders {
		ordersByNumber[order.Number] = order
	}

	_, ok := ordersByNumber[processedOrderNumber]
	require.False(t, ok)

	_, ok = ordersByNumber[invalidOrderNumber]
	require.False(t, ok)
}
