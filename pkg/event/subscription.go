package event

import (
	"context"
	"errors"
	"time"

	"github.com/seibert-media/golibs/log"
	"go.uber.org/zap"
)

// Subscribe adds a new channel as receiver for events and unsubscribes on a closed ctx
func (b *Broker) Subscribe(ctx context.Context, out chan<- Event) {
	b.new <- out
	go func() {
		<-ctx.Done()
		select {
		case b.closed <- out:
		case <-time.After(1 * time.Second):
			log.From(ctx).Warn("closing subscription", zap.String("reason", "ctx closed"), zap.Error(errors.New("timeout")))
		}
	}()
}

// Unsubscribe removes the provided channel from the active listeners and tells the broker to clean up
func (b *Broker) Unsubscribe(ctx context.Context, out chan<- Event) {
	b.closed <- out
}
