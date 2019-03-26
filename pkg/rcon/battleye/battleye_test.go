package battleye_test

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"gopkg.in/tomb.v2"

	be_proto "github.com/playnet-public/battleye/battleye"
	be_mocks "github.com/playnet-public/battleye/mocks"
	"github.com/playnet-public/gorcon/pkg/mocks"
	be "github.com/playnet-public/gorcon/pkg/rcon/battleye"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBattlEye(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BattlEye Suite")
}

var _ = Describe("Client", func() {
	var (
		ctx context.Context
		c   *be.Client
	)

	BeforeEach(func() {
		ctx = context.Background()
		c = be.New(ctx)
	})

	Describe("NewConnection", func() {
		It("does not return nil", func() {
			Expect(c.NewConnection(ctx)).NotTo(BeNil())
		})
	})
})

var _ = Describe("Connection", func() {
	setup := func() (ctx context.Context, c *be.Client, con *be.Connection, dial *mocks.UDPDialer, proto *be_mocks.Protocol, udp *mocks.UDPConnection) {
		ctx = context.Background()
		c = be.New(ctx)
		con = c.NewConnection(ctx).(*be.Connection)

		dial = &mocks.UDPDialer{}
		con.Dialer = dial

		proto = &be_mocks.Protocol{}
		con.Protocol = proto

		udp = &mocks.UDPConnection{}
		dial.DialUDPReturns(udp, nil)
		con.UDP = udp

		return
	}

	Describe("Open", func() {
		It("does not return error", func() {
			ctx, _, con, dial, _, _ := setup()
			con.UDP = nil
			con.Dialer = dial
			Expect(con.Open(ctx)).To(BeNil())
		})
		It("returns an error if there already is a udp connection", func() {
			ctx, _, con, _, _, _ := setup()
			con.UDP = &net.UDPConn{}
			Expect(con.Open(ctx)).NotTo(BeNil())
		})
		It("calls DialUDP once", func() {
			ctx, _, con, dial, _, _ := setup()
			con.UDP = nil
			con.Dialer = dial
			Expect(con.Open(ctx)).To(BeNil())
			Expect(dial.DialUDPCallCount()).To(BeEquivalentTo(1))
		})
		It("calls DialUDP with the correct address", func() {
			ctx, _, con, dial, _, _ := setup()
			con.UDP = nil
			con.Dialer = dial
			con.Addr, _ = net.ResolveUDPAddr("udp", "127.0.0.1:8080")
			Expect(con.Open(ctx)).To(BeNil())
			_, _, addr := dial.DialUDPArgsForCall(0)
			Expect(addr).To(BeEquivalentTo(con.Addr))
		})
		It("is setting the udp connection", func() {
			ctx, _, con, dial, _, _ := setup()
			con.UDP = nil
			con.Dialer = dial
			Expect(con.Open(ctx)).To(BeNil())
			Expect(con.UDP).NotTo(BeNil())
		})
		It("calls deadline setters", func() {
			ctx, _, con, dial, _, udp := setup()
			con.UDP = nil
			con.Dialer = dial
			Expect(con.Open(ctx)).To(BeNil())
			Expect(udp.SetReadDeadlineCallCount()).To(BeEquivalentTo(1))
			Expect(udp.SetWriteDeadlineCallCount()).To(BeEquivalentTo(1))
		})
		It("returns error if dial fails", func() {
			ctx, _, con, dial, _, _ := setup()
			con.UDP = nil
			con.Dialer = dial
			dial.DialUDPReturns(nil, errors.New("test"))
			Expect(con.Open(ctx)).NotTo(BeNil())
		})
		It("does send a login packet", func() {
			ctx, _, con, dial, _, udp := setup()
			con.UDP = nil
			con.Dialer = dial
			Expect(con.Open(ctx)).To(BeNil())
			args := udp.WriteArgsForCall(0)
			Expect(args).To(BeEquivalentTo(con.Protocol.BuildLoginPacket("test")))
		})
		It("does use the stored credentials for building login packets", func() {
			ctx, _, con, dial, _, udp := setup()
			con.UDP = nil
			con.Dialer = dial
			con.Password = "password"
			Expect(con.Open(ctx)).To(BeNil())
			args := udp.WriteArgsForCall(0)
			Expect(args).To(BeEquivalentTo(con.Protocol.BuildLoginPacket("password")))
		})
		It("does return error if sending login packet fails", func() {
			ctx, _, con, dial, _, udp := setup()
			con.UDP = nil
			con.Dialer = dial
			udp.WriteReturns(0, errors.New("test"))
			Expect(con.Open(ctx)).NotTo(BeNil())
		})
		It("does call read after sending login", func() {
			ctx, _, con, dial, _, udp := setup()
			con.UDP = nil
			con.Dialer = dial
			Expect(con.Open(ctx)).To(BeNil())
			Expect(udp.ReadCallCount()).To(BeNumerically(">", 0))
		})
		It("does return error if reading from udp fails", func() {
			ctx, _, con, _, _, udp := setup()
			udp.ReadReturns(0, errors.New("test"))
			Expect(con.Open(ctx)).NotTo(BeNil())
		})
		It("does return error on invalid login response", func() {
			ctx, _, con, _, proto, _ := setup()
			proto.VerifyLoginReturns(errors.New("test"))
			Expect(con.Open(ctx)).NotTo(BeNil())
		})
		It("does return error on invalid login credentials", func() {
			ctx, _, con, _, proto, _ := setup()
			proto.VerifyLoginReturns(errors.New("test"))
			Expect(con.Open(ctx)).NotTo(BeNil())
		})
	})

	Describe("WriterLoop", func() {
		It("does send at least one keepAlive packet", func() {
			ctx, _, con, _, _, udp := setup()
			con.UDP = udp
			con.KeepAliveTimeout = 0
			con.Tomb.Go(con.WriterLoop(ctx))
			<-time.After(time.Second * 1)
			Expect(con.Tomb.Err()).To(BeEquivalentTo(tomb.ErrStillAlive))
			Expect(udp.WriteCallCount()).To(BeNumerically(">", 0))
			con.Close(ctx)
		})
		It("does exit on close", func() {
			ctx, _, con, _, _, udp := setup()
			con.UDP = udp
			con.KeepAliveTimeout = 100
			go func() {
				time.Sleep(time.Second * 1)
				con.Close(ctx)
			}()
			Expect(con.WriterLoop(ctx)()).To(BeEquivalentTo(tomb.ErrDying))

		})
		It("does return error if udp is nil", func() {
			ctx, _, con, _, _, _ := setup()
			con.UDP = nil
			con.KeepAliveTimeout = 0
			Expect(con.WriterLoop(ctx)()).NotTo(BeNil())
		})
	})

	Describe("ReaderLoop", func() {
		It("does return error if udp is nil", func() {
			ctx, _, con, _, _, _ := setup()
			con.KeepAliveTimeout = 0
			con.UDP = nil
			Expect(con.ReaderLoop(ctx)()).NotTo(BeNil())
		})
		It("does not return on timeout", func() {
			ctx, _, con, _, _, udp := setup()
			con.UDP = udp
			con.KeepAliveTimeout = 0

			udp.ReadReturns(0, &timeoutError{})
			con.Tomb.Go(con.ReaderLoop(ctx))
			<-time.After(time.Second * 1)
			Expect(con.Tomb.Err()).To(BeEquivalentTo(tomb.ErrStillAlive))
			Expect(udp.ReadCallCount()).To(BeNumerically(">", 0))
			con.Close(ctx)
		})
		It("does return on non-timeout error", func() {
			ctx, _, con, _, _, udp := setup()
			con.UDP = udp
			con.KeepAliveTimeout = 0

			udp.ReadReturns(0, errors.New("test"))
			Expect(con.ReaderLoop(ctx)()).NotTo(BeNil())
		})
	})

	Describe("Close", func() {
		It("does not return error", func() {
			ctx, _, con, _, _, udp := setup()
			con.UDP = udp
			con.Hold(ctx)
			Expect(con.Close(ctx)).To(BeNil())
		})
		It("does return error if udp connection is nil", func() {
			ctx, _, con, _, _, udp := setup()
			con.UDP = udp
			con.Tomb.Go(func() error {
				for {
					<-con.Tomb.Dying()
					return nil
				}
			})
			con.UDP = nil
			Expect(con.Close(ctx)).NotTo(BeNil())
		})
		It("calls close on the udp connection", func() {
			ctx, _, con, _, _, udp := setup()
			con.UDP = udp
			con.Hold(ctx)
			con.Close(ctx)
			Expect(udp.CloseCallCount()).To(BeEquivalentTo(1))
		})
		It("does return error if udp close fails", func() {
			ctx, _, con, _, _, udp := setup()
			con.UDP = udp
			con.Hold(ctx)
			udp.CloseReturns(errors.New("test"))
			Expect(con.Close(ctx)).NotTo(BeNil())
		})
		It("does reset the udp after closing", func() {
			ctx, _, con, _, _, udp := setup()
			con.UDP = udp
			con.Hold(ctx)
			con.Close(ctx)
			Expect(con.UDP).To(BeNil())
		})
	})

	Describe("Write", func() {
		It("does not return error", func() {
			ctx, _, con, _, proto, _ := setup()
			proto.BuildCmdPacketStub = be_proto.New().BuildCmdPacket
			con.ResetSequence()
			_, err := con.Write(ctx, "")
			Expect(err).To(BeNil())
		})
		It("does return error if udp connection is nil", func() {
			ctx, _, con, _, proto, _ := setup()
			proto.BuildCmdPacketStub = be_proto.New().BuildCmdPacket
			con.ResetSequence()
			con.UDP = nil
			_, err := con.Write(ctx, "")
			Expect(err).NotTo(BeNil())
		})
		It("does call con.Write", func() {
			ctx, _, con, _, proto, udp := setup()
			proto.BuildCmdPacketStub = be_proto.New().BuildCmdPacket
			con.ResetSequence()
			_, err := con.Write(ctx, "test")
			Expect(err).To(BeNil())
			Expect(udp.WriteCallCount()).To(BeEquivalentTo(1))
		})
		It("does return error on failed write", func() {
			ctx, _, con, _, proto, udp := setup()
			proto.BuildCmdPacketStub = be_proto.New().BuildCmdPacket
			con.ResetSequence()
			udp.WriteReturns(0, errors.New("test"))
			_, err := con.Write(ctx, "")
			Expect(err).NotTo(BeNil())
		})
		It("does write correct command packet", func() {
			ctx, _, con, _, proto, udp := setup()
			con.UDP = udp
			proto.BuildCmdPacketStub = be_proto.New().BuildCmdPacket
			con.ResetSequence()
			_, err := con.Write(ctx, "test")
			Expect(err).To(BeNil())
			Expect(udp.WriteArgsForCall(0)).To(BeEquivalentTo(con.Protocol.BuildCmdPacket([]byte("test"), 1)))
		})
		It("does increase sequence after write", func() {
			ctx, _, con, _, proto, udp := setup()
			con.UDP = udp
			proto.BuildCmdPacketStub = be_proto.New().BuildCmdPacket
			con.ResetSequence()
			seq := con.Sequence()
			_, err := con.Write(ctx, "")
			Expect(err).To(BeNil())
			Expect(con.Sequence() == seq+1).To(BeTrue())
		})
		It("does add transmission to connection at write", func() {
			ctx, _, con, _, proto, udp := setup()
			con.UDP = udp
			proto.BuildCmdPacketStub = be_proto.New().BuildCmdPacket
			con.ResetSequence()
			_, err := con.Write(ctx, "test")
			Expect(err).To(BeNil())
			Expect(con.GetTransmission(1)).NotTo(BeNil())
		})
	})
})

type timeoutError struct {
	Err error
}

func (t *timeoutError) Error() string   { return t.Err.Error() }
func (t *timeoutError) Timeout() bool   { return true }
func (t *timeoutError) Temporary() bool { return false }
