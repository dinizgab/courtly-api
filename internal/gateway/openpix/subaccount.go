package openpix

type Subaccount struct {
	Name    string `json:"name"`
	PixKey  string `json:"pixKey"`
	Balance int64  `json:"balance,omitempty"`
}

type SubAccountResponse struct {
	Subaccount Subaccount `json:"SubAccount"`
}
