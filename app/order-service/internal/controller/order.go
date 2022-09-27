package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/justin831201/trading-service/app/order-service/internal/model"
	order_producer "github.com/justin831201/trading-service/app/order-service/internal/order-sender"
	"github.com/justin831201/trading-service/pkg/model"
	"github.com/sirupsen/logrus"
	"net/http"
)

const logPrefixOrder = "[OrderController]"

type OrderController struct {
	orderProducer order_producer.Interface
}

func NewOrderController(orderProducer order_producer.Interface) OrderController {
	return OrderController{orderProducer: orderProducer}
}

func (controller *OrderController) SendOrder(ctx *gin.Context) {
	requestOrder := order_model.RequestOrder{}
	if err := ctx.ShouldBindJSON(&requestOrder); err != nil {
		logrus.Debug(err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, order_model.ErrorMessage{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	orderId := uuid.NewString()

	order := &model.Order{
		OrderId:   orderId,
		UserId:    requestOrder.UserId,
		OrderType: requestOrder.OrderType,
		Quantity:  requestOrder.Quantity,
		PriceType: requestOrder.PriceType,
	}

	var err error
	switch requestOrder.OrderType {
	case "buy", "sell":
		err = controller.orderProducer.SendOrder(order)
	default:
		logrus.Tracef("%s Unknown order type %s", logPrefixOrder, requestOrder.OrderType)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, order_model.ErrorMessage{
			Status:  http.StatusBadRequest,
			Message: "invalid order type",
		})
		return
	}

	if err != nil {
		logrus.Warningf("%s Sending order error: %v", logPrefixOrder, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, order_model.ErrorMessage{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"order_id": order.OrderId,
	})
}
