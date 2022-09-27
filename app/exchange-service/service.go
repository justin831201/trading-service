package exchange_service

import (
	"github.com/justin831201/trading-service/app/exchange-service/internal/env"
	"github.com/justin831201/trading-service/pkg/kafkautil"
	"github.com/justin831201/trading-service/pkg/logger"
	"github.com/justin831201/trading-service/pkg/model"
	log "github.com/sirupsen/logrus"
)

const (
	logPrefix = "[ExchangeService]"
)

var consumer *kafkautil.Consumer

func StartService(configFile string) {

	config := env.LoadConfig(configFile)
	logger.SetupLogger(config.Logger)

	log.Infof("%s Service Start", logPrefix)

	consumer = &kafkautil.Consumer{}
	if err := consumer.Start(config.Kafka, config.Kafka.Topic[model.OrderEventTopicKey], processMessage); err != nil {
		log.Errorf("%s Kafka consumer error %v", logPrefix, err)
	}
}

func DisposeService() {
	log.Infof("%s Service Close", logPrefix)
	consumer.Close()
}
