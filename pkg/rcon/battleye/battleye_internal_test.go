package battleye

import (
	"context"
	"time"

	be_mocks "github.com/playnet-public/battleye/mocks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/playnet-public/gorcon/pkg/rcon"
)

var _ = Describe("Connection", func() {
	var (
		con   *Connection
		proto *be_mocks.Protocol
		ctx   context.Context
	)

	BeforeEach(func() {
		proto = &be_mocks.Protocol{}
		ctx = context.Background()
		con = NewConnection(ctx)
		con.Protocol = proto
	})

	Describe("Subscribe", func() {
		BeforeEach(func() {
			con = NewConnection(ctx)
			ctx = context.Background()
		})
		It("does add subscription on Listen", func() {
			l := len(con.subscriptions)
			ctx, _ := context.WithCancel(ctx)
			con.Subscribe(ctx, make(chan *rcon.Event))
			time.Sleep(1 * time.Second)
			con.subscriptionsMutex.RLock()
			Expect(len(con.subscriptions)).To(BeEquivalentTo(l + 1))
			con.subscriptionsMutex.RUnlock()
		})
		It("does end subscription on closed context", func() {
			ctx, close := context.WithCancel(ctx)
			l := len(con.subscriptions)
			con.Subscribe(ctx, make(chan *rcon.Event))
			con.subscriptionsMutex.RLock()
			Expect(len(con.subscriptions)).To(BeEquivalentTo(l + 1))
			con.subscriptionsMutex.RUnlock()
			close()
			time.Sleep(1 * time.Second)
			con.subscriptionsMutex.RLock()
			Expect(len(con.subscriptions)).To(BeEquivalentTo(l))
			con.subscriptionsMutex.RUnlock()
		})
		It("does remove the correct subscription on closed context", func(done Done) {
			ctx, _ := context.WithCancel(ctx)
			c1 := make(chan *rcon.Event)
			c2 := make(chan *rcon.Event)
			ctx2, close2 := context.WithCancel(ctx)
			c3 := make(chan *rcon.Event)
			con.Subscribe(ctx, c1)
			con.Subscribe(ctx2, c2)
			con.Subscribe(ctx, c3)

			go func() {
				Expect(<-c1).NotTo(BeNil())
			}()
			go func() {
				Expect(<-c2).NotTo(BeNil())
			}()
			go func() {
				Expect(<-c3).NotTo(BeNil())
			}()
			e := &rcon.Event{}
			con.subscriptionsMutex.RLock()
			defer con.subscriptionsMutex.RUnlock()
			for _, l := range con.subscriptions {
				go func(l chan *rcon.Event) { l <- e }(l)
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
				case <-time.After(300 * time.Millisecond):
					Expect(true).To(BeTrue())
				}
				close(done)
			}()
			go func() {
				Expect(<-c3).NotTo(BeNil())
			}()
			con.subscriptionsMutex.RLock()
			defer con.subscriptionsMutex.RUnlock()
			for _, l := range con.subscriptions {
				go func(l chan *rcon.Event) { l <- e }(l)
			}
		})
	})
})
