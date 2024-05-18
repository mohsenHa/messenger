package rabbitmq

type Message struct {
	Ack  func() error
	Body []byte
}
