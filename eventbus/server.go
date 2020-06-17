package eventbus

import (
	"context"
	"gocloud.dev/pubsub"
	"sync"
)

type Request struct {
	Context context.Context
	Message *pubsub.Message
}

type HandlerFunc func(*Request)

type Consumer struct {
	MaxHandler int
	Consumer   *pubsub.Subscription
	Handler    HandlerFunc
	cancelCtx  context.CancelFunc
}

func (s *Consumer) serve(ctx context.Context) {
	defer func() {
		_ = s.Consumer.Shutdown(ctx)
	}()

	// Loop on received messages. We can use a channel as a semaphore to limit how
	// many goroutines we have active at a time as well as wait on the goroutines
	// to finish before exiting.
	sem := make(chan struct{}, s.MaxHandler)
recvLoop:
	for {
		msg, err := s.Consumer.Receive(ctx)
		if err != nil {
			// Errors from Receive indicate that Receive will no longer succeed.
			s.cancelCtx()
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

			// Do work based on the message
			// When operation is successful, must acknowledge message
			ctxHdl, _ := context.WithCancel(ctx)
			s.Handler(&Request{
				Context: ctxHdl,
				Message: msg,
			})
		}()
	}

	// We're no longer receiving messages. Wait to finish handling any
	// unacknowledged messages by totally acquiring the semaphore.
	for n := 0; n < s.MaxHandler; n++ {
		sem <- struct{}{}
	}
}

type Server struct {
	Consumers   []*Consumer
	rootContext context.Context
	mtx         *sync.Mutex
}

func NewServer(ctx context.Context, cs ...*Consumer) *Server {
	return &Server{
		Consumers:   cs,
		rootContext: ctx,
		mtx:         new(sync.Mutex),
	}
}

func (s *Server) AddConsumer(c *Consumer) {
	s.Consumers = append(s.Consumers, c)
}

func (s *Server) Serve() error {
	for _, c := range s.Consumers {
		ctxSub, cancel := context.WithCancel(s.rootContext)
		c.cancelCtx = cancel
		go c.serve(ctxSub)
	}

	select {
	case <-s.rootContext.Done():
		return nil
	}
}

func (s *Server) Close() {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.rootContext.Done()
}
