package order_model

type ErrorMessage struct {
	Status  int    `json:"status"`  // Ths status value would be same as status code
	Message string `json:"message"` // The message provides information about the current situation.
}
