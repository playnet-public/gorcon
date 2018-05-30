package watcher

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	context "github.com/seibert-media/golibs/log"
)

var _ = Describe("Watcher", func() {
	var (
		w   *Watcher
		ctx context.Context
	)

	BeforeEach(func() {
		ctx = context.NewNop()
		w = &Watcher{}
	})

	Describe("Subscribe", func() {
		BeforeEach(func() {
			w = &Watcher{}
			ctx = context.NewNop()
		})
		It("does add subscription on Listen", func() {
			l := len(w.subscriptions)
			ctx, _ := context.WithCancel(ctx)
			w.Subscribe(ctx, make(chan *Event))
			time.Sleep(1 * time.Second)
			w.subscriptionsMutex.RLock()
			Expect(len(w.subscriptions)).To(BeEquivalentTo(l + 1))
			w.subscriptionsMutex.RUnlock()
		})
		It("does end subscription on closed context", func() {
			ctx, close := context.WithCancel(ctx)
			l := len(w.subscriptions)
			w.Subscribe(ctx, make(chan *Event))
			w.subscriptionsMutex.RLock()
			Expect(len(w.subscriptions)).To(BeEquivalentTo(l + 1))
			w.subscriptionsMutex.RUnlock()
			close()
			time.Sleep(1 * time.Second)
			w.subscriptionsMutex.RLock()
			Expect(len(w.subscriptions)).To(BeEquivalentTo(l))
			w.subscriptionsMutex.RUnlock()
		})
		It("does remove the correct subscription on closed context", func(done Done) {
			ctx, _ := context.WithCancel(ctx)
			c1 := make(chan *Event)
			c2 := make(chan *Event)
			ctx2, close2 := context.WithCancel(ctx)
			c3 := make(chan *Event)
			w.Subscribe(ctx, c1)
			w.Subscribe(ctx2, c2)
			w.Subscribe(ctx, c3)

			go func() {
				Expect(<-c1).NotTo(BeNil())
			}()
			go func() {
				Expect(<-c2).NotTo(BeNil())
			}()
			go func() {
				Expect(<-c3).NotTo(BeNil())
			}()
			e := &Event{}
			w.subscriptionsMutex.RLock()
			defer w.subscriptionsMutex.RUnlock()
			for _, l := range w.subscriptions {
				go func(l chan *Event) { l <- e }(l)
			}
			close2()
			go func() {
				Expect(<-c1).NotTo(BeNil())
			}()
			go func() {
				defer GinkgoRecover()
				select {
				case <-c2:
					Fail("should not receive on c2")
				case <-time.After(500 * time.Millisecond):
					Expect(true).To(BeTrue())
				}
				close(done)
			}()
			go func() {
				Expect(<-c3).NotTo(BeNil())
			}()
			time.Sleep(500 * time.Microsecond)
			w.subscriptionsMutex.RLock()
			defer w.subscriptionsMutex.RUnlock()
			for _, l := range w.subscriptions {
				go func(l chan *Event) { l <- e }(l)
			}
		})
	})
})
