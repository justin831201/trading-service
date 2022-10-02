package order_model

type RequestOrder struct {
	UserId    string `json:"user_id"`
	OrderType uint8  `json:"order_type"`
	Quantity  int    `json:"quantity"`
	PriceType uint8  `json:"price_type"`
	Price     int    `json:"price,omitempty"`
}
