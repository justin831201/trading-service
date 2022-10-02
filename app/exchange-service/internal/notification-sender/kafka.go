package notification_sender

import log "github.com/sirupsen/logrus"

const logPrefixKafka = "[KafkaNotificationSender]"

type KafkaNotificationSender struct{}

func (sender *KafkaNotificationSender) SendNotification(orderId string, quantity, price int) {
	// TODO: implement notification service and update this function
	log.Debugf("%s Order %s trade %d quantity successfully on price %d", logPrefixKafka, orderId, quantity, price)
}
