package openpix

type Refund struct {
	ID            string `json:"id"`
	EndToEndID    string `json:"endToEndId"`
	CorrelationID string `json:"correlationID"`
	RefundedAt    string `json:"time"`
	Value         int64  `json:"value"`
	Status        string `json:"status"`
}
