package battleye

import (
	"context"

	be_mocks "github.com/playnet-public/battleye/mocks"
	"github.com/playnet-public/gorcon/pkg/event"

	. "github.com/onsi/ginkgo"
	//. "github.com/onsi/gomega"
)

var _ = Describe("Connection", func() {
	var (
		ctx   context.Context
		c     chan event.Event
		b     *event.Broker
		con   *Connection
		proto *be_mocks.Protocol
	)

	BeforeEach(func() {
		ctx = context.Background()
		c = make(chan event.Event)
		b = event.NewBroker(ctx, c)
		con = NewConnection(ctx, b, c)
		proto = &be_mocks.Protocol{}
		con.Protocol = proto
	})

	Describe("Subscribe", func() {
		BeforeEach(func() {
			con = NewConnection(ctx, b, c)
			ctx = context.Background()
		})
	})
})
