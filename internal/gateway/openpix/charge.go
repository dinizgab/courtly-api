package openpix

type Customer struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type Split struct {
	Value     int64  `json:"value"`
	PixKey    string `json:"pixKey"`
	SplitType string `json:"splitType"`
}

// TODO - See to add comment to charge
type CreateChargeRequest struct {
	CorrelationID string   `json:"correlationID"`
	Value         int64    `json:"value"`
	Customer      Customer `json:"customer"`
	Splits        []Split  `json:"splits"`
	Subaccount    string   `json:"subaccount"`
	ExpiresIn     int64    `json:"expiresIn"`
}

type CreateChargeResponse struct {
	Charge Charge `json:"charge"`
}

type Charge struct {
	Status         string `json:"status"`
	Value          int64  `json:"value"`
    GasPrice      int64  `json:"gasPrice"`
	CorrelationID  string `json:"correlationID"`
	PaymentLinkID  string `json:"paymentLinkID"`
	PaymentLinkURL string `json:"paymentLinkUrl"`
	QrCodeImage    string `json:"qrCodeImage"`
	ExpiresDate    string `json:"expiresDate"`
	Brcode         string `json:"brCode"`
	PaidAt         string `json:"paidAt"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`

	// TODO - Add Payer information to payment
	//Payer Payer `json:"payer"`
}

type ChargeWebhookEvent struct {
	WebhookEvent  string `json:"event"`
	Charge Charge `json:"charge"`
}
