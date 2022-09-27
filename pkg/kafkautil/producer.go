package kafkautil

import (
	"errors"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	log "github.com/sirupsen/logrus"
	"strings"
)

const logPrefixProducer = "[KafkaProducer]"

type Producer struct {
	kafkaProducer *kafka.Producer
}

func (p *Producer) Init(config *Config) error {
	if p.kafkaProducer != nil {
		return nil
	}

	var err error
	p.kafkaProducer, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":      strings.Join(config.BootstrapServers, ","),
		"compression.codec":      "gzip",
		"queue.buffering.max.ms": 1000,
	})
	if err != nil {
		panic(err)
	}

	go func() {
		for e := range p.kafkaProducer.Events() {
			switch event := e.(type) {
			case *kafka.Message:
				if event.TopicPartition.Error != nil {
					log.Errorf("%s Delivery failed: %v", logPrefixProducer, event.TopicPartition)
				} else {
					log.Debugf("%s Delivery succeed: %v", logPrefixProducer, event.TopicPartition)
				}
			}
		}
	}()

	return err
}

func (p *Producer) SendMessage(topic string, msg []byte) error {
	if p.kafkaProducer == nil {
		return errors.New("sending message with a nil producer")
	}
	log.Debugf("%s Send message: %s", logPrefixProducer, string(msg))
	return p.kafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          msg,
	}, nil)
}

func (p *Producer) Close() {
	p.kafkaProducer.Close()
}
