package rcon_test

import (
	"errors"
	"testing"

	"github.com/playnet-public/gorcon/pkg/mocks"
	"github.com/playnet-public/gorcon/pkg/rcon"
	context "github.com/seibert-media/golibs/log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGoRcon(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rcon Suite")
}

var _ = Describe("Rcon", func() {
	var (
		r              *rcon.Rcon
		mockClient     *mocks.RconClient
		mockConnection *mocks.RconConnection
		ctx            context.Context
	)

	BeforeEach(func() {
		mockClient = &mocks.RconClient{}
		r = &rcon.Rcon{
			Client: mockClient,
		}
		mockConnection = &mocks.RconConnection{}
		mockClient.NewConnectionReturns(mockConnection)
		ctx = context.NewNop()
	})

	Describe("Connect", func() {
		It("returns no error", func() {
			Expect(r.Connect(ctx)).To(BeNil())
		})
		It("returns error if Client is nil", func() {
			r.Client = nil
			Expect(r.Connect(ctx)).NotTo(BeNil())
		})
		It("calls Client.NewConnection once", func() {
			r.Connect(ctx)
			Expect(mockClient.NewConnectionCallCount()).To(BeEquivalentTo(1))
		})
		It("returns error if Client returns nil connection", func() {
			mockClient.NewConnectionReturns(nil)
			Expect(r.Connect(ctx)).NotTo(BeNil())
		})
		It("returns error if current connection is not nil", func() {
			r.Con = &mocks.RconConnection{}
			Expect(r.Connect(ctx)).NotTo(BeNil())
		})
		It("sets the connection", func() {
			Expect(r.Con).To(BeNil())
			Expect(r.Connect(ctx)).To(BeNil())
			Expect(r.Con).NotTo(BeNil())
		})
		It("opens the new connection", func() {
			r.Connect(ctx)
			Expect(mockConnection.OpenCallCount()).To(BeEquivalentTo(1))
		})
		It("returns error if opening the connection fails", func() {
			mockConnection.OpenReturns(errors.New("test"))
			Expect(r.Connect(ctx)).NotTo(BeNil())
		})
	})

	Describe("Write", func() {
		BeforeEach(func() {
			r.Con = mockConnection
		})
		It("does return error on nil connection", func() {
			r.Con = nil
			_, err := r.Write(ctx, "")
			Expect(err).NotTo(BeNil())
		})
		It("does call Write on connection", func() {
			r.Write(ctx, "")
			Expect(mockConnection.WriteCallCount()).To(BeEquivalentTo(1))
		})
	})

	Describe("Reconnect", func() {
		BeforeEach(func() {
			r.Con = mockConnection
		})
		It("returns no error", func() {
			Expect(r.Reconnect(ctx)).To(BeNil())
		})
		It("returns error if Client is nil", func() {
			r.Client = nil
			Expect(r.Reconnect(ctx)).NotTo(BeNil())
		})
		It("calls Con.Close once", func() {
			r.Reconnect(ctx)
			Expect(mockConnection.CloseCallCount()).To(BeEquivalentTo(1))
		})
		It("calls Con.Close once before the connection gets overwritten", func() {
			tempCon := &mocks.RconConnection{}
			r.Con = tempCon
			r.Reconnect(ctx)
			Expect(tempCon.CloseCallCount()).To(BeEquivalentTo(1))
		})
		It("calls Client.NewConnection once", func() {
			r.Reconnect(ctx)
			Expect(mockClient.NewConnectionCallCount()).To(BeEquivalentTo(1))
		})
		It("returns error if Client returns nil connection", func() {
			mockClient.NewConnectionReturns(nil)
			Expect(r.Reconnect(ctx)).NotTo(BeNil())
		})
	})

	Describe("Disconnect", func() {
		BeforeEach(func() {
			r.Con = mockConnection
		})
		It("returns no error", func() {
			Expect(r.Disconnect(ctx)).To(BeNil())
		})
		It("returns error if current connection is nil", func() {
			r.Con = nil
			Expect(r.Disconnect(ctx)).NotTo(BeNil())
		})
		It("calls Con.Close once", func() {
			r.Disconnect(ctx)
			Expect(mockConnection.CloseCallCount()).To(BeEquivalentTo(1))
		})
		It("returns error if Con.Close fails", func() {
			mockConnection.CloseReturns(errors.New("test"))
			Expect(r.Disconnect(ctx)).NotTo(BeNil())
		})
		It("sets the connection to nil after disconnect", func() {
			r.Disconnect(ctx)
			Expect(r.Con).To(BeNil())
		})
	})
})
