package event

import (
	"context"
	"errors"

	"github.com/seibert-media/golibs/log"
	"go.uber.org/zap"
)

var (
	// ErrInputClosed is returned when the input channel for the broker gets closed
	ErrInputClosed = errors.New("input channel closed")
)

// Broker for subscribing to an eventsource with multiple subscriptions automatically canceled on ctx.Close
type Broker struct {
	new    chan chan<- Event
	active map[chan<- Event]struct{}
	closed chan chan<- Event

	in <-chan Event
}

// NewBroker with the provided input channel as event source
// A running broker will handle all incoming events by sending them to all active subscriptions
func NewBroker(ctx context.Context, in <-chan Event) *Broker {
	return &Broker{
		new:    make(chan chan<- Event),
		active: make(map[chan<- Event]struct{}),
		closed: make(chan chan<- Event),

		in: in,
	}
}

// Run the broker and listen for new subscriptions, events and unsubscribes
// The broker will run until either it's parent context closes or the incoming event channel gets closed
func (b *Broker) Run(ctx context.Context) error {
	defer func() {
		for s := range b.active {
			delete(b.active, s)
			close(s)
		}
	}()
	for {
		select {
		case <-ctx.Done():
			log.From(ctx).Info("stopping background loop", zap.Error(ctx.Err()))
			return ctx.Err()

		case s := <-b.new:
			b.active[s] = struct{}{}
			log.From(ctx).Debug("subscribing", zap.Int("count", len(b.active)))

		case s := <-b.closed:
			delete(b.active, s)
			close(s)
			log.From(ctx).Debug("unsubscribing", zap.Int("count", len(b.active)))

		case event, ok := <-b.in:
			if !ok {
				log.From(ctx).Info("stopping background loop")
				return ErrInputClosed
			}

			for s := range b.active {
				log.From(ctx).Debug("handling event", zap.String("data", event.Data()))
				// TODO(kwiesmueller): make sure we don't leak here and need an unsubscribing timeout
				go func(s chan<- Event) { s <- event }(s)
			}
		}
	}
}
