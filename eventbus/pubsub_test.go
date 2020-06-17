package eventbus

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEvent(t *testing.T) {
	event := NewEvent("author", EventIntegration, PriorityMid, ProviderRabbitMQ, []byte("message 1"))
	t.Log(event)
	event2 := NewEvent("media", EventDomain, PriorityHigh, ProviderKafka, []byte("message 2"))
	t.Log(event2)
	assert.NotEqual(t, event.TransactionID, event2.TransactionID, "Distributed ID are not unique")
	assert.NotEqual(t, event.ID, event2.ID, "Event ID are not unique")

	event = NewEvent("author", EventIntegration, PriorityMid, ProviderRabbitMQ, []byte("message 1"))
	t.Log(event)
	assert.NotEqual(t, event.TransactionID, "")
}

func BenchmarkNewEvent(b *testing.B) {
	for n := 0; n < b.N; n++ {
		NewEvent("author", EventIntegration, PriorityMid, ProviderRabbitMQ, []byte("message 1"))
	}
}
