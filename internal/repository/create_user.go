package repository

import (
	"context"
	"fmt"

	"github.com/Vla8islav/gophemart/internal/helpers"
	"github.com/Vla8islav/gophemart/internal/models"
)

func (s *PostgresStorage) CreateUser(ctx context.Context, user models.User) (int64, error) {
	var userID int64
	hash, err := helpers.HashPassword(user.Password)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate the hash for the new user %s: %w", user.Login, err)
	}

	err = s.withRetry(ctx, func() error {
		err := s.db.QueryRowContext(ctx,
			`INSERT INTO users (login, password_hash)
			 		VALUES ($1, $2)
			 		RETURNING id`,
			user.Login,
			hash,
		).Scan(&userID)
		if err != nil {
			return fmt.Errorf("failed to create a new user %s: %w", user.Login, err)
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return userID, nil
}
