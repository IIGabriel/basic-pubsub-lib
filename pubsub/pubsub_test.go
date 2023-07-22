package pubsub_test

import (
	"os"
	"testing"
	"time"

	"github.com/IIGabriel/basic-pubsub-lib/pubsub"
)

func TestSubscribe(t *testing.T) {
	ps := pubsub.NewPubSub[string]()
	ch := ps.Subscribe("testTopic")

	ps.Publish("testTopic", "testMessage")
	msg := <-ch

	if msg != "testMessage" {
		t.Errorf("Expected 'testMessage', but got %s", msg)
	}

}

func TestSubscribeWithHandler(t *testing.T) {
	ps := pubsub.NewPubSub[string]()
	var receivedMsg string

	unsub := ps.SubscribeWithHandler("testTopic", func(msg string) {
		receivedMsg = msg
	})

	ps.Publish("testTopic", "testMessageWithHandler")
	time.Sleep(5 * time.Millisecond) // Give some time for the handler to execute

	if receivedMsg != "testMessageWithHandler" {
		t.Errorf("Expected 'testMessageWithHandler', but got %s", receivedMsg)
	}

	unsub()
}

func TestUnsubscribe(t *testing.T) {
	ps := pubsub.NewPubSub[string]()
	ch := ps.Subscribe("testTopic")

	ps.Unsubscribe("testTopic", ch)
	ps.Publish("testTopic", "testMessage")

	if _, open := <-ch; open {
		t.Error("Channel not closed")
	}
}

func TestPublishToFile(t *testing.T) {
	ps := pubsub.NewPubSub[string]()
	testArqName := "test.log"
	ps.ActiveFileWriter(testArqName)
	defer ps.CloseFileWriter()

	ps.Publish("testTopic", "testMessageToFile")
	if _, err := os.Stat(testArqName); os.IsNotExist(err) {
		t.Error("test.log does not exist after publishing")
	}
	os.Remove(testArqName)
}
