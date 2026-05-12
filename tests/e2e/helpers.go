package e2e

import (
	"context"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/Vla8islav/gophemart/internal/config"
	"github.com/Vla8islav/gophemart/internal/repository"
	"github.com/Vla8islav/gophemart/internal/run"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func initE2ETestServer(t *testing.T) *config.OptionsServer {
	t.Helper()

	lg, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to initialize lg: %v", err)
	}
	t.Cleanup(func() {
		_ = lg.Sync()
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	cfg := config.ReadFlagsServer(nil)

	wrappedDB := repository.InitTestPostgresStorage(t, cfg)

	go func() {
		err := run.Run(ctx, wrappedDB, cfg, lg)
		require.NoError(t, err)
	}()
	waitForServer(t, cfg)

	return cfg
}

func waitForServer(t *testing.T, cfg *config.OptionsServer) {
	pingURL := "http://" + cfg.ServerAddress.Value + "/api/ping"

	t.Helper()

	client := http.Client{
		Timeout: 200 * time.Millisecond,
	}

	deadline := time.Now().Add(5 * time.Second)
	var lastErr error

	for time.Now().Before(deadline) {
		resp, err := client.Get(pingURL)
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return
			}
		}

		lastErr = err
		time.Sleep(50 * time.Millisecond)
	}

	require.NoError(t, lastErr, "сервер не стал готов за необходимое время")
}
