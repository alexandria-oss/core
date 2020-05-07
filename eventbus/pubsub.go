package eventbus

import (
	"context"
	"github.com/google/uuid"
	"strings"
	"time"

	"gocloud.dev/pubsub"
)

// Event Represents an event log record
/*
*	Service Name = Service who dispatched the event
*	Transaction ID = Distributed transaction ID *Only for SAGA transaction
*	Event Type = Type of the event dispatched (integration or domain)
*	Content = Message body, mostly bytes or JSON-into string
*	Importance = Event's importance
*	Provider = Message Broker Provider (Kafka, RabbitMQ)
*	Dispatch Time = Event's dispatching timestamp
 */
type Event struct {
	ID            string `json:"event_id"`
	ServiceName   string `json:"service_name"`
	TransactionID string `json:"transaction_id,omitempty"`
	EventType     string `json:"event_type"`
	Content       string `json:"content"`
	Importance    string `json:"importance"`
	Provider      string `json:"provider"`
	DispatchTime  int64  `json:"dispatch_time"`
}

func NewEvent(serviceName, eventType, content, importance, provider string) *Event {
	return &Event{
		ID:            uuid.New().String(),
		ServiceName:   strings.ToUpper(serviceName),
		TransactionID: uuid.New().String(),
		EventType:     strings.ToUpper(eventType),
		Content:       content,
		Importance:    strings.ToUpper(importance),
		Provider:      strings.ToUpper(provider),
		// Unix to millis
		DispatchTime: time.Now().UnixNano() / 1000000,
	}
}

// ListenSubscriber Listen to a Pub/Sub subscription concurrently
func ListenSubscriber(ctx context.Context, subscription *pubsub.Subscription) {
	defer subscription.Shutdown(ctx)

	// Loop on received messages. We can use a channel as a semaphore to limit how
	// many goroutines we have active at a time as well as wait on the goroutines
	// to finish before exiting.
	const maxHandlers = 10
	sem := make(chan struct{}, maxHandlers)
recvLoop:
	for {
		msg, err := subscription.Receive(ctx)
		if err != nil {
			// Errors from Receive indicate that Receive will no longer succeed.
			// logger.Log("msg", err.Error())
			break
		}

		// Wait if there are too many active handle goroutines and acquire the
		// semaphore. If the context is canceled, stop waiting and start shutting
		// down.
		select {
		case sem <- struct{}{}:
		case <-ctx.Done():
			break recvLoop
		}

		// Handle the message in a new goroutine.
		go func() {
			defer func() { <-sem }() // Release the semaphore.
			defer msg.Ack()          // Messages must always be acknowledged with Ack.

			// Do work based on the message, for example:
			// logger.Log("msg", string(msg.Body))
			// logger.Log("msg", fmt.Sprintf("%v", msg.Metadata))
		}()
	}

	// We're no longer receiving messages. Wait to finish handling any
	// unacknowledged messages by totally acquiring the semaphore.
	for n := 0; n < maxHandlers; n++ {
		sem <- struct{}{}
	}
}
