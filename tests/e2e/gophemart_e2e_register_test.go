package e2e

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPing(t *testing.T) {

	cfg := initE2ETestServer(t)

	pingURL := "http://" + cfg.ServerAddress.Value + "/api/ping"

	resp, err := http.Get(pingURL)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRegisterUserDuplicateLogin(t *testing.T) {

	cfg := initE2ETestServer(t)

	registerURL := "http://" + cfg.ServerAddress.Value + "/api/user/register"
	login := fmt.Sprintf("e2e-user-%d", time.Now().UnixNano())
	body := []byte(fmt.Sprintf(`{"login":%q,"password":"test-password"}`, login))

	firstResp, err := http.Post(registerURL, "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer firstResp.Body.Close()
	require.Equal(t, http.StatusOK, firstResp.StatusCode)
	require.NotEmpty(t, firstResp.Cookies(), "successful registration should set auth cookie")

	secondResp, err := http.Post(registerURL, "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer secondResp.Body.Close()
	require.Equal(t, http.StatusConflict, secondResp.StatusCode)
}

func TestRegisterUserRejectsEmptyPassword(t *testing.T) {

	cfg := initE2ETestServer(t)

	registerURL := "http://" + cfg.ServerAddress.Value + "/api/user/register"
	login := fmt.Sprintf("e2e-user-empty-password-%d", time.Now().UnixNano())
	body := []byte(fmt.Sprintf(`{"login":%q,"password":""}`, login))

	resp, err := http.Post(registerURL, "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
