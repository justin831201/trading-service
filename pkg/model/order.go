package model

const (
	OrderEventTopicKey = "order_event"
)

type Order struct {
	OrderId   string `json:"order_id"`
	UserId    string `json:"user_id"`
	Quantity  int    `json:"quantity"`
	OrderType uint8  `json:"order_type"` // 0 means buy, 1 means sell
	PriceType uint8  `json:"price_type"` // 0 means market price, 1 means limit price
	Price     int    `json:"price,omitempty"`
}
