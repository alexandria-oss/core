package eventbus

import (
	"github.com/google/uuid"
	"strings"
	"sync"
	"time"
)

// Event Represents an event log record for metadata
//	It is formed by the following fields:
//	- Service Name = Service who dispatched the event, aka. Event source
//	- Transaction ID = Distributed transaction ID *Only for SAGA transaction
//	- Event Type = Type of the event dispatched (integration or domain)
//	- Content = Message body, mostly bytes or marshalled JSON
//	- Priority = Event's priority type
//	- Provider = Message Broker/Queue-Notification Provider (Kafka, RabbitMQ, AWS)
//	- Dispatch Time = Event's dispatching timestamp
type Event struct {
	ID string `json:"event_id"`
	// ServiceName Service who dispatched the event, aka. Event source
	ServiceName string `json:"service_name"`
	// Event Type Type of the event dispatched (integration or domain)
	EventType string `json:"event_type"`
	// Content Message body, mostly bytes or marshalled JSON
	Content []byte `json:"content"`
	// Priority Event's priority type
	Priority string `json:"importance"`
	// Provider Message Broker/Queue-Notification Provider (Kafka, RabbitMQ, AWS)
	Provider string `json:"provider"`
	// DispatchTime Event's dispatching timestamp
	DispatchTime time.Time `json:"dispatch_time"`
}

// Transaction represents a SAGA-like transaction entity
type Transaction struct {
	// ID Transaction ID
	ID string `json:"transaction_id"`
	// RootID Aggregate/Entity's ID
	RootID string `json:"root_id"`
	// SpanID OpenTracing/OpenCensus root span ID
	SpanID string `json:"span_id,omitempty"`
	// TraceID OpenTracing/OpenCensus trace ID
	TraceID string `json:"trace_id,omitempty"`
	// Operation Kind of operation to perform
	Operation string `json:"operation"`
	// Backup Aggregate/Entity's backup for update-like operations
	Backup string `json:"backup,omitempty"`
	// Code HTTP-like status code
	Code string `json:"code"`
	// Message Custom message for logging
	Message string `json:"message,omitempty"`
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
func NewEvent(serviceName, eventType, priority, provider string, content []byte) *Event {
	mtx.Lock()
	defer mtx.Unlock()

	// Validate payload
	eventType = strings.ToUpper(eventType)
	eventType = isEventTypeValid(eventType)

	priority = strings.ToUpper(priority)
	priority = isPriorityValid(priority)

	provider = strings.ToUpper(provider)
	provider = isProviderValid(provider)

	return &Event{
		ID:           uuid.New().String(),
		ServiceName:  strings.ToUpper(serviceName),
		EventType:    eventType,
		Content:      content,
		Priority:     priority,
		Provider:     provider,
		DispatchTime: time.Now(),
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
