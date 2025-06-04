package openpix

type Subaccount struct {
	Name    string `json:"name"`
	PixKey  string `json:"pixKey"`
	Balance int64  `json:"balance,omitempty"`
}

type CreateSubAccountRequest struct {
	Name   string `json:"name"`
	PixKey string `json:"pixKey"`
}

type SubAccountResponse struct {
	Subaccount Subaccount `json:"SubAccount"`
}
