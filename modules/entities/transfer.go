package entities

type TransferRequest struct {
	Type   string  `json:"type"`
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
}
