package model

type RequestOrder struct {
	UserId    string `json:"user_id"`
	OrderType string `json:"order_type"`
	Quantity  int    `json:"quantity"`
	PriceType string `json:"price_type"`
	Price     int    `json:"price,omitempty"`
}
