package order_service

import (
	"github.com/gin-gonic/gin"
	"github.com/justin831201/trading-service/app/order-service/internal"
)

func RegisterRoutes(engine *gin.Engine) {
	engine.GET("/", internal.GetIndex)
	engine.POST("/order", internal.SendOrder)
}
