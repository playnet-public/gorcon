package battleye_test

import (
	"context"

	"github.com/playnet-public/gorcon/pkg/mocks"
	"github.com/playnet-public/gorcon/pkg/rcon"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	be_proto "github.com/playnet-public/battleye/battleye"
	be_mocks "github.com/playnet-public/battleye/mocks"
	be "github.com/playnet-public/gorcon/pkg/rcon/battleye"
)

var _ = Describe("Reader", func() {
	var (
		con *be.Connection
		pr  *be_mocks.Protocol
		udp *mocks.UDPConnection
	)

	BeforeEach(func() {
		con = be.NewConnection(context.Background())
		pr = &be_mocks.Protocol{}
		udp = &mocks.UDPConnection{}
		con.Protocol = pr
		con.UDP = udp
	})

	Describe("HandlePacket", func() {
		BeforeEach(func() {
			pr.VerifyReturns(nil)
			pr.DataReturns([]byte("test"), nil)
			pr.TypeReturns(be_proto.Command, nil)
		})
		It("does return nil", func() {
			pr.TypeReturns(0x12, nil)
			Expect(con.HandlePacket(nil)).To(BeNil())
		})
		It("should call Verify", func() {
			con.HandlePacket([]byte(""))
			Expect(pr.VerifyCallCount()).To(BeEquivalentTo(1))
		})
		It("does call Verify with packet", func() {
			con.HandlePacket([]byte("test"))
			Expect(pr.VerifyArgsForCall(0)).To(BeEquivalentTo(be_proto.Packet("test")))
		})
		It("does return error on invalid packet", func() {
			pr.VerifyReturns(errors.New("test"))
			Expect(con.HandlePacket(nil)).NotTo(BeNil())
		})
		It("does call Data", func() {
			con.HandlePacket([]byte(""))
			Expect(pr.DataCallCount()).To(BeEquivalentTo(1))
		})
		It("should call Data with packet", func() {
			con.HandlePacket([]byte("test"))
			Expect(pr.DataArgsForCall(0)).To(BeEquivalentTo(be_proto.Packet("test")))
		})
		It("does return error on corrupt data", func() {
			pr.DataReturns(nil, errors.New("test"))
			Expect(con.HandlePacket(nil)).NotTo(BeNil())
		})
		It("does call Type", func() {
			con.HandlePacket([]byte("test"))
			Expect(pr.TypeCallCount()).To(BeEquivalentTo(1))
		})
		It("does call Type with packet", func() {
			con.HandlePacket([]byte("test"))
			Expect(pr.TypeArgsForCall(0)).To(BeEquivalentTo(be_proto.Packet("test")))
		})
		It("does return error if Type returns error", func() {
			pr.TypeReturns(0x00, errors.New("test"))
			Expect(con.HandlePacket(nil)).NotTo(BeNil())
		})
		Context("when given empty data", func() {
			It("does increase pingback", func() {
				pb := con.Pingback()
				pr.DataReturns([]byte(""), nil)
				con.HandlePacket([]byte(""))
				Expect(con.Pingback()).To(BeNumerically(">", pb))
			})
		})
		It("does return nil when handling ServerMessage", func() {
			pr.TypeReturns(be_proto.ServerMessage, nil)
			Expect(con.HandlePacket(nil)).To(BeNil())
		})
	})

	Describe("HandleResponse", func() {
		BeforeEach(func() {
			pr.SequenceReturns(0, nil)
			trm := be.NewTransmission("test")
			go func() {
				<-trm.Done()
			}()
			con.AddTransmission(0, trm)
			pr.TypeReturns(be_proto.Command, nil)
			pr.DataReturns([]byte("test data"), nil)
			pr.MultiReturns(2, 1, false)
		})
		It("does not return error", func() {
			Expect(con.HandleResponse(be_proto.Packet("test"))).To(BeNil())
		})
		It("does return error if no transmission is present", func() {
			con.AddTransmission(0, nil)
			Expect(con.HandleResponse(be_proto.Packet("test"))).NotTo(BeNil())
		})
		It("does call Sequence with packet", func() {
			con.HandleResponse([]byte("test"))
			Expect(pr.SequenceCallCount()).To(BeEquivalentTo(1))
			Expect(pr.SequenceArgsForCall(0)).To(BeEquivalentTo(be_proto.Packet("test")))
		})
		It("does return error if Sequence returns error", func() {
			pr.SequenceReturns(0, errors.New("test"))
			Expect(con.HandleResponse(nil)).NotTo(BeNil())
		})
		It("does call Type with packet", func() {
			con.HandleResponse([]byte("test"))
			Expect(pr.TypeCallCount()).To(BeEquivalentTo(1))
			Expect(pr.TypeArgsForCall(0)).To(BeEquivalentTo(be_proto.Packet("test")))
		})
		It("does return error if Type returns error", func() {
			pr.TypeReturns(0x02, errors.New("test"))
			Expect(con.HandleResponse(nil)).NotTo(BeNil())
		})
		It("does call Data with packet", func() {
			con.HandleResponse([]byte("test"))
			Expect(pr.DataCallCount()).To(BeEquivalentTo(1))
			Expect(pr.DataArgsForCall(0)).To(BeEquivalentTo(be_proto.Packet("test")))
		})
		It("does return error if Data returns error", func() {
			pr.DataReturns([]byte(""), errors.New("test"))
			Expect(con.HandleResponse(nil)).NotTo(BeNil())
		})
		Context("on multi command", func() {
			BeforeEach(func() {
				pr.TypeReturns(be_proto.MultiCommand, nil)
			})
			It("does add correct index to buffer", func() {
				pr.DataReturns([]byte("test "), nil)
				pr.MultiReturns(2, 0, false)
				con.HandleResponse(nil)
				pr.DataReturns([]byte("data"), nil)
				pr.MultiReturns(2, 1, false)
				con.HandleResponse(nil)
				trm := con.GetTransmission(0)
				Expect(trm.Response()).To(BeEquivalentTo("test data"))
			})
		})
		Context("on single command", func() {
			BeforeEach(func() {
				pr.TypeReturns(be_proto.Command, nil)
			})
			It("does add correct index to buffer", func() {
				pr.DataReturns([]byte("test data"), nil)
				con.HandleResponse(nil)
				trm := con.GetTransmission(0)
				Expect(trm.Response()).To(BeEquivalentTo("test data"))
			})
		})
		It("does send to done channel", func() {
			trm := be.NewTransmission("test")
			con.AddTransmission(0, trm)
			go func() {
				pr.TypeReturns(be_proto.Command, nil)
				pr.DataReturns([]byte("test data"), nil)
				con.HandleResponse(nil)
			}()
			Expect(<-trm.Done()).To(BeTrue())
		})
		It("does timeout on blocking done channel", func() {
			trm := be.NewTransmission("test")
			con.AddTransmission(0, trm)
			Expect(con.HandleResponse(nil)).To(BeNil())
		})
	})

	Describe("HandleServerMessage", func() {
		BeforeEach(func() {
			pr.SequenceReturns(0, nil)
			udp.WriteReturns(0, nil)
		})
		It("does return nil", func() {
			Expect(con.HandleServerMessage(be_proto.Packet("test"))).To(BeNil())
		})
		It("does call Sequence with packet", func() {
			con.HandleServerMessage([]byte("test"))
			Expect(pr.SequenceCallCount()).To(BeEquivalentTo(1))
			Expect(pr.SequenceArgsForCall(0)).To(BeEquivalentTo(be_proto.Packet("test")))
		})
		It("does return error if Sequence returns error", func() {
			pr.SequenceReturns(0, errors.New("test"))
			Expect(con.HandleServerMessage(nil)).NotTo(BeNil())
		})
		It("does send event to channel", func() {
			c := make(chan *rcon.Event)
			con.Listen(c)
			con.HandleServerMessage([]byte("test"))
			event := <-c
			Expect(event.Message).NotTo(BeEquivalentTo(""))
		})
		It("does set correct type when handling chat event", func() {
			c := make(chan *rcon.Event)
			con.Listen(c)
			con.HandleServerMessage([]byte("(Group) Test"))
			event := <-c
			Expect(event.Type).To(BeEquivalentTo(rcon.TypeChat))
		})
		It("does return error if UDP.Write fails", func() {
			udp.WriteReturns(0, errors.New("test"))
			Expect(con.HandleServerMessage([]byte("test"))).NotTo(BeNil())
		})
	})
})
