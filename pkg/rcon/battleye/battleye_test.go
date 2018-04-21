package battleye_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/playnet-public/gorcon/pkg/rcon"
	be "github.com/playnet-public/gorcon/pkg/rcon/battleye"
)

func TestBattlEye(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BattlEye Suite")
}

var _ = Describe("Client", func() {
	var (
		c *be.Client
	)

	BeforeEach(func() {
		c = &be.Client{}
	})

	Describe("NewConnection", func() {
		It("does not return nil", func() {
			Expect(c.NewConnection()).NotTo(BeNil())
		})
	})
})

var _ = Describe("Connection", func() {
	var (
		con *be.Connection
	)

	BeforeEach(func() {
		con = &be.Connection{}
	})

	Describe("Open", func() {
		It("does not return error", func() {
			Expect(con.Open()).To(BeNil())
		})
	})

	Describe("Close", func() {
		It("does not return error", func() {
			Expect(con.Close()).To(BeNil())
		})
	})

	Describe("Write", func() {
		It("does not return error", func() {
			Expect(con.Write("")).To(BeNil())
		})
	})

	Describe("Listen", func() {
		It("does not return error", func() {
			ch := make(chan<- rcon.Event)
			Expect(con.Listen(ch)).To(BeNil())
		})
	})
})
