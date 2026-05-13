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

type e2eServerManager struct {
	mu sync.Mutex

	cfg  *config.OptionsServer
	stop func()
	err  error

	refCount     int
	destroyTimer *time.Timer
}

var e2eServer e2eServerManager

func initE2ETestServer(t *testing.T) *config.OptionsServer {
	t.Helper()

	cfg, release, err := e2eServer.acquire(t)
	require.NoError(t, err)

	t.Cleanup(release)

	return cfg
}

func (m *e2eServerManager) acquire(t *testing.T) (*config.OptionsServer, func(), error) {
	t.Helper()

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.destroyTimer != nil {
		m.destroyTimer.Stop()
		m.destroyTimer = nil
	}

	if m.cfg == nil {
		cfg, stop, err := startE2ETestServer(t)
		if err != nil {
			return nil, nil, err
		}

		m.cfg = cfg
		m.stop = stop
		m.err = nil
	}

	m.refCount++

	released := false
	release := func() {
		m.mu.Lock()
		defer m.mu.Unlock()

		if released {
			return
		}
		released = true

		m.refCount--
		if m.refCount != 0 {
			return
		}

		m.destroyTimer = time.AfterFunc(time.Second, func() {
			m.mu.Lock()
			defer m.mu.Unlock()

			if m.refCount != 0 {
				return
			}

			if m.stop != nil {
				m.stop()
			}

			m.cfg = nil
			m.stop = nil
			m.err = nil
			m.destroyTimer = nil
		})
	}

	return m.cfg, release, m.err
}

func startE2ETestServer(t *testing.T) (*config.OptionsServer, func(), error) {
	t.Helper()

	lg, err := zap.NewProduction()
	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	cfg := config.ReadFlagsServer(nil)

	wrappedDB := repository.InitTestPostgresStorage(t, cfg)

	go func() {
		err := run.Run(ctx, wrappedDB, cfg, lg)
		if err != nil && ctx.Err() == nil {
			lg.Error("e2e server stopped", zap.Error(err))
		}
	}()

	waitForServer(t, cfg)

	stop := func() {
		cancel()
	}

	return cfg, stop, nil
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
