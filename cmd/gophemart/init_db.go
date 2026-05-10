package main

import (
	"context"
	"fmt"

	"github.com/Vla8islav/gophemart/internal/config"
	"github.com/Vla8islav/gophemart/internal/domain"
	"github.com/Vla8islav/gophemart/internal/repository"
	"go.uber.org/zap"
)

func initDB(ctx context.Context, currentConfig *config.OptionsServer, logger *zap.Logger) (domain.GophemartRepository, error) {
	var db domain.GophemartRepository
	var err error

	// Case 1
	if currentConfig.DatabaseDSN.BeenSet {
		db, err = repository.NewPostgresStorage(currentConfig, currentConfig.MigrationsFolder.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize metrics repository: %w", err)
		}
		return db, nil
	}

	return nil, fmt.Errorf("something strange happened: " +
		"restore and connection string parameters are incorrect")
}
