package repository

import (
	"context"
	"testing"

	"github.com/Vla8islav/gophemart/internal/config"
	"github.com/Vla8islav/gophemart/internal/helpers"
	"github.com/Vla8islav/gophemart/internal/models"
	"github.com/stretchr/testify/require"
)

func TestPostgresStorage_CreateUser(t *testing.T) {
	cfg := config.ReadFlagsServer(nil)
	storage := InitTestPostgresStorage(t, cfg)

	ctx := context.Background()

	user := models.User{
		Login:    helpers.UniqueLogin("create-user-test"),
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
	require.NotEqual(t, user.Password, passwordHash)
	require.NoError(t, helpers.CompareHashAndPassword(passwordHash, user.Password))
}

func TestPostgresStorage_CreateUser_DuplicateLogin(t *testing.T) {
	cfg := config.ReadFlagsServer(nil)
	storage := InitTestPostgresStorage(t, cfg)

	ctx := context.Background()

	user := models.User{
		Login:    helpers.UniqueLogin("duplicate-user-test"),
		Password: "hashed-password",
	}

	userID, err := storage.CreateUser(ctx, user)
	require.NoError(t, err)
	require.Greater(t, userID, int64(0))

	duplicateUserID, err := storage.CreateUser(ctx, user)
	require.Error(t, err)
	require.Zero(t, duplicateUserID)
}

func TestPostgresStorage_CreateUser_HashesSamePasswordDifferently(t *testing.T) {
	cfg := config.ReadFlagsServer(nil)
	storage := InitTestPostgresStorage(t, cfg)

	ctx := context.Background()

	firstUser := models.User{
		Login:    helpers.UniqueLogin("same-password-user-1"),
		Password: "same-hashed-password",
	}
	secondUser := models.User{
		Login:    helpers.UniqueLogin("same-password-user-2"),
		Password: "same-hashed-password",
	}

	firstUserID, err := storage.CreateUser(ctx, firstUser)
	require.NoError(t, err)
	require.Greater(t, firstUserID, int64(0))

	secondUserID, err := storage.CreateUser(ctx, secondUser)
	require.NoError(t, err)
	require.Greater(t, secondUserID, int64(0))
	require.NotEqual(t, firstUserID, secondUserID)

	var firstPasswordHash string
	var secondPasswordHash string

	err = storage.db.QueryRowContext(ctx,
		`SELECT password_hash
		 FROM users
		 WHERE id = $1`,
		firstUserID,
	).Scan(&firstPasswordHash)
	require.NoError(t, err)

	err = storage.db.QueryRowContext(ctx,
		`SELECT password_hash
		 FROM users
		 WHERE id = $1`,
		secondUserID,
	).Scan(&secondPasswordHash)
	require.NoError(t, err)

	require.NotEqual(t, firstUser.Password, firstPasswordHash)
	require.NotEqual(t, secondUser.Password, secondPasswordHash)
	require.NotEqual(t, firstPasswordHash, secondPasswordHash)
	require.NoError(t, helpers.CompareHashAndPassword(firstPasswordHash, firstUser.Password))
	require.NoError(t, helpers.CompareHashAndPassword(secondPasswordHash, secondUser.Password))
}
