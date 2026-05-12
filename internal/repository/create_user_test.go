package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Vla8islav/gophemart/internal/config"
	"github.com/Vla8islav/gophemart/internal/models"
	"github.com/stretchr/testify/require"
)

func uniqueLogin(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

func TestPostgresStorage_CreateUser(t *testing.T) {
	cfg := config.ReadFlagsServer(nil)
	storage := InitTestPostgresStorage(t, cfg)

	ctx := context.Background()

	user := models.User{
		Login:    uniqueLogin("create-user-test"),
		Password: "hashed-password",
	}

	userID, err := storage.CreateUser(ctx, user)
	require.NoError(t, err)
	require.Greater(t, userID, int64(0))

	var login string
	var passwordHash string

	err = storage.db.QueryRowContext(ctx,
		`SELECT login, password_hash
		 FROM users
		 WHERE id = $1`,
		userID,
	).Scan(&login, &passwordHash)
	require.NoError(t, err)

	require.Equal(t, user.Login, login)
	require.Equal(t, user.Password, passwordHash)
}

func TestPostgresStorage_CreateUser_DuplicateLogin(t *testing.T) {
	cfg := config.ReadFlagsServer(nil)
	storage := InitTestPostgresStorage(t, cfg)

	ctx := context.Background()

	user := models.User{
		Login:    uniqueLogin("duplicate-user-test"),
		Password: "hashed-password",
	}

	userID, err := storage.CreateUser(ctx, user)
	require.NoError(t, err)
	require.Greater(t, userID, int64(0))

	duplicateUserID, err := storage.CreateUser(ctx, user)
	require.Error(t, err)
	require.Zero(t, duplicateUserID)
}

func TestPostgresStorage_CreateUser_AllowsSamePasswordHash(t *testing.T) {
	cfg := config.ReadFlagsServer(nil)
	storage := InitTestPostgresStorage(t, cfg)

	ctx := context.Background()

	firstUser := models.User{
		Login:    uniqueLogin("same-password-user-1"),
		Password: "same-hashed-password",
	}
	secondUser := models.User{
		Login:    uniqueLogin("same-password-user-2"),
		Password: "same-hashed-password",
	}

	firstUserID, err := storage.CreateUser(ctx, firstUser)
	require.NoError(t, err)
	require.Greater(t, firstUserID, int64(0))

	secondUserID, err := storage.CreateUser(ctx, secondUser)
	require.NoError(t, err)
	require.Greater(t, secondUserID, int64(0))
	require.NotEqual(t, firstUserID, secondUserID)
}
