package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	order_service "github.com/justin831201/trading-service/app/order-service"
	"github.com/justin831201/trading-service/pkg/logger"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "c", "", "Configuration file path.")
	flag.Parse()

	config := order_service.LoadConfig(configFile)

	logger.SetupLogger(config.Logger)

	engine := gin.New()
	order_service.RegisterRoutes(engine)
	err := engine.Run(":8080")
	if err != nil {
		return
	}
}
