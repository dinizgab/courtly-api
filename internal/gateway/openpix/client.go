package openpix

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dinizgab/booking-mvp/internal/config"
)

type OpenPixClient interface {
	CreateSubaccount(ctx context.Context, in CreateSubAccountRequest) (Subaccount, error)
}

type openPixClientImpl struct {
	baseURL    string
	appId      string
	httpClient *http.Client
}

func NewOpenPixClient(cfg *config.OpenPixConfig) OpenPixClient {
	return &openPixClientImpl{
		baseURL:    cfg.BaseURL,
		appId:      cfg.AppID,
		httpClient: &http.Client{},
	}
}

func (c *openPixClientImpl) CreateSubaccount(ctx context.Context, in CreateSubAccountRequest) (Subaccount, error) {
	body, err := json.Marshal(in)
	if err != nil {
		return Subaccount{}, fmt.Errorf("failed to marshal request: %w", err)
	}

    fmt.Println(c.appId)

	url := fmt.Sprintf("%s/api/v1/subaccount", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return Subaccount{}, fmt.Errorf("OpenPixClient.CreateSubaccount - failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.appId)

    fmt.Println(req.Header)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return Subaccount{}, fmt.Errorf("OpenPixClient.CreateSubaccount - failed to send request: %w", err)
	}
	defer res.Body.Close()

    fmt.Println("Response Status:", res)
	if res.StatusCode != http.StatusOK {
		return Subaccount{}, fmt.Errorf("OpenPixClient.CreateSubaccount - failed to create subaccount with status: %s", res.Status)
	}

	var out CreateSubAccountResponse
	err = json.NewDecoder(res.Body).Decode(&out)
	if err != nil {
		return Subaccount{}, fmt.Errorf("OpenPixClient.CreateSubaccount - failed to decode response: %w", err)
	}

	return out.Subaccount, err
}
