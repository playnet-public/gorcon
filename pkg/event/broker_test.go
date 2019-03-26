package event_test

import (
	"context"
	"testing"
	"time"

	"github.com/playnet-public/gorcon/pkg/event"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/seibert-media/golibs/log"
)

const debug = false

func TestEvent(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Event Suite")
}

type fakeEvent struct{}

func (f *fakeEvent) Timestamp() time.Time {
	t, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	return t
}
func (f *fakeEvent) Kind() string { return "fake" }
func (f *fakeEvent) Data() string { return "fake" }

var _ = Describe("Event", func() {

	setup := func() (ctx context.Context, in chan event.Event, b *event.Broker) {
		ctx = context.Background()
		l := log.New("", debug)
		ctx = log.WithLogger(ctx, l)
		in = make(chan event.Event)
		b = event.NewBroker(ctx, in)
		return
	}

	Describe("Run", func() {
		It("does exit on closed context", func() {
			ctx, _, b := setup()

			ctx, close := context.WithCancel(ctx)
			close()
			Expect(b.Run(ctx)).To(BeEquivalentTo(context.Canceled))
		})
		It("does return error on closed input", func() {
			ctx, in, b := setup()

			close(in)
			Expect(b.Run(ctx)).To(BeEquivalentTo(event.ErrInputClosed))
		})
		It("does cleanup and close all subscriptions on exit", func() {
			ctx, in, b := setup()

			go b.Run(ctx)

			c1 := make(chan event.Event)
			b.Subscribe(ctx, c1)

			c2 := make(chan event.Event)
			b.Subscribe(ctx, c2)

			close(in)

			e1, ok1 := <-c1
			e2, ok2 := <-c2
			Expect(e1).To(BeNil())
			Expect(e2).To(BeNil())
			Expect(ok1).To(BeFalse())
			Expect(ok2).To(BeFalse())
		})
		It("does not block on inactive subscriptions", func() {
			ctx, in, b := setup()

			go b.Run(ctx)

			c1 := make(chan event.Event)
			b.Subscribe(ctx, c1)

			c2 := make(chan event.Event)
			b.Subscribe(ctx, c2)

			go func() { in <- &fakeEvent{} }()

			Expect(<-c1).NotTo(BeNil())

			go func() { in <- &fakeEvent{} }()

			Expect(<-c1).NotTo(BeNil())
		})
	})

	Describe("Subscribe", func() {
		It("does forward events to new subscriptions", func() {
			ctx, in, b := setup()

			go b.Run(ctx)
			c := make(chan event.Event)
			b.Subscribe(ctx, c)
			go func() { in <- &fakeEvent{} }()
			Expect(<-c).NotTo(BeNil())
		})
		It("does unsubscribe on closed context", func() {
			ctx, in, b := setup()

			go b.Run(ctx)
			ctx, close := context.WithCancel(ctx)
			c := make(chan event.Event)
			b.Subscribe(ctx, c)
			go func() { in <- &fakeEvent{} }()
			Expect(<-c).NotTo(BeNil())

			close()
			// wait until we check, as it might take one event
			<-time.After(1 * time.Millisecond)

			go func() { in <- &fakeEvent{} }()

			e, ok := <-c
			Expect(ok).To(BeFalse())
			Expect(e).To(BeNil())
		})
	})

	Describe("Unsubscribe", func() {
		It("does not forward events to closed subscriptions", func() {
			ctx, in, b := setup()

			go b.Run(ctx)
			c := make(chan event.Event)
			b.Subscribe(ctx, c)
			b.Unsubscribe(ctx, c)
			go func() { in <- &fakeEvent{} }()
			select {
			case e := <-c:
				Expect(e).To(BeNil())
			case <-time.After(5 * time.Millisecond):
			}
		})
	})
})
