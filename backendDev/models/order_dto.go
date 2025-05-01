package models

type LiveOrderDTO struct {
	ID        uint    `json:"id"`
	Symbol    string  `json:"symbol"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Direction string  `json:"direction"`
	Status    string  `json:"status"`
}

type DraftOrderDTO struct {
	ID       uint    `json:"id"`
	Symbol   string  `json:"symbol"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	Status   string  `json:"status"`
}
