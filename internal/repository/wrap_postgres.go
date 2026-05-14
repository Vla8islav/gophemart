package repository

import (
	"fmt"

	"github.com/Vla8islav/gophemart/internal/config"
	"github.com/Vla8islav/gophemart/internal/domain"
)

func WrapPostgres(currentConfig *config.OptionsServer) (domain.GophermartRepository, error) {
	var db domain.GophermartRepository
	var err error

	// Case 1
	if currentConfig.DatabaseURI.BeenSet {
		db, err = NewPostgresStorage(currentConfig, currentConfig.MigrationsFolder.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize metrics repository: %w", err)
		}
		return db, nil
	}

	return nil, fmt.Errorf("something strange happened: " +
		"restore and connection string parameters are incorrect")
}
