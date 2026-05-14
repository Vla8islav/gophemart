package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
	t.Parallel()

	cfg := initE2ETestServer(t)

	user := map[string]string{
		"login":    "login-user",
		"password": "password",
	}

	registerBody, err := json.Marshal(user)
	require.NoError(t, err)

	registerResp, err := http.Post("http://"+cfg.ServerAddress.Value+"/api/user/register",
		"application/json", bytes.NewReader(registerBody))
	require.NoError(t, err)
	defer registerResp.Body.Close()

	require.Equal(t, http.StatusOK, registerResp.StatusCode)
	require.NotEmpty(t, registerResp.Cookies())

	loginBody, err := json.Marshal(user)
	require.NoError(t, err)

	loginResp, err := http.Post("http://"+cfg.ServerAddress.Value+"/api/user/login",
		"application/json", bytes.NewReader(loginBody))
	require.NoError(t, err)
	defer loginResp.Body.Close()

	require.Equal(t, http.StatusOK, loginResp.StatusCode)
	require.NotEmpty(t, loginResp.Cookies())
}
