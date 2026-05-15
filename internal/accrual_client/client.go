package accrual_client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Vla8islav/gophemart/internal/domain"
)

type AccrualClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewAccrualClient(baseURL string, httpClient *http.Client) *AccrualClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL = strings.TrimRight(baseURL, "/")
	if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
		baseURL = "http://" + baseURL
	}

	return &AccrualClient{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

func (c *AccrualClient) GetOrderInfo(ctx context.Context, orderNumber string) (*domain.AccrualOrderInfoResponse, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.baseURL+"/api/orders/"+orderNumber,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create accrual request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call accrual service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected accrual status code: %d", resp.StatusCode)
	}

	var result domain.AccrualOrderInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode accrual response: %w", err)
	}

	return &result, nil
}
