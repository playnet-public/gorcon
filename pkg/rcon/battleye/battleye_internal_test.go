package battleye

import (
	"testing"
	"time"

	be_mocks "github.com/playnet-public/battleye/mocks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/playnet-public/gorcon/pkg/rcon"
	context "github.com/seibert-media/golibs/log"
)

func TestBattlEye(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BattlEye Internal Suite")
}

var _ = Describe("Connection", func() {
	var (
		con   *Connection
		proto *be_mocks.Protocol
		ctx   context.Context
	)

	BeforeEach(func() {
		proto = &be_mocks.Protocol{}
		ctx = context.NewNop()
		con = NewConnection(ctx)
		con.Protocol = proto
	})

	Describe("Subscribe", func() {
		BeforeEach(func() {
			con = NewConnection(ctx)
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
		It("does remove the correct subscription on closed context", func() {
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
			c1 <- e
			c2 <- e
			c3 <- e
			close2()
			go func() {
				Expect(<-c1).NotTo(BeNil())
			}()
			go func() {
				select {
				case ev := <-c2:
					Expect(ev).To(BeNil())
				case <-time.After(100 * time.Millisecond):
					Expect(true).To(BeTrue())
				}
			}()
			go func() {
				Expect(<-c3).NotTo(BeNil())
			}()
			c1 <- e
			go func() { c2 <- e }()
			c3 <- e
		})
	})
})
