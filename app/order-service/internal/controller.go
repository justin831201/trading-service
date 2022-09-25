package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/justin831201/trading-service/app/order-service/internal/model"
	"net/http"
)

func GetIndex(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func SendOrder(ctx *gin.Context) {
	requestOrder := model.RequestOrder{}
	if err := ctx.ShouldBindJSON(&requestOrder); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorMessage{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	// TODO: implement trading algorithm

	ctx.JSON(http.StatusOK, gin.H{
		"order_id": requestOrder.OrderId,
	})
}
