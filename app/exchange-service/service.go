package exchange_service

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/justin831201/trading-service/app/exchange-service/internal/controller"
	"github.com/justin831201/trading-service/app/exchange-service/internal/env"
	notification_sender "github.com/justin831201/trading-service/app/exchange-service/internal/notification-sender"
	"github.com/justin831201/trading-service/pkg/collection"
	"github.com/justin831201/trading-service/pkg/kafkautil"
	"github.com/justin831201/trading-service/pkg/logger"
	"github.com/justin831201/trading-service/pkg/model"
	procedure_group "github.com/justin831201/trading-service/pkg/procedure-group"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	logPrefix        = "[ExchangeService]"
	exchangeDuration = 5 * time.Second
)

func processMessageFunc(orders *collection.ThreadSafeArray) func(message *kafka.Message) error {
	return func(message *kafka.Message) error {
		order := new(model.Order)
		if err := json.Unmarshal(message.Value, order); err != nil {
			return err
		}
		if order.PriceType < 0 && order.PriceType > 1 {
			log.Debugf("%s Order %s has unknown price type %d", logPrefix, order.OrderId, order.PriceType)
			return nil
		}
		if order.OrderType < 0 && order.OrderType > 1 {
			log.Debugf("%s Order %s has unknown order type %d", logPrefix, order.OrderId, order.OrderType)
			return nil
		}

		orders.StartTransaction()
		defer orders.EndTransaction()
		orders.Add(order)
		return nil
	}
}

func StartService(configFile string) {

	config := env.LoadConfig(configFile)
	logger.SetupLogger(config.Logger)

	log.Infof("%s Service Start", logPrefix)

	orders := &collection.ThreadSafeArray{}

	procedureGroup := procedure_group.ProcedureGroup{}

	var consumer *kafkautil.Consumer
	procedureGroup.Add(func() error {
		consumer = &kafkautil.Consumer{}
		return consumer.Start(
			config.Kafka,
			config.Kafka.Topic[model.OrderEventTopicKey],
			processMessageFunc(orders),
		)
	}, func(err error) {
		consumer.Close()
	})

	procedureGroup.Add(func() error {
		kafkaNotificationSender := &notification_sender.KafkaNotificationSender{}
		exchangeController := controller.NewExchangeController(orders, kafkaNotificationSender)
		for true {
			time.Sleep(exchangeDuration)
			exchangeController.Exchange()
		}
		return nil
	}, func(err error) {
		return
	})

	_ = procedureGroup.Run()
}
