package exchange_service

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/justin831201/trading-service/pkg/model"
	log "github.com/sirupsen/logrus"
)

func processMessage(message *kafka.Message) error {
	order := new(model.Order)
	if err := json.Unmarshal(message.Value, order); err != nil {
		return err
	}
	log.Debugf("%s Receive order event from kafka: %v", logPrefix, order)
	return nil
}
