package controller

import (
	notification_sender "github.com/justin831201/trading-service/app/exchange-service/internal/notification-sender"
	"github.com/justin831201/trading-service/pkg/collection"
	"github.com/justin831201/trading-service/pkg/model"
	log "github.com/sirupsen/logrus"
	"sort"
)

const logPrefixExchange = "[ExchangeController]"

type ExchangeController struct {
	orders             *collection.ThreadSafeArray
	notificationSender notification_sender.Interface
}

func NewExchangeController(orders *collection.ThreadSafeArray, notificationSender notification_sender.Interface) ExchangeController {
	return ExchangeController{
		orders:             orders,
		notificationSender: notificationSender,
	}
}

// Exchange The function would exchange orders
func (controller *ExchangeController) Exchange() {
	/*
	 * The function would exchange orders, and it can separate into four steps
	 * Step1: Separate orders into different arrays by order type and price type
	 * Step2: Sort limit price order arrays
	 * Step3: handle market price orders
	 * Step4: handle limit price orders
	 * Step5: remove order which quantity is 0 after exchanging
	 */
	log.Debugf("%s Start to exchange...", logPrefixExchange)
	controller.orders.StartTransaction()
	defer controller.orders.EndTransaction()
	var marketPriceBuyOrders, marketPriceSellOrders, limitPriceBuyOrders, limitPriceSellOrders []*model.Order

	// Step1
	// separate orders into different arrays by order type and price type
	for idx := 0; idx < controller.orders.Size(); idx++ {
		order := controller.orders.Get(idx).(*model.Order)
		log.Debugf("%s Order %v", logPrefixExchange, *order)
		if order.OrderType == 0 {
			if order.PriceType == 0 {
				marketPriceBuyOrders = append(marketPriceBuyOrders, order)
			} else {
				limitPriceBuyOrders = append(limitPriceBuyOrders, order)
			}
		} else {
			if order.PriceType == 0 {
				marketPriceSellOrders = append(marketPriceSellOrders, order)
			} else {
				limitPriceSellOrders = append(limitPriceSellOrders, order)
			}
		}
	}

	// Step2
	// descending sort limit price buy orders
	// Example [100, 99, 98, 97...]
	sort.Slice(limitPriceBuyOrders, func(i, j int) bool {
		return limitPriceBuyOrders[i].Price > limitPriceBuyOrders[j].Price
	})
	// ascending sort limit price sell orders
	// Example [103, 102, 101, 100...]
	sort.Slice(limitPriceSellOrders, func(i, j int) bool {
		return limitPriceSellOrders[i].Price < limitPriceSellOrders[j].Price
	})

	var marketPriceBuyOrderIdx, marketPriceSellOrderIdx, limitPriceBuyOrderIdx, limitPriceSellOrderIdx int

	// Step3
	// handle market price orders
	for marketPriceBuyOrderIdx < len(marketPriceBuyOrders) {
		if limitPriceSellOrderIdx >= len(limitPriceSellOrders) {
			break
		}

		quantity := min(marketPriceBuyOrders[marketPriceBuyOrderIdx].Quantity, limitPriceSellOrders[limitPriceSellOrderIdx].Quantity)
		price := limitPriceSellOrders[limitPriceSellOrderIdx].Price

		marketPriceBuyOrders[marketPriceBuyOrderIdx].Quantity -= quantity
		limitPriceSellOrders[limitPriceSellOrderIdx].Quantity -= quantity

		controller.notificationSender.SendNotification(marketPriceBuyOrders[marketPriceBuyOrderIdx].OrderId, quantity, price)
		controller.notificationSender.SendNotification(limitPriceSellOrders[limitPriceSellOrderIdx].OrderId, quantity, price)

		if marketPriceBuyOrders[marketPriceBuyOrderIdx].Quantity == 0 {
			marketPriceBuyOrderIdx++
		}
		if limitPriceSellOrders[limitPriceSellOrderIdx].Quantity == 0 {
			limitPriceSellOrderIdx++
		}
	}

	for marketPriceSellOrderIdx < len(marketPriceSellOrders) {
		if limitPriceBuyOrderIdx >= len(limitPriceBuyOrders) {
			break
		}

		quantity := min(marketPriceSellOrders[marketPriceSellOrderIdx].Quantity, limitPriceBuyOrders[limitPriceBuyOrderIdx].Quantity)
		price := limitPriceBuyOrders[limitPriceBuyOrderIdx].Price

		marketPriceSellOrders[marketPriceSellOrderIdx].Quantity -= quantity
		limitPriceBuyOrders[limitPriceBuyOrderIdx].Quantity -= quantity

		controller.notificationSender.SendNotification(marketPriceSellOrders[marketPriceSellOrderIdx].OrderId, quantity, price)
		controller.notificationSender.SendNotification(limitPriceBuyOrders[limitPriceBuyOrderIdx].OrderId, quantity, price)

		if marketPriceSellOrders[marketPriceSellOrderIdx].Quantity == 0 {
			marketPriceSellOrderIdx++
		}
		if limitPriceBuyOrders[limitPriceBuyOrderIdx].Quantity == 0 {
			limitPriceBuyOrderIdx++
		}
	}

	// Step4
	// handle limit price orders
	for limitPriceBuyOrderIdx < len(limitPriceBuyOrders) && limitPriceSellOrderIdx < len(limitPriceSellOrders) {
		if limitPriceBuyOrders[limitPriceBuyOrderIdx].Price < limitPriceSellOrders[limitPriceSellOrderIdx].Price {
			break
		}

		quantity := min(limitPriceBuyOrders[limitPriceBuyOrderIdx].Quantity, limitPriceSellOrders[limitPriceSellOrderIdx].Quantity)
		price := (limitPriceBuyOrders[limitPriceBuyOrderIdx].Price + limitPriceSellOrders[limitPriceSellOrderIdx].Price) / 2

		limitPriceBuyOrders[limitPriceBuyOrderIdx].Quantity -= quantity
		limitPriceSellOrders[limitPriceSellOrderIdx].Quantity -= quantity

		controller.notificationSender.SendNotification(limitPriceBuyOrders[limitPriceBuyOrderIdx].OrderId, quantity, price)
		controller.notificationSender.SendNotification(limitPriceSellOrders[limitPriceSellOrderIdx].OrderId, quantity, price)

		if limitPriceBuyOrders[limitPriceBuyOrderIdx].Quantity == 0 {
			limitPriceBuyOrderIdx++
		}
		if limitPriceSellOrders[limitPriceSellOrderIdx].Quantity == 0 {
			limitPriceSellOrderIdx++
		}
	}

	// Step5
	// remove order which quantity is 0 after exchanging
	var idx int
	for idx < controller.orders.Size() {
		order := controller.orders.Get(idx).(*model.Order)
		if order.Quantity == 0 {
			controller.orders.Remove(idx)
		} else {
			idx++
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
