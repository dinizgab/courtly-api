package openpix

type Refund struct {
    ID                string `json:"id"`
    EndToEndID         string `json:"transactionEndToEndId"`
    CorrelationID    string `json:"correlationId"`
    RefundedAt       string `json:"time"`
    Value            int64  `json:"value"`
    Status           string `json:"status"`
}
