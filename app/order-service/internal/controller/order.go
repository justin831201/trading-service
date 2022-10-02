package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/justin831201/trading-service/app/order-service/internal/model"
	order_sender "github.com/justin831201/trading-service/app/order-service/internal/order-sender"
	"github.com/justin831201/trading-service/pkg/model"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const logPrefixOrder = "[OrderController]"

type OrderController struct {
	orderProducer order_sender.Interface
}

func NewOrderController(orderProducer order_sender.Interface) OrderController {
	return OrderController{orderProducer: orderProducer}
}

func (controller *OrderController) SendOrder(ctx *gin.Context) {
	requestOrder := order_model.RequestOrder{}
	if err := ctx.ShouldBindJSON(&requestOrder); err != nil {
		log.Debug(err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, order_model.ErrorMessage{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	if requestOrder.OrderType != 0 && requestOrder.OrderType != 1 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, order_model.ErrorMessage{
			Status:  http.StatusBadRequest,
			Message: "invalid order type",
		})
		return
	}
	if requestOrder.PriceType != 0 && (requestOrder.PriceType != 1) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, order_model.ErrorMessage{
			Status:  http.StatusBadRequest,
			Message: "invalid price type",
		})
		return
	}

	if requestOrder.PriceType == 1 && requestOrder.Price == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, order_model.ErrorMessage{
			Status:  http.StatusBadRequest,
			Message: "price needed when price type is limit price",
		})
		return
	}

	orderId := uuid.NewString()

	order := &model.Order{
		OrderId:   orderId,
		UserId:    requestOrder.UserId,
		OrderType: requestOrder.OrderType,
		Quantity:  requestOrder.Quantity,
		PriceType: requestOrder.PriceType,
		Price:     requestOrder.Price,
	}

	if err := controller.orderProducer.SendOrder(order); err != nil {
		log.Warningf("%s Sending order error: %v", logPrefixOrder, err)
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
