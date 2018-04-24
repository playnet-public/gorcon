package battleye_test

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"gopkg.in/tomb.v2"

	"github.com/playnet-public/battleye/protocol"

	"github.com/playnet-public/gorcon/pkg/mocks"

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
		c   *be.Client
		ctx context.Context
	)

	BeforeEach(func() {
		c = &be.Client{}
		ctx = context.Background()
	})

	Describe("NewConnection", func() {
		It("does not return nil", func() {
			Expect(c.NewConnection(ctx)).NotTo(BeNil())
		})
	})
})

var _ = Describe("Connection", func() {
	var (
		con   *be.Connection
		dial  *mocks.UDPDialer
		udp   *mocks.UDPConnection
		proto *mocks.BattlEyeProtocol
		ctx   context.Context
	)

	BeforeEach(func() {
		dial = &mocks.UDPDialer{}
		proto = &mocks.BattlEyeProtocol{}
		ctx = context.Background()
		con = be.NewConnection(ctx)
		con.Dialer = dial
		con.Protocol = proto

		udp = &mocks.UDPConnection{}
		dial.DialUDPReturns(udp, nil)
	})

	Describe("Open", func() {
		BeforeEach(func() {
			con.Password = "test"
			proto.VerifyLoginReturns(protocol.PacketResponse.LoginOk, nil)
		})
		It("does not return error", func() {
			Expect(con.Open()).To(BeNil())
		})
		It("returns an error if there already is a udp connection", func() {
			con.UDP = &net.UDPConn{}
			Expect(con.Open()).NotTo(BeNil())
		})
		It("calls DialUDP once", func() {
			con.Open()
			Expect(dial.DialUDPCallCount()).To(BeEquivalentTo(1))
		})
		It("calls DialUDP with the correct address", func() {
			con.Addr, _ = net.ResolveUDPAddr("udp", "127.0.0.1:8080")
			con.Open()
			_, _, addr := dial.DialUDPArgsForCall(0)
			Expect(addr).To(BeEquivalentTo(con.Addr))
		})
		It("is setting the udp connection", func() {
			con.Open()
			Expect(con.UDP).NotTo(BeNil())
		})
		It("calls deadline setters", func() {
			con.Open()
			Expect(udp.SetReadDeadlineCallCount()).To(BeEquivalentTo(1))
			Expect(udp.SetWriteDeadlineCallCount()).To(BeEquivalentTo(1))
		})
		It("returns error if dial fails", func() {
			dial.DialUDPReturns(nil, errors.New("test"))
			Expect(con.Open()).NotTo(BeNil())
		})
		It("does send a login packet", func() {
			con.Open()
			args := udp.WriteArgsForCall(0)
			Expect(args).To(BeEquivalentTo(proto.BuildLoginPacket("test")))
		})
		It("does use the stored credentials for building login packets", func() {
			con.Password = "password"
			con.Open()
			args := udp.WriteArgsForCall(0)
			Expect(args).To(BeEquivalentTo(proto.BuildLoginPacket("password")))
		})
		It("does return error if sending login packet fails", func() {
			udp.WriteReturns(0, errors.New("test"))
			Expect(con.Open()).NotTo(BeNil())
		})
		It("does call read after sending login", func() {
			con.Open()
			Expect(udp.ReadCallCount()).To(BeEquivalentTo(1))
		})
		It("does return error if reading from udp fails", func() {
			udp.ReadReturns(0, errors.New("test"))
			Expect(con.Open()).NotTo(BeNil())
		})
		It("does return error on invalid login response", func() {
			proto.VerifyLoginReturns(0, errors.New("test"))
			Expect(con.Open()).NotTo(BeNil())
		})
		It("does return error on invalid login credentials", func() {
			proto.VerifyLoginReturns(protocol.PacketResponse.LoginFail, nil)
			Expect(con.Open()).NotTo(BeNil())
		})
	})

	Describe("WriterLoop", func() {
		BeforeEach(func() {
			con.UDP = udp
			con.KeepAliveTimeout = 0
		})
		It("does send at least one keepAlive packet", func() {
			con.KeepAliveTimeout = 0
			con.Tomb.Go(con.WriterLoop)
			<-time.After(time.Millisecond * 2)
			Expect(con.Tomb.Err()).To(BeEquivalentTo(tomb.ErrStillAlive))
			Expect(udp.WriteCallCount()).To(BeNumerically(">", 0))
			con.Close()
		})
		It("does exit on close", func() {
			con.KeepAliveTimeout = 100
			go func() {
				time.Sleep(time.Millisecond * 5)
				con.Close()
			}()
			Expect(con.WriterLoop()).To(BeEquivalentTo(tomb.ErrDying))

		})
		It("does return error if udp is nil", func() {
			con.UDP = nil
			Expect(con.WriterLoop()).NotTo(BeNil())
		})
	})

	Describe("ReaderLoop", func() {
		BeforeEach(func() {
			con.UDP = udp
			con.KeepAliveTimeout = 0
		})
		It("does return error if udp is nil", func() {
			con.UDP = nil
			Expect(con.ReaderLoop()).NotTo(BeNil())
		})
		It("does not return on timeout", func() {
			udp.ReadReturns(0, &timeoutError{})
			con.Tomb.Go(con.ReaderLoop)
			<-time.After(time.Millisecond * 2)
			Expect(con.Tomb.Err()).To(BeEquivalentTo(tomb.ErrStillAlive))
			Expect(udp.ReadCallCount()).To(BeNumerically(">", 0))
			con.Close()
		})
		It("does return on non-timeout error", func() {
			udp.ReadReturns(0, errors.New("test"))
			Expect(con.ReaderLoop()).NotTo(BeNil())
		})
	})

	Describe("Close", func() {
		BeforeEach(func() {
			con.UDP = udp
			con.Hold()
		})
		It("does not return error", func() {
			Expect(con.Close()).To(BeNil())
		})
		It("does return error if udp connection is nil", func() {
			con.UDP = nil
			Expect(con.Close()).NotTo(BeNil())
		})
		It("calls close on the udp connection", func() {
			con.Close()
			Expect(udp.CloseCallCount()).To(BeEquivalentTo(1))
		})
		It("does return error if udp close fails", func() {
			udp.CloseReturns(errors.New("test"))
			Expect(con.Close()).NotTo(BeNil())
		})
		It("does reset the udp after closing", func() {
			con.Close()
			Expect(con.UDP).To(BeNil())
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

type timeoutError struct {
	Err error
}

func (t *timeoutError) Error() string   { return t.Err.Error() }
func (t *timeoutError) Timeout() bool   { return true }
func (t *timeoutError) Temporary() bool { return false }
