package kafkautil

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	log "github.com/sirupsen/logrus"
	"strings"
)

const logPrefixConsumer = "[KafkaConsumer]"

type Consumer struct {
	c         *kafka.Consumer
	isRunning bool
}

func (consumer *Consumer) Start(config *Config, topic string, processMessage func(msg *kafka.Message) error) error {
	if consumer.c != nil {
		log.Warningf("%s Already running, operation cancelled", logPrefixConsumer)
		return nil
	}

	var err error

	consumer.c, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        strings.Join(config.BootstrapServers, ","),
		"group.id":                 config.Consumer.GroupId,
		"enable.auto.commit":       true,
		"auto.offset.reset":        "earliest",
		"go.events.channel.enable": false,
	})

	if err != nil {
		log.Errorf("%s Initialize Kafka consumer error: %v", logPrefixConsumer, err)
		panic(err)
	}

	log.Infof("%s Start consuming...", logPrefixConsumer)

	err = consumer.c.Subscribe(topic, nil)
	if err != nil {
		log.Errorf("%s Subscribe topic %s error: %v", logPrefixConsumer, topic, err)
		panic(err)
	}

	consumer.isRunning = true
	for consumer.isRunning {
		e := consumer.c.Poll(1000)
		switch event := e.(type) {
		case *kafka.Message:
			log.Debugf("%s Message on %s: %s", logPrefixConsumer, event.TopicPartition, string(event.Value))
			if err := processMessage(event); err != nil {
				log.Debugf("%s Process message error %v", logPrefixConsumer, err)
			}
		case kafka.Error:
			log.Debugf("%s Error: %v (%v)", logPrefixConsumer, err, e)

			// Terminate the application if all brokers are down
			if event.Code() == kafka.ErrAllBrokersDown {
				consumer.isRunning = false
			}
		default:
			log.Tracef("%s Ignored %v", logPrefixConsumer, event)
		}
	}

	return err
}

func (consumer *Consumer) Close() {
	log.Infof("%s Closing %v...", logPrefixConsumer, consumer.c)
	consumer.isRunning = false
	if consumer.c != nil {
		_ = consumer.c.Close()
		log.Infof("%s Closed", logPrefixConsumer)
	} else {
		log.Warningf("%s Unable to close nil consumer", logPrefixConsumer)
	}
}
