package battleye_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	be_proto "github.com/playnet-public/battleye/battleye"
	be_mocks "github.com/playnet-public/battleye/mocks"
	be "github.com/playnet-public/gorcon/pkg/rcon/battleye"
)

var _ = Describe("Reader", func() {
	var (
		con *be.Connection
		pr  *be_mocks.Protocol
	)

	BeforeEach(func() {
		con = be.NewConnection(context.Background())
		pr = &be_mocks.Protocol{}
		con.Protocol = pr
	})

	Describe("HandlePacket", func() {
		It("should call VerifyPacket", func() {
			con.HandlePacket([]byte(""))
			Expect(pr.VerifyCallCount()).To(BeEquivalentTo(1))
		})
		It("should call VerifyPacket with data", func() {
			con.HandlePacket([]byte("test"))
			Expect(pr.VerifyArgsForCall(0)).To(BeEquivalentTo(be_proto.Packet("test")))
		})
	})
})
