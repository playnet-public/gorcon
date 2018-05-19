package battleye_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	be "github.com/playnet-public/gorcon/pkg/rcon/battleye"
)

var _ = Describe("Connection Helpers", func() {
	var (
		con *be.Connection
	)

	BeforeEach(func() {
		con = be.NewConnection(context.Background())
	})

	Describe("Sequence", func() {
		BeforeEach(func() {
			con.ResetSequence()
		})
		It("should increase the sequence", func() {
			Expect(con.Sequence()).To(BeEquivalentTo(0))
			con.AddSequence()
			Expect(con.Sequence()).To(BeEquivalentTo(1))
			con.ResetSequence()
			Expect(con.Sequence()).To(BeEquivalentTo(0))
		})
		It("should return 1 when calling Add", func() {
			Expect(con.AddSequence()).To(BeEquivalentTo(1))
		})
	})

	Describe("Pingback", func() {
		BeforeEach(func() {
			con.ResetPingback()
		})
		It("should increase and reset pingback", func() {
			Expect(con.Pingback()).To(BeEquivalentTo(0))
			con.AddPingback()
			Expect(con.Pingback()).To(BeEquivalentTo(1))
			con.ResetPingback()
			Expect(con.Pingback()).To(BeEquivalentTo(0))
		})
		It("should return 1 when calling Add", func() {
			Expect(con.AddPingback()).To(BeEquivalentTo(1))
		})
	})
	Describe("KeepAlive", func() {
		BeforeEach(func() {
			con.ResetKeepAlive()
		})
		It("should increase and reset keepalive", func() {
			Expect(con.KeepAlive()).To(BeEquivalentTo(0))
			con.AddKeepAlive()
			Expect(con.KeepAlive()).To(BeEquivalentTo(1))
			con.ResetKeepAlive()
			Expect(con.KeepAlive()).To(BeEquivalentTo(0))
		})
		It("should return 1 when calling Add", func() {
			Expect(con.AddKeepAlive()).To(BeEquivalentTo(1))
		})
	})

	Describe("Transmission", func() {
		It("should return nil on invalid sequence", func() {
			Expect(con.GetTransmission(999)).To(BeNil())
		})
		It("should return valid transmission if present", func() {
			con.AddTransmission(0, be.NewTransmission("test"))
			Expect(con.GetTransmission(0)).NotTo(BeNil())
		})
		It("should remove said transmission on deltet", func() {
			con.AddTransmission(0, be.NewTransmission("test"))
			con.DeleteTransmission(0)
			Expect(con.GetTransmission(0)).To(BeNil())
		})
	})
})
