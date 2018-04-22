package battleye_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	be "github.com/playnet-public/gorcon/pkg/rcon/battleye"
)

var _ = Describe("Connection Helpers", func() {
	var (
		con *be.Connection
	)

	BeforeEach(func() {
		con = be.NewConnection()
	})

	Describe("Sequence", func() {
		It("should increase the sequence", func() {
			Expect(con.Sequence()).To(BeEquivalentTo(0))
			con.AddSequence()
			Expect(con.Sequence()).To(BeEquivalentTo(1))
		})
	})

	Describe("Pingback", func() {
		It("should increase and reset pingback", func() {
			Expect(con.Pingback()).To(BeEquivalentTo(0))
			con.AddPingback()
			Expect(con.Pingback()).To(BeEquivalentTo(1))
			con.ResetPingback()
			Expect(con.Pingback()).To(BeEquivalentTo(0))
		})
	})
	Describe("KeepAlive", func() {
		It("should increase and reset keepalive", func() {
			Expect(con.KeepAlive()).To(BeEquivalentTo(0))
			con.AddKeepAlive()
			Expect(con.KeepAlive()).To(BeEquivalentTo(1))
			con.ResetKeepAlive()
			Expect(con.KeepAlive()).To(BeEquivalentTo(0))
		})
	})
})
