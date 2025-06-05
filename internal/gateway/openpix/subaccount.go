package openpix

type Subaccount struct {
	Name    string `json:"name"`
	PixKey  string `json:"pixKey"`
	Balance int64  `json:"balance,omitempty"`
}

type Withdraw struct {
    ID string `json:"id"`
    CompanyId string `json:"companyId"`
    Value int64 `json:"value"`
    CorrelationId string `json:"correlationID"`
    DestinationAlias string `json:"destinationAlias"`
    CreatedAt string `json:"createdAt"`
}

type SubAccountResponse struct {
	Subaccount Subaccount `json:"SubAccount"`
}
