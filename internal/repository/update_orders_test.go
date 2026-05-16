package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Vla8islav/gophemart/internal/config"
	"github.com/Vla8islav/gophemart/internal/domain"
	"github.com/Vla8islav/gophemart/internal/helpers"
	"github.com/stretchr/testify/require"
)

func TestPostgresStorage_UpdateOrders(t *testing.T) {
	cfg := config.ReadFlagsServer(nil)
	storage := InitTestPostgresStorage(t, cfg)

	ctx := context.Background()

	userID, err := storage.CreateUser(ctx, domain.CreateUserParams{
		Login:        helpers.UniqueLogin("update-orders-user-test"),
		PasswordHash: "hashed-password",
	})
	require.NoError(t, err)
	require.Greater(t, userID, int64(0))

	processedOrderNumber := helpers.UniqueLogin("update-order-processed")
	invalidOrderNumber := helpers.UniqueLogin("update-order-invalid")
	untouchedOrderNumber := helpers.UniqueLogin("update-order-untouched")

	_, err = storage.db.ExecContext(ctx,
		`INSERT INTO orders (number, user_id, status, accrual)
		 VALUES
		     ($1, $2, 'NEW', NULL),
		     ($3, $2, 'PROCESSING', NULL),
		     ($4, $2, 'NEW', NULL)`,
		processedOrderNumber, userID,
		invalidOrderNumber,
		untouchedOrderNumber,
	)
	require.NoError(t, err)

	accrual := int64(50000)
	updates := []domain.UpdateOrderParams{
		{
			Number:  processedOrderNumber,
			Status:  "PROCESSED",
			Accrual: &accrual,
		},
		{
			Number:  invalidOrderNumber,
			Status:  "INVALID",
			Accrual: nil,
		},
	}

	err = storage.UpdateOrders(ctx, updates)
	require.NoError(t, err)

	var processedStatus string
	var processedAccrual sql.NullInt64
	err = storage.db.QueryRowContext(ctx,
		`SELECT status, accrual FROM orders WHERE number = $1`,
		processedOrderNumber,
	).Scan(&processedStatus, &processedAccrual)
	require.NoError(t, err)
	require.Equal(t, "PROCESSED", processedStatus)
	require.True(t, processedAccrual.Valid)
	require.Equal(t, accrual, processedAccrual.Int64)

	var invalidStatus string
	var invalidAccrual sql.NullInt64
	err = storage.db.QueryRowContext(ctx,
		`SELECT status, accrual FROM orders WHERE number = $1`,
		invalidOrderNumber,
	).Scan(&invalidStatus, &invalidAccrual)
	require.NoError(t, err)
	require.Equal(t, "INVALID", invalidStatus)
	require.False(t, invalidAccrual.Valid)

	var untouchedStatus string
	var untouchedAccrual sql.NullInt64
	err = storage.db.QueryRowContext(ctx,
		`SELECT status, accrual FROM orders WHERE number = $1`,
		untouchedOrderNumber,
	).Scan(&untouchedStatus, &untouchedAccrual)
	require.NoError(t, err)
	require.Equal(t, "NEW", untouchedStatus)
	require.False(t, untouchedAccrual.Valid)
}

func TestPostgresStorage_UpdateOrdersEmptyUpdates(t *testing.T) {
	cfg := config.ReadFlagsServer(nil)
	storage := InitTestPostgresStorage(t, cfg)

	ctx := context.Background()

	err := storage.UpdateOrders(ctx, nil)
	require.NoError(t, err)

	err = storage.UpdateOrders(ctx, []domain.UpdateOrderParams{})
	require.NoError(t, err)
}
