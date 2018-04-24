package battleye

import (
	"context"
	"net"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	tomb "gopkg.in/tomb.v2"

	be "github.com/playnet-public/battleye/protocol"
	"github.com/playnet-public/gorcon/pkg/rcon"
)

// Client is a BattlEye specific implementation of rcon.Client to create new BattlEye rcon connections
type Client struct {
}

// NewConnection from the current client's configuration
func (c *Client) NewConnection(ctx context.Context) rcon.Connection {
	return NewConnection(ctx)
}

// Connection is a BattlEye specific implementation of rcon.Connection offering all required rcon generics
type Connection struct {
	Addr     *net.UDPAddr
	Password string
	Dialer   udpDialer

	UDP      UDPConnection
	Protocol protocol

	KeepAliveTimeout int
	keepAliveCount   int64
	seq              uint32
	pingbackCount    int64

	errors chan error
	Tomb   *tomb.Tomb
}

// NewConnection from the passed in configuration
func NewConnection(ctx context.Context) *Connection {
	c := &Connection{
		errors: make(chan error),
	}
	atomic.StoreUint32(&c.seq, 0)
	atomic.StoreInt64(&c.keepAliveCount, 0)
	atomic.StoreInt64(&c.pingbackCount, 0)
	c.Tomb, ctx = tomb.WithContext(ctx)
	return c
}

//go:generate counterfeiter -o ../../mocks/udp_dialer.go --fake-name UDPDialer . udpDialer
type udpDialer interface {
	DialUDP(string, *net.UDPAddr, *net.UDPAddr) (UDPConnection, error)
}

//go:generate counterfeiter -o ../../mocks/battleye_protocol.go --fake-name BattlEyeProtocol . protocol
type protocol interface {
	BuildLoginPacket(string) []byte
	VerifyLogin([]byte) (byte, error)
}

// UDPConnection interface defines all udp functions required and is used primarily for mocking
//go:generate counterfeiter -o ../../mocks/udp_connection.go --fake-name UDPConnection . UDPConnection
type UDPConnection interface {
	Close() error
	Read([]byte) (int, error)
	Write([]byte) (int, error)

	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
}

// Open the connection
func (c *Connection) Open() error {
	if c.UDP != nil {
		return errors.New("connection already open")
	}
	udp, err := c.Dialer.DialUDP("udp", nil, c.Addr)
	if err != nil {
		return errors.Wrap(err, "dialing udp failed")
	}
	c.UDP = udp
	c.UDP.SetReadDeadline(time.Now().Add(time.Second * 2)) // TODO: Evaluate if this is required
	c.UDP.SetWriteDeadline(time.Now().Add(time.Millisecond * 100))

	buf := make([]byte, 9)
	_, err = c.UDP.Write(c.Protocol.BuildLoginPacket(c.Password))
	if err != nil {
		return errors.Wrap(err, "sending login packet failed")
	}

	n, err := c.UDP.Read(buf)
	if err != nil {
		return errors.Wrap(err, "reading login response failed")
	}

	resp, err := c.Protocol.VerifyLogin(buf[:n])
	if err != nil {
		return errors.Wrap(err, "verifying login response failed")
	}
	if resp == be.PacketResponse.LoginFail {
		return errors.New("logging in failed with invalid credentials")
	}
	c.Hold()
	return nil
}

// Hold the connection by sending keepalive packets as required by the battleye protocol
func (c *Connection) Hold() {
	c.Tomb.Go(c.WriterLoop)
	c.Tomb.Go(c.ReaderLoop)
}

// WriterLoop for keeping the connection alive
func (c *Connection) WriterLoop() error {
	for {
		select {
		case <-c.Tomb.Dying():
			return tomb.ErrDying
		case <-time.After(time.Second * time.Duration(c.KeepAliveTimeout)):
			if c.UDP != nil {
				c.UDP.Write(be.BuildKeepAlivePacket(c.Sequence()))
				c.AddKeepAlive()
			}
			return errors.New("udp connection must not be nil")
		}
	}
}

// ReaderLoop for keeping the connection alive
func (c *Connection) ReaderLoop() error {
	for {
		select {
		case <-c.Tomb.Dying():
			return tomb.ErrDying
		default:
			if c.UDP != nil {
				buf := make([]byte, 4096)
				_, err := c.UDP.Read(buf)
				if err, ok := err.(net.Error); ok && err.Timeout() {
					// TODO: add debug logging here
					continue
				}
				if err != nil {
					return errors.Wrap(err, "reading udp failed")
				}
			}
			return errors.New("udp connection must not be nil")
		}
	}
}

// Close the connection for graceful shutdown or reconnect
func (c *Connection) Close() error {
	c.Tomb.Kill(errors.New("SIGCLOSE"))
	c.Tomb.Wait()
	if c.UDP == nil {
		return errors.New("connection must not be nil")
	}
	if err := c.UDP.Close(); err != nil {
		return errors.Wrap(err, "closing udp failed")
	}
	// TODO: verify if this should happen always because not resetting will make Open() impossible
	c.UDP = nil
	return nil
}

// Write a command to the connection
func (c *Connection) Write(string) error {
	return nil
}

// Listen for events on the connection. This is a blocking call sending on the passed in channel and returning once an error occurs
func (c *Connection) Listen(chan<- rcon.Event) error {
	return nil
}
