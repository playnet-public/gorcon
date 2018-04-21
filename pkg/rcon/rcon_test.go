package rcon_test

import (
	"errors"
	"testing"

	"github.com/playnet-public/gorcon/pkg/mocks"
	"github.com/playnet-public/gorcon/pkg/rcon"

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
	)

	BeforeEach(func() {
		mockClient = &mocks.RconClient{}
		r = &rcon.Rcon{
			Client: mockClient,
		}
		mockConnection = &mocks.RconConnection{}
		mockClient.NewConnectionReturns(mockConnection)
	})

	Describe("Connect", func() {
		It("returns no error", func() {
			Expect(r.Connect()).To(BeNil())
		})
		It("returns error if Client is nil", func() {
			r.Client = nil
			Expect(r.Connect()).NotTo(BeNil())
		})
		It("calls Client.NewConnection once", func() {
			r.Connect()
			Expect(mockClient.NewConnectionCallCount()).To(BeEquivalentTo(1))
		})
		It("returns error if Client returns nil connection", func() {
			mockClient.NewConnectionReturns(nil)
			Expect(r.Connect()).NotTo(BeNil())
		})
		It("returns error if current connection is not nil", func() {
			r.Con = &mocks.RconConnection{}
			Expect(r.Connect()).NotTo(BeNil())
		})
		It("sets the connection", func() {
			Expect(r.Con).To(BeNil())
			Expect(r.Connect()).To(BeNil())
			Expect(r.Con).NotTo(BeNil())
		})
		It("opens the new connection", func() {
			r.Connect()
			Expect(mockConnection.OpenCallCount()).To(BeEquivalentTo(1))
		})
		It("returns error if opening the connection fails", func() {
			mockConnection.OpenReturns(errors.New("test"))
			Expect(r.Connect()).NotTo(BeNil())
		})
	})

	Describe("Reconnect", func() {
		BeforeEach(func() {
			r.Con = mockConnection
		})
		It("returns no error", func() {
			Expect(r.Reconnect()).To(BeNil())
		})
		It("returns error if Client is nil", func() {
			r.Client = nil
			Expect(r.Reconnect()).NotTo(BeNil())
		})
		It("calls Con.Close once", func() {
			r.Reconnect()
			Expect(mockConnection.CloseCallCount()).To(BeEquivalentTo(1))
		})
		It("calls Con.Close once before the connection gets overwritten", func() {
			tempCon := &mocks.RconConnection{}
			r.Con = tempCon
			r.Reconnect()
			Expect(tempCon.CloseCallCount()).To(BeEquivalentTo(1))
		})
		It("calls Client.NewConnection once", func() {
			r.Reconnect()
			Expect(mockClient.NewConnectionCallCount()).To(BeEquivalentTo(1))
		})
		It("returns error if Client returns nil connection", func() {
			mockClient.NewConnectionReturns(nil)
			Expect(r.Reconnect()).NotTo(BeNil())
		})
	})

	Describe("Disconnect", func() {
		BeforeEach(func() {
			r.Con = mockConnection
		})
		It("returns no error", func() {
			Expect(r.Disconnect()).To(BeNil())
		})
		It("returns error if current connection is nil", func() {
			r.Con = nil
			Expect(r.Disconnect()).NotTo(BeNil())
		})
		It("calls Con.Close once", func() {
			r.Disconnect()
			Expect(mockConnection.CloseCallCount()).To(BeEquivalentTo(1))
		})
		It("returns error if Con.Close fails", func() {
			mockConnection.CloseReturns(errors.New("test"))
			Expect(r.Disconnect()).NotTo(BeNil())
		})
		It("sets the connection to nil after disconnect", func() {
			r.Disconnect()
			Expect(r.Con).To(BeNil())
		})
	})
})
