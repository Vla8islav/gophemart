package e2e

import (
	"context"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/Vla8islav/gophemart/internal/config"
	"github.com/Vla8islav/gophemart/internal/repository"
	"github.com/Vla8islav/gophemart/internal/run"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var (
	e2eServerOnce sync.Once
	e2eServerCfg  *config.OptionsServer
	e2eServerErr  error
)

func initE2ETestServer(t *testing.T) *config.OptionsServer {
	t.Helper()

	e2eServerOnce.Do(func() {
		e2eServerCfg, e2eServerErr = startE2ETestServer(t)
	})

	require.NoError(t, e2eServerErr)

	return e2eServerCfg
}

func startE2ETestServer(t *testing.T) (*config.OptionsServer, error) {
	t.Helper()

	lg, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	cfg := config.ReadFlagsServer(nil)

	wrappedDB := repository.InitTestPostgresStorage(t, cfg)

	go func() {
		err := run.Run(ctx, wrappedDB, cfg, lg)
		if err != nil {
			lg.Error("e2e server stopped", zap.Error(err))
		}
	}()

	waitForServer(t, cfg)

	return cfg, nil
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
