package main

import (
	"flag"
	service "github.com/justin831201/trading-service/app/order-service"
	"os"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "c", "", "Configuration file path.")
	flag.Parse()

	if configFile == "" {
		flag.Usage()
		os.Exit(0)
	}

	defer service.DisposeService()
	service.StartService(configFile)
}
