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
	CreateSubaccount(ctx context.Context, subaccount Subaccount) (Subaccount, error)
	CreateCharge(ctx context.Context, subaccountKey string, booking entity.Booking) (Charge, error)
	GetCompanyBalance(ctx context.Context, pixKey string) (int64, error)
	WithdrawSubaccount(ctx context.Context, pixKey string) (Withdraw, error)
    RefundCharge(ctx context.Context, payment entity.Payment) (Refund, error)
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

func (c *openPixClientImpl) CreateSubaccount(ctx context.Context, subaccount Subaccount) (Subaccount, error) {
	body, err := json.Marshal(subaccount)
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

	var out SubAccountResponse
	err = json.NewDecoder(res.Body).Decode(&out)
	if err != nil {
		return Subaccount{}, fmt.Errorf("OpenPixClient.CreateSubaccount - failed to decode response: %w", err)
	}

	return out.Subaccount, err
}

// TODO - Check charges with a large amounts of money, its giving an error with split
// {"error":"O valor total do split de pagamento não pode ser igual ou maior que o valor da cobrança menos a taxa esperada"}
func (c *openPixClientImpl) CreateCharge(ctx context.Context, subaccountKey string, booking entity.Booking) (Charge, error) {
	correlationId := fmt.Sprintf("booking-%s", booking.ID)
	totalPrice := int64(math.Round(booking.TotalPrice * 100))
	gasPrice := (totalPrice * 5 + 50) / 100
	in := CreateChargeRequest{
		CorrelationID: correlationId,
		Value:         totalPrice,
		Customer: Customer{
			Name:  booking.GuestName,
			Email: booking.GuestEmail,
			Phone: booking.GuestPhone,
		},
		Splits: []Split{{
			Value:     totalPrice - gasPrice,
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

	if res.StatusCode != http.StatusOK {
		return Charge{}, fmt.Errorf("OpenPixClient.CreateCharge - failed to create charge with status: %s", res.Status)
	}

	var out CreateChargeResponse
	err = json.NewDecoder(res.Body).Decode(&out)
	if err != nil {
		return Charge{}, fmt.Errorf("OpenPixClient.CreateCharge - failed to decode response: %w", err)
	}

	return out.Charge, nil
}

func (c *openPixClientImpl) GetCompanyBalance(ctx context.Context, pixKey string) (int64, error) {
	url := fmt.Sprintf("%s/api/v1/subaccount/%s", c.baseURL, pixKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("OpenPixClient.GetCompanyBalance - failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.appId)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("OpenPixClient.GetCompanyBalance - failed to send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("OpenPixClient.GetCompanyBalance - failed to get balance with status: %s", res.Status)
	}

	var out SubAccountResponse
	err = json.NewDecoder(res.Body).Decode(&out)
	if err != nil {
		return 0, fmt.Errorf("OpenPixClient.GetCompanyBalance - failed to decode response: %w", err)
	}

	return out.Subaccount.Balance, nil
}

func (c *openPixClientImpl) WithdrawSubaccount(ctx context.Context, pixKey string) (Withdraw, error) {
	url := fmt.Sprintf("%s/api/v1/subaccount/%s/withdraw", c.baseURL, pixKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return Withdraw{}, fmt.Errorf("OpenPixClient.WithdrawSubaccount - failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.appId)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return Withdraw{}, fmt.Errorf("OpenPixClient.WithdrawSubaccount - failed to send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Withdraw{}, fmt.Errorf("OpenPixClient.WithdrawSubaccount - failed to withdraw with status: %s", res.Status)
	}

	var out struct {
		Withdraw Withdraw `json:"transaction"`
	}
	err = json.NewDecoder(res.Body).Decode(&out)
	if err != nil {
		return Withdraw{}, fmt.Errorf("OpenPixClient.WithdrawSubaccount - failed to decode response: %w", err)
	}

	return out.Withdraw, nil
}

func (c *openPixClientImpl) RefundCharge(ctx context.Context, payment entity.Payment) (Refund, error) {
    refundCorrelationID := fmt.Sprintf("refund-%s", payment.ID)
    in := Refund{
        EndToEndID: payment.ID,
        CorrelationID: refundCorrelationID,
    }
	body, err := json.Marshal(in)
	if err != nil {
		return Refund{}, fmt.Errorf("OpenPixClient.CreateSubaccount - failed to marshal request: %w", err)
	}

    url := fmt.Sprintf("%s/api/v1/charge/%s/refund", c.baseURL, payment.CorrelationID)

    req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
    if err != nil {
        return Refund{}, fmt.Errorf("OpenPixClient.RefundCharge - failed to create request: %w", err)
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", c.appId)

    res, err := c.httpClient.Do(req)
    if err != nil {
        return Refund{}, fmt.Errorf("OpenPixClient.RefundCharge - failed to send request: %w", err)
    }
    defer res.Body.Close()

    if res.StatusCode != http.StatusOK {
        return Refund{}, fmt.Errorf("OpenPixClient.RefundCharge - failed to refund with status: %s", res.Status)
    }

    var out struct {
        Refund Refund `json:"refund"`
    }
    err = json.NewDecoder(res.Body).Decode(&out)
    if err != nil {
        return Refund{}, fmt.Errorf("OpenPixClient.RefundCharge - failed to decode response: %w", err)
    }

    return out.Refund, nil
}
