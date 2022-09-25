package main

import (
	"github.com/gin-gonic/gin"
	order_service "github.com/justin831201/trading-service/app/order-service"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	engine := gin.New()
	order_service.RegisterRoutes(engine)
	err := engine.Run(":8080")
	if err != nil {
		return
	}
}
