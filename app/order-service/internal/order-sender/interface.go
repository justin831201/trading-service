package order_sender

import "github.com/justin831201/trading-service/pkg/model"

type Interface interface {
	SendOrder(order *model.Order) error
}
