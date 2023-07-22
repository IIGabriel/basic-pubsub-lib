# Basic PubSub Library
A lightweight and efficient publish-subscribe system in Go, enhanced with optional file logging features.

## Key Features:

- Easy-to-use API for publishing and subscribing to topics.
- Built with concurrency in mind to ensure safety across multiple operations.
- Option to log messages to a file for auditing or debugging purposes.

## Usage Guide:

### Initialization

Initialize the PubSub system with:
```go
ps := pubsub.NewPubSub[string]()
```
It is possible to notice that generics are used, you can choose the type of message
### Activating File Writer

To log messages to a file, activate the file writer. Note that this will panic if there's an issue opening the file:
```go
ps.ActiveFileWriter("path_to_file.txt")
```
### Deactivating and Closing File Writer

To stop logging and close the file:
```go
ps.CloseFileWriter()
```
### Subscribing to Topics

Register for a specific topic and get a channel to receive its messages:
```go
ch := ps.Subscribe("your_topic")
```
### Subscribing with a Handler

Register a subscriber that uses a handler function to process messages for a specific topic:
```go
unsubscribe := ps.SubscribeWithHandler("your_topic", yourHandlerFunction)
```

### Unsubscribing from Topics

To stop receiving messages from a topic:
```go
ps.Unsubscribe("your_topic", ch)
```
### Publishing Messages

Send a message to all subscribers of a specific topic. If the file writer is active, the message will also be logged to
the file. A panic will occur if there's an error writing to the file:
```go
ps.Publish("your_topic", "your_message")
```