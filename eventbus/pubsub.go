package eventbus

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sony/sonyflake"
	"gocloud.dev/pubsub"
)

// Event Represents an event log record
//	It is formed by the following fields:
//	- Service Name = Service who dispatched the event, aka. Event source
//	- Transaction ID = Distributed transaction ID *Only for SAGA transaction
//	- Event Type = Type of the event dispatched (integration or domain)
//	- Content = Message body, mostly bytes or marshalled JSON
//	- Priority = Event's priority type
//	- Provider = Message Broker/Queue-Notification Provider (Kafka, RabbitMQ)
//	- Dispatch Time = Event's dispatching timestamp
type Event struct {
	ID            string `json:"event_id"`
	ServiceName   string `json:"service_name"`
	TransactionID uint64 `json:"transaction_id,omitempty"`
	EventType     string `json:"event_type"`
	Content       []byte `json:"content"`
	Priority      string `json:"importance"`
	Provider      string `json:"provider"`
	DispatchTime  int64  `json:"dispatch_time"`
}

const (
	// EventDomain Domain event type
	EventDomain = "EVENT_DOMAIN"
	// EventIntegration Integration event type
	EventIntegration = "EVENT_INTEGRATION"
	// ProviderKafka Apache Kafka provider type
	ProviderKafka = "PROVIDER_KAFKA"
	// ProviderRabbitMQ RabbitMQ provider type
	ProviderRabbitMQ = "PROVIDER_RABBITMQ"
	// ProviderNATS NATS provider type
	ProviderNATS = "PROVIDER_NATS"
	// ProviderAWS AWS SQS/SNS provider type
	ProviderAWS = "PROVIDER_AWS"
	// PriorityLow Low event's priority
	PriorityLow = "PRIORITY_LOW"
	// PriorityMid Middle event's priority
	PriorityMid = "PRIORITY_MID"
	// PriorityHigh High event's priority
	PriorityHigh = "PRIORITY_HIGH"
)

var mtx *sync.Mutex

func init() {
	mtx = new(sync.Mutex)
}

// NewEvent returns a new Event ready to be used by an Event Bus
func NewEvent(serviceName, eventType, priority, provider string, content []byte, isTransaction bool) *Event {
	mtx.Lock()
	defer mtx.Unlock()

	// Generate distributed ID
	var err error
	var distID uint64
	distID = 0

	if isTransaction {
		flake := sonyflake.NewSonyflake(sonyflake.Settings{
			StartTime:      time.Time{},
			MachineID:      nil,
			CheckMachineID: nil,
		})

		distID, err = flake.NextID()
		if err != nil {
			return nil
		}
	}

	// Validate payload
	eventType = strings.ToUpper(eventType)
	eventType = isEventTypeValid(eventType)

	priority = strings.ToUpper(priority)
	priority = isPriorityValid(priority)

	provider = strings.ToUpper(provider)
	provider = isProviderValid(provider)

	return &Event{
		ID:            uuid.New().String(),
		ServiceName:   strings.ToUpper(serviceName),
		TransactionID: distID,
		EventType:     eventType,
		Content:       content,
		Priority:      priority,
		Provider:      provider,
		// Unix to millis
		DispatchTime: time.Now().UnixNano() / 1000000,
	}
}

func isEventTypeValid(eventType string) string {
	if eventType != EventDomain && eventType != EventIntegration {
		return EventDomain
	}

	return eventType
}

func isPriorityValid(priority string) string {
	if priority != PriorityLow && priority != PriorityMid && priority != PriorityHigh {
		return PriorityLow
	}

	return priority
}

func isProviderValid(provider string) string {
	if provider != ProviderKafka && provider != ProviderRabbitMQ && provider != ProviderAWS && provider != ProviderNATS {
		return ProviderKafka
	}

	return provider
}

// ListenSubscriber Listen to a subscription concurrently
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
