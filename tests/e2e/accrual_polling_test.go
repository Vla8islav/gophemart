package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Vla8islav/gophemart/internal/helpers"
	"github.com/stretchr/testify/require"
)

type accrualPollingOrderResponse struct {
	Number  string   `json:"number"`
	Status  string   `json:"status"`
	Accrual *float64 `json:"accrual,omitempty"`
}

func TestAccrualPollingUpdatesOrders(t *testing.T) {

	cfg := initE2ETestServer(t)

	user := map[string]string{
		"login":    helpers.UniqueLogin("accrual-polling-user"),
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

	authCookie := registerResp.Cookies()[0]

	orderNumbers := []string{
		"9278923470",
		"12345678903",
		"346436439",
	}

	client := http.Client{Timeout: 2 * time.Second}

	for _, orderNumber := range orderNumbers {
		req, err := http.NewRequest(http.MethodPost,
			"http://"+cfg.ServerAddress.Value+"/api/user/orders",
			bytes.NewReader([]byte(orderNumber)),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "text/plain")
		req.AddCookie(authCookie)

		resp, err := client.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusAccepted, resp.StatusCode)
		require.NoError(t, resp.Body.Close())
	}

	initialOrders := getUserOrders(t, client, cfg.ServerAddress.Value, authCookie)
	initialOrdersByNumber := ordersByNumber(initialOrders)

	for _, orderNumber := range orderNumbers {
		order, ok := initialOrdersByNumber[orderNumber]
		require.True(t, ok)
		require.Equal(t, "NEW", order.Status)
		require.Nil(t, order.Accrual)
	}

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	timeout := time.After(30 * time.Second)

	for {
		select {
		case <-timeout:
			lastOrders := getUserOrders(t, client, cfg.ServerAddress.Value, authCookie)
			require.Failf(t,
				"orders were not updated by accrual polling within 30 seconds",
				"initial orders: %+v\nlast orders: %+v",
				initialOrders,
				lastOrders,
			)
			return
		case <-ticker.C:
			currentOrders := getUserOrders(t, client, cfg.ServerAddress.Value, authCookie)
			currentOrdersByNumber := ordersByNumber(currentOrders)

			for _, orderNumber := range orderNumbers {
				initialOrder := initialOrdersByNumber[orderNumber]
				currentOrder, ok := currentOrdersByNumber[orderNumber]
				require.True(t, ok)

				if currentOrder.Status != initialOrder.Status || !sameAccrual(currentOrder.Accrual, initialOrder.Accrual) {
					return
				}
			}
		}
	}
}

func getUserOrders(t *testing.T, client http.Client, serverAddress string, authCookie *http.Cookie) []accrualPollingOrderResponse {
	t.Helper()

	req, err := http.NewRequest(http.MethodGet,
		"http://"+serverAddress+"/api/user/orders",
		nil,
	)
	require.NoError(t, err)
	req.AddCookie(authCookie)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var orders []accrualPollingOrderResponse
	err = json.NewDecoder(resp.Body).Decode(&orders)
	require.NoError(t, err)

	return orders
}

func ordersByNumber(orders []accrualPollingOrderResponse) map[string]accrualPollingOrderResponse {
	result := make(map[string]accrualPollingOrderResponse, len(orders))
	for _, order := range orders {
		result[order.Number] = order
	}

	return result
}

func sameAccrual(left *float64, right *float64) bool {
	if left == nil || right == nil {
		return left == right
	}

	return *left == *right
}
