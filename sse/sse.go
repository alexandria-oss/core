package sse

import (
	"encoding/json"
	"fmt"
	"github.com/alexandria-oss/core"
	"github.com/go-kit/kit/log"
	"net/http"
	"sync"
)

// Event entity, required for SSE operations
type Event struct {
	ID      uint64
	Message []byte
	// Consumer client's distributed ID
	Consumer uint64
}

// NewEvent returns an event entity with a valid event ID
func NewEvent(msg []byte, c uint64) Event {
	// Generate unique id for each event
	id := core.NewSonyflakeID()

	return Event{
		ID:       id,
		Message:  msg,
		Consumer: c,
	}
}

// String parses an event into a single string
func (e Event) String() string {
	return fmt.Sprintf("%d: [%s,%d)", e.ID, string(e.Message), e.Consumer)
}

// Broker manages all SSE event transactions and contains a consumer pool
type Broker struct {
	// consumers subscriber pool using which assigns a Distributed ID for each client
	consumers map[chan Event]uint64
	logger    log.Logger
	mtx       *sync.Mutex
}

// NewBroker returns a valid SSE broker
func NewBroker(logger log.Logger) *Broker {
	return &Broker{
		consumers: make(map[chan Event]uint64),
		mtx:       new(sync.Mutex),
		logger:    logger,
	}
}

// Subscribe returns a new broker consumer; listens to broker's events and generates a
// valid consumer ID
func (b *Broker) Subscribe() chan Event {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	// Generate unique id for each client
	id := core.NewSonyflakeID()

	c := make(chan Event)
	b.consumers[c] = id

	_ = b.logger.Log("resource", "transport.eventsource.broker", "msg",
		fmt.Sprintf("client %d connected", id))
	return c
}

// Unsubscribe removes a consumer from broker's pool
func (b *Broker) Unsubscribe(c chan Event) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	id := b.consumers[c]
	close(c)
	delete(b.consumers, c)
	_ = b.logger.Log("resource", "transport.eventsource.broker", "msg",
		fmt.Sprintf("client %d killed, %d remaining", id, len(b.consumers)))
}

// Publish issues a new event to either one or many consumers
func (b *Broker) Publish(e Event) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	pubMsg := 0
	for s, id := range b.consumers {
		if e.Consumer > 0 {
			// Push to specific consumer
			if id == e.Consumer {
				s <- e
				pubMsg++
				break
			}
		} else {
			// Push to every consumer
			e.Consumer = id
			s <- e
			// Reset unused consumer
			e.Consumer = 0
			pubMsg++
		}
	}

	_ = b.logger.Log("resource", "transport.eventsource.broker", "msg",
		fmt.Sprintf("published message to %d subscribers", pubMsg))
}

// Close removes any channels leftovers from broker's pool
func (b *Broker) Close() {
	for k, _ := range b.consumers {
		close(k)
		delete(b.consumers, k)
	}
}

// ServeHTTP receives and attach a new broker subscription to every HTTP request
// using required streaming HTTP headers
func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming is not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")

	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create new client channel for stream events
	c := b.Subscribe()
	defer b.Unsubscribe(c)

	// Send client it's new ID
	go b.Publish(NewEvent(nil, b.consumers[c]))

	for {
		select {
		case msg := <-c:
			// MIME Type (text/event-stream) formatted, DO NOT MODIFY IT
			msgJSON, err := json.Marshal(struct {
				ID       uint64 `json:"event_id"`
				Message  string `json:"message"`
				Consumer uint64 `json:"consumer_id"`
			}{msg.ID, string(msg.Message), msg.Consumer})
			if err != nil {
				_, _ = fmt.Fprintf(w, "data: %s\n\n", msg)
			} else {
				_, _ = fmt.Fprintf(w, "data: %s\n\n", msgJSON)
			}
			f.Flush()
		case <-ctx.Done():
			return
		}
	}
}
