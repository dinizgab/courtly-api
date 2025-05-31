package openpix

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"

	"github.com/dinizgab/booking-mvp/internal/config"
	"github.com/dinizgab/booking-mvp/internal/entity"
)

type OpenPixClient interface {
	CreateSubaccount(ctx context.Context, in CreateSubAccountRequest) (Subaccount, error)
	CreateCharge(ctx context.Context, subaccountKey string, booking entity.Booking) (Charge, error)
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

// TODO - Change parameters to receive a Subaccount entity
func (c *openPixClientImpl) CreateSubaccount(ctx context.Context, in CreateSubAccountRequest) (Subaccount, error) {
	body, err := json.Marshal(in)
	if err != nil {
		return Subaccount{}, fmt.Errorf("OpenPixClient.CreateSubaccount - failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/subaccount", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return Subaccount{}, fmt.Errorf("OpenPixClient.CreateSubaccount - failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.appId)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return Subaccount{}, fmt.Errorf("OpenPixClient.CreateSubaccount - failed to send request: %w", err)
	}
	defer res.Body.Close()

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

func (c *openPixClientImpl) CreateCharge(ctx context.Context, subaccountKey string, booking entity.Booking) (Charge, error) {
	correlationId := fmt.Sprintf("booking-%s", booking.ID)
	parsedPrice := int64(math.Round(booking.TotalPrice * 100))
	in := CreateChargeRequest{
		CorrelationID: correlationId,
		Value:         parsedPrice,
		Customer: Customer{
			Name:  booking.GuestName,
			Email: booking.GuestEmail,
			Phone: booking.GuestPhone,
		},
		Splits: []Split{{
			Value:     parsedPrice - 300,
			PixKey:    subaccountKey,
			SplitType: "SPLIT_SUB_ACCOUNT",
		}},
		ExpiresIn: 1800,
	}

	body, err := json.Marshal(in)
	if err != nil {
		return Charge{}, fmt.Errorf("OpenPixClient.CreateCharge - failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/charge", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return Charge{}, fmt.Errorf("OpenPixClient.CreateCharge - failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.appId)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return Charge{}, fmt.Errorf("OpenPixClient.CreateCharge- failed to send request: %w", err)
	}
	defer res.Body.Close()

	var out CreateChargeResponse
	err = json.NewDecoder(res.Body).Decode(&out)
	if err != nil {
		return Charge{}, fmt.Errorf("OpenPixClient.CreateCharge - failed to decode response: %w", err)
	}

	return out.Charge, nil
}
