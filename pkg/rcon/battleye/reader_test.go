package battleye_test

import (
	"github.com/playnet-public/gorcon/pkg/mocks"
	"github.com/playnet-public/gorcon/pkg/rcon"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	be_proto "github.com/playnet-public/battleye/battleye"
	be_mocks "github.com/playnet-public/battleye/mocks"
	be "github.com/playnet-public/gorcon/pkg/rcon/battleye"
	context "github.com/seibert-media/golibs/log"
)

var _ = Describe("Reader", func() {
	var (
		con *be.Connection
		pr  *be_mocks.Protocol
		udp *mocks.UDPConnection
		ctx context.Context
	)

	BeforeEach(func() {
		con = be.NewConnection(context.Background())
		pr = &be_mocks.Protocol{}
		udp = &mocks.UDPConnection{}
		con.Protocol = pr
		con.UDP = udp
		ctx = context.NewNop()
	})

	Describe("HandlePacket", func() {
		BeforeEach(func() {
			pr.VerifyReturns(nil)
			pr.DataReturns([]byte("test"), nil)
			pr.TypeReturns(be_proto.Command, nil)
		})
		It("does return nil", func() {
			pr.TypeReturns(0x12, nil)
			Expect(con.HandlePacket(ctx, nil)).To(BeNil())
		})
		It("should call Verify", func() {
			con.HandlePacket(ctx, []byte(""))
			Expect(pr.VerifyCallCount()).To(BeEquivalentTo(1))
		})
		It("does call Verify with packet", func() {
			con.HandlePacket(ctx, []byte("test"))
			Expect(pr.VerifyArgsForCall(0)).To(BeEquivalentTo(be_proto.Packet("test")))
		})
		It("does return error on invalid packet", func() {
			pr.VerifyReturns(errors.New("test"))
			Expect(con.HandlePacket(ctx, nil)).NotTo(BeNil())
		})
		It("does call Data", func() {
			con.HandlePacket(ctx, []byte(""))
			Expect(pr.DataCallCount()).To(BeEquivalentTo(1))
		})
		It("should call Data with packet", func() {
			con.HandlePacket(ctx, []byte("test"))
			Expect(pr.DataArgsForCall(0)).To(BeEquivalentTo(be_proto.Packet("test")))
		})
		It("does return error on corrupt data", func() {
			pr.DataReturns(nil, errors.New("test"))
			Expect(con.HandlePacket(ctx, nil)).NotTo(BeNil())
		})
		It("does call Type", func() {
			con.HandlePacket(ctx, []byte("test"))
			Expect(pr.TypeCallCount()).To(BeEquivalentTo(1))
		})
		It("does call Type with packet", func() {
			con.HandlePacket(ctx, []byte("test"))
			Expect(pr.TypeArgsForCall(0)).To(BeEquivalentTo(be_proto.Packet("test")))
		})
		It("does return error if Type returns error", func() {
			pr.TypeReturns(0x00, errors.New("test"))
			Expect(con.HandlePacket(ctx, nil)).NotTo(BeNil())
		})
		Context("when given empty data", func() {
			It("does increase pingback", func() {
				pb := con.Pingback()
				pr.DataReturns([]byte(""), nil)
				con.HandlePacket(ctx, []byte(""))
				Expect(con.Pingback()).To(BeNumerically(">", pb))
			})
		})
		It("does return nil when handling ServerMessage", func() {
			pr.TypeReturns(be_proto.ServerMessage, nil)
			Expect(con.HandlePacket(ctx, nil)).To(BeNil())
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
			Expect(con.HandleResponse(ctx, be_proto.Packet("test"))).To(BeNil())
		})
		It("does return error if no transmission is present", func() {
			con.AddTransmission(0, nil)
			Expect(con.HandleResponse(ctx, be_proto.Packet("test"))).NotTo(BeNil())
		})
		It("does call Sequence with packet", func() {
			con.HandleResponse(ctx, []byte("test"))
			Expect(pr.SequenceCallCount()).To(BeEquivalentTo(1))
			Expect(pr.SequenceArgsForCall(0)).To(BeEquivalentTo(be_proto.Packet("test")))
		})
		It("does return error if Sequence returns error", func() {
			pr.SequenceReturns(0, errors.New("test"))
			Expect(con.HandleResponse(ctx, nil)).NotTo(BeNil())
		})
		It("does call Type with packet", func() {
			con.HandleResponse(ctx, []byte("test"))
			Expect(pr.TypeCallCount()).To(BeEquivalentTo(1))
			Expect(pr.TypeArgsForCall(0)).To(BeEquivalentTo(be_proto.Packet("test")))
		})
		It("does return error if Type returns error", func() {
			pr.TypeReturns(0x02, errors.New("test"))
			Expect(con.HandleResponse(ctx, nil)).NotTo(BeNil())
		})
		It("does call Data with packet", func() {
			con.HandleResponse(ctx, []byte("test"))
			Expect(pr.DataCallCount()).To(BeEquivalentTo(1))
			Expect(pr.DataArgsForCall(0)).To(BeEquivalentTo(be_proto.Packet("test")))
		})
		It("does return error if Data returns error", func() {
			pr.DataReturns([]byte(""), errors.New("test"))
			Expect(con.HandleResponse(ctx, nil)).NotTo(BeNil())
		})
		Context("on multi command", func() {
			BeforeEach(func() {
				pr.TypeReturns(be_proto.MultiCommand, nil)
			})
			It("does add correct index to buffer", func() {
				pr.DataReturns([]byte("test "), nil)
				pr.MultiReturns(2, 0, false)
				con.HandleResponse(ctx, nil)
				pr.DataReturns([]byte("data"), nil)
				pr.MultiReturns(2, 1, false)
				con.HandleResponse(ctx, nil)
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
				con.HandleResponse(ctx, nil)
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
				con.HandleResponse(ctx, nil)
			}()
			Expect(<-trm.Done()).To(BeTrue())
		})
		It("does timeout on blocking done channel", func() {
			trm := be.NewTransmission("test")
			con.AddTransmission(0, trm)
			Expect(con.HandleResponse(ctx, nil)).To(BeNil())
		})
	})

	Describe("HandleServerMessage", func() {
		BeforeEach(func() {
			pr.SequenceReturns(0, nil)
			udp.WriteReturns(0, nil)
		})
		It("does return nil", func() {
			Expect(con.HandleServerMessage(ctx, be_proto.Packet("test"))).To(BeNil())
		})
		It("does call Sequence with packet", func() {
			con.HandleServerMessage(ctx, []byte("test"))
			Expect(pr.SequenceCallCount()).To(BeEquivalentTo(1))
			Expect(pr.SequenceArgsForCall(0)).To(BeEquivalentTo(be_proto.Packet("test")))
		})
		It("does return error if Sequence returns error", func() {
			pr.SequenceReturns(0, errors.New("test"))
			Expect(con.HandleServerMessage(ctx, nil)).NotTo(BeNil())
		})
		It("does send event to channel", func() {
			c := make(chan *rcon.Event)
			con.Subscribe(ctx, c)
			con.HandleServerMessage(ctx, []byte("test"))
			event := <-c
			Expect(event.Payload).NotTo(BeEquivalentTo(""))
		})
		It("does set correct type when handling chat event", func() {
			c := make(chan *rcon.Event)
			con.Subscribe(ctx, c)
			con.HandleServerMessage(ctx, []byte("(Group) Test"))
			event := <-c
			Expect(event.Type).To(BeEquivalentTo(rcon.TypeChat))
		})
		It("does return error if UDP.Write fails", func() {
			udp.WriteReturns(0, errors.New("test"))
			Expect(con.HandleServerMessage(ctx, []byte("test"))).NotTo(BeNil())
		})
	})
})
