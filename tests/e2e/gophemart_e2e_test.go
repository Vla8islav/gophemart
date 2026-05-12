package e2e

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPing(t *testing.T) {
	t.Parallel()

	cfg := initE2ETestServer(t)

	pingURL := "http://" + cfg.ServerAddress.Value + "/api/ping"

	resp, err := http.Get(pingURL)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}
