package interfaces

type PubSubInterface[T any] interface {
	CloseFileWriter() error
	ActiveFileWriter(string) PubSubInterface[T]
	Subscribe(string) <-chan T
	SubscribeWithHandler(topic string, handler func(msg T)) func()
	Unsubscribe(topic string, ch <-chan T)
	Publish(topic string, msg T)
}
