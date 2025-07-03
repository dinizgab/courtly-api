package entity

import "time"

type Payment struct {
	ID                string    `json:"id"`
	BookingID         string    `json:"booking_id"`
	CompanyID         string    `json:"company_id"`
	SubaccountID      string    `json:"subaccount_id"`
	CorrelationID     string    `json:"correlation_id"`
	ChargeID          string    `json:"charge_id"`
	BrCode            string    `json:"brcode"`
	QrCodeImage       string    `json:"qr_code_image"`
	ValueTotal        int64     `json:"value_total"`
	ValueCommission   int64     `json:"value_commission"`
	ValueCompany      int64     `json:"value_company"`
	Status            string    `json:"status"`
	ExpiresAt         time.Time `json:"expires_at"`
	PaidAt            time.Time `json:"paid_at,omitempty"`
	RefundedAt        time.Time `json:"refunded_at,omitempty"`
	RefundRequestedAt time.Time `json:"refund_requested_at,omitempty"`
	RefundEndToEndID  string    `json:"refund_end_to_end_id,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}

type Subaccount struct {
	ID           string    `json:"id"`
	CompanyID    string    `json:"company_id"`
	SubaccountID string    `json:"subaccount_id"`
	PixKey       string    `json:"pix_key"`
    PixKeyType   string    `json:"pix_key_type"`
	CreatedAt    time.Time `json:"created_at"`
}
