package order_sender

import (
	"encoding/json"
	"fmt"
	"github.com/justin831201/trading-service/app/order-service/internal/env"
	"github.com/justin831201/trading-service/pkg/kafkautil"
	"github.com/justin831201/trading-service/pkg/model"
)

type KafkaOrderSender struct {
	orderTopic string
	producer   *kafkautil.Producer
}

func NewKafkaOrderSender(config *env.Config) *KafkaOrderSender {
	producer := new(kafkautil.Producer)
	if err := producer.Init(config.Kafka); err != nil {
		panic(err)
	} else if orderEventTopic, ok := config.Kafka.Topic[model.OrderEventTopicKey]; !ok {
		panic(fmt.Sprintf("unknown order topic key %s", model.OrderEventTopicKey))
	} else {
		return &KafkaOrderSender{
			orderTopic: orderEventTopic,
			producer:   producer,
		}
	}
}

func (sender *KafkaOrderSender) SendOrder(order *model.Order) error {
	if msg, err := json.Marshal(order); err != nil {
		return err
	} else {
		return sender.producer.SendMessage(sender.orderTopic, msg)
	}
}

func (sender *KafkaOrderSender) Close() {
	sender.producer.Close()
}
