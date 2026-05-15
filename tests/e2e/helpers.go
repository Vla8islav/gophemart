package e2e

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"
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

const e2eAccrualDatabaseURI = "postgres://default_user:default_password@localhost:5432/accrual?sslmode=disable"

var e2eAccrualOrderNumbers = []string{
	"9278923470",
	"12345678903",
	"346436439",
}

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

	accrualStop, accrualAddress, err := startE2EAccrualServer(t)
	if err != nil {
		cancel()
		return nil, nil, err
	}
	cfg.AccrualAddress.Value = accrualAddress

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
		accrualStop()
	}

	return cfg, stop, nil
}

func startE2EAccrualServer(t *testing.T) (func(), string, error) {
	t.Helper()

	address, err := freeLocalAddress()
	if err != nil {
		return nil, "", err
	}

	binaryPath := filepath.Join("..", "..", "cmd", "accrual", fmt.Sprintf("accrual_%s_%s", runtime.GOOS, runtime.GOARCH))
	if runtime.GOOS == "windows" {
		binaryPath += ".exe"
	}

	cmd := exec.Command(binaryPath,
		"-a", address,
		"-d", e2eAccrualDatabaseURI,
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, "", err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, "", err
	}

	if err := cmd.Start(); err != nil {
		return nil, "", err
	}

	logProcessOutput(t, "accrual stdout", stdout)
	logProcessOutput(t, "accrual stderr", stderr)

	stop := func() {
		if cmd.Process == nil {
			return
		}

		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	}

	waitForAccrualServer(t, address)
	seedAccrualServer(t, address)

	return stop, address, nil
}

func logProcessOutput(t *testing.T, name string, reader io.Reader) {
	t.Helper()

	go func() {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			t.Logf("%s: %s", name, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			t.Logf("%s read error: %v", name, err)
		}
	}()
}

func seedAccrualServer(t *testing.T, address string) {
	t.Helper()

	client := http.Client{Timeout: 2 * time.Second}
	baseURL := "http://" + address

	postJSON(t, client, baseURL+"/api/goods", map[string]any{
		"match":       "Bork",
		"reward":      10,
		"reward_type": "%",
	})

	for _, orderNumber := range e2eAccrualOrderNumbers {
		postJSON(t, client, baseURL+"/api/orders", map[string]any{
			"order": orderNumber,
			"goods": []map[string]any{
				{
					"description": "Чайник Bork",
					"price":       7000,
				},
			},
		})
	}
}

func postJSON(t *testing.T, client http.Client, url string, payload any) {
	t.Helper()

	body, err := json.Marshal(payload)
	require.NoError(t, err)

	resp, err := client.Post(url, "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Truef(t,
		resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices,
		"unexpected status %d for POST %s",
		resp.StatusCode,
		url,
	)
}

func freeLocalAddress() (string, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", err
	}
	defer listener.Close()

	return listener.Addr().String(), nil
}

func waitForAccrualServer(t *testing.T, address string) {
	t.Helper()

	client := http.Client{
		Timeout: 200 * time.Millisecond,
	}

	url := "http://" + address + "/api/orders/12345678903"
	deadline := time.Now().Add(5 * time.Second)
	var lastErr error

	for time.Now().Before(deadline) {
		resp, err := client.Get(url)
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent {
				return
			}
		}

		lastErr = err
		time.Sleep(50 * time.Millisecond)
	}

	require.NoError(t, lastErr, "accrual server did not become ready in time")
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
