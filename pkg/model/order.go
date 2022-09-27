package model

const (
	OrderEventTopicKey = "order_event"
)

type Order struct {
	OrderId   string `json:"order_id"`
	UserId    string `json:"user_id"`
	Quantity  int    `json:"quantity"`
	OrderType string `json:"order_type"`
	PriceType string `json:"price_type"`
	Price     int    `json:"price,omitempty"`
}
