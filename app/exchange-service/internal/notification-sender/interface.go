package notification_sender

type Interface interface {
	SendNotification(orderId string, quantity, price int)
}
