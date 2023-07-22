package pubsub

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/IIGabriel/basic-pubsub-lib/interfaces"
)

type pubSub[T any] struct {
	subscribers map[string][]chan T
	mu          sync.RWMutex
	writerConfig
}

type writerConfig struct {
	active     bool
	fileWriter *os.File
}

// NewPubSub initializes and returns a new PubSub instance.
func NewPubSub[T any]() interfaces.PubSubInterface[T] {
	return &pubSub[T]{
		subscribers: make(map[string][]chan T),
	}
}

// ActiveFileWriter activates a file writer for the PubSub with the given file name.
// Will panic if the file fails to open.
func (ps *pubSub[T]) ActiveFileWriter(fileName string) interfaces.PubSubInterface[T] {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to open pubsub: %v", err))
	}

	ps.writerConfig = writerConfig{true, file}

	return ps
}

// CloseFileWriter deactivates the file writer and closes the file.
// Returns any error encountered during file closing.
func (ps *pubSub[T]) CloseFileWriter() error {
	ps.active = false
	return ps.fileWriter.Close()
}

// Subscribe registers a subscriber for a specific topic and returns a channel to receive messages.
func (ps *pubSub[T]) Subscribe(topic string) <-chan T {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan T, 1)
	ps.subscribers[topic] = append(ps.subscribers[topic], ch)
	return ch
}

// SubscribeWithHandler registers a subscriber with a handler function for a specific topic.
// The handler is called whenever a message is published on the topic.
// Returns a function to unsubscribe.
func (ps *pubSub[T]) SubscribeWithHandler(topic string, handler func(msg T)) (Unsubscribe func()) {
	ps.mu.Lock()

	ch := make(chan T, 1)
	chQuit := make(chan struct{})

	ps.subscribers[topic] = append(ps.subscribers[topic], ch)
	ps.mu.Unlock()

	go func() {
		for {
			select {
			case msg := <-ch:
				handler(msg)
			case <-chQuit:
				return
			}
		}
	}()

	return func() {
		close(chQuit)
		ps.Unsubscribe(topic, ch)
	}

}

// Unsubscribe removes a subscriber from a specific topic.
func (ps *pubSub[T]) Unsubscribe(topic string, ch <-chan T) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	subscribers, ok := ps.subscribers[topic]
	if !ok {
		return
	}

	for i, subscriber := range subscribers {
		if subscriber == ch {
			ps.subscribers[topic] = append(subscribers[:i], subscribers[i+1:]...)
			close(subscriber)
			break
		}
	}
}

// Publish sends a message to all subscribers of a specific topic.
// If file writing is active, writes the message to the file.
// Will panic if writing to the file fails.
func (ps *pubSub[T]) Publish(topic string, msg T) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if ps.active {
		if _, err := ps.fileWriter.Write([]byte(fmt.Sprintf("Topic: %s, Time: %s, Message: %v\n", topic, time.Now().Format("2006/01/02 15:04:5"), msg))); err != nil {
			panic(err)
		}
	}

	subscribers, ok := ps.subscribers[topic]
	if !ok {
		return
	}

	for _, subscriber := range subscribers {
		go func(sub chan T) {
			sub <- msg
		}(subscriber)
	}
}
