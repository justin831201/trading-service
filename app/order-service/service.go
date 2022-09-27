package order_service

import (
	"github.com/gin-gonic/gin"
	"github.com/justin831201/trading-service/app/order-service/internal/controller"
	"github.com/justin831201/trading-service/app/order-service/internal/env"
	sender "github.com/justin831201/trading-service/app/order-service/internal/order-sender"
	"github.com/justin831201/trading-service/pkg/logger"
	log "github.com/sirupsen/logrus"
)

const (
	logPrefix = "[OrderService]"
)

// orderSender can inject mock sender for testing
var orderSender sender.Interface

func registerRoutes(engine *gin.Engine) {
	engine.GET("/", controller.GetIndex)

	orderController := controller.NewOrderController(orderSender)
	engine.POST("/order", orderController.SendOrder)
}

func StartService(configFile string) {
	config := env.LoadConfig(configFile)
	logger.SetupLogger(config.Logger)

	log.Infof("%s Service Start", logPrefix)

	if orderSender == nil {
		orderSender = sender.NewKafkaOrderSender(config)
	}
	engine := gin.Default()
	registerRoutes(engine)
	if err := engine.Run(config.Address); err != nil {
		panic(err)
	}
}

func DisposeService() {
	log.Infof("%s Service Close", logPrefix)

	if p, ok := orderSender.(*sender.KafkaOrderSender); ok {
		p.Close()
	}
}
