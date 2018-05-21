package battleye

import (
	"net"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"

	"github.com/pkg/errors"
	context "github.com/seibert-media/golibs/log"
	tomb "gopkg.in/tomb.v2"

	be_proto "github.com/playnet-public/battleye/battleye"
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
	Protocol be_proto.Protocol

	KeepAliveTimeout    int
	keepAliveCount      int64
	seq                 uint32
	pingbackCount       int64
	transmissions       map[be_proto.Sequence]*Transmission
	transmissionsMutext sync.RWMutex

	listeners      []chan *rcon.Event
	listenersMutex sync.RWMutex

	Tomb *tomb.Tomb
}

// NewConnection from the passed in configuration
func NewConnection(ctx context.Context) *Connection {
	c := &Connection{
		Protocol:  be_proto.New(),
		listeners: []chan *rcon.Event{},
	}
	atomic.StoreUint32(&c.seq, 0)
	atomic.StoreInt64(&c.keepAliveCount, 0)
	atomic.StoreInt64(&c.pingbackCount, 0)
	c.transmissions = make(map[be_proto.Sequence]*Transmission)
	c.Tomb, _ = tomb.WithContext(ctx)
	return c
}

//go:generate counterfeiter -o ../../mocks/udp_dialer.go --fake-name UDPDialer . udpDialer
type udpDialer interface {
	DialUDP(string, *net.UDPAddr, *net.UDPAddr) (UDPConnection, error)
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
func (c *Connection) Open(ctx context.Context) error {
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

	err = c.Protocol.VerifyLogin(buf[:n])
	if err != nil {
		return errors.Wrap(err, "login failed")
	}
	c.Hold(ctx)
	return nil
}

// Hold the connection by sending keepalive packets as required by the battleye protocol
func (c *Connection) Hold(ctx context.Context) {
	c.Tomb.Go(c.WriterLoop(ctx))
	c.Tomb.Go(c.ReaderLoop(ctx))
}

// WriterLoop for keeping the connection alive
func (c *Connection) WriterLoop(ctx context.Context) func() error {
	return func() error {
		for {
			select {
			case <-c.Tomb.Dying():
				return tomb.ErrDying
			case <-time.After(time.Second * time.Duration(c.KeepAliveTimeout)):
				if c.UDP != nil {
					c.UDP.Write(c.Protocol.BuildKeepAlivePacket(c.Sequence()))
					c.AddKeepAlive()
					continue
				}
				return errors.New("udp connection must not be nil")
			}
		}
	}
}

// ReaderLoop for keeping the connection alive
func (c *Connection) ReaderLoop(ctx context.Context) func() error {
	return func() error {
		for {
			select {
			case <-c.Tomb.Dying():
				return tomb.ErrDying
			default:
				if c.UDP != nil {
					buf := make([]byte, 4096)
					_, err := c.UDP.Read(buf)
					if err, ok := err.(net.Error); ok && err.Timeout() {
						ctx.Debug("timeout", zap.Error(err))
						continue
					}
					if err != nil {
						return errors.Wrap(err, "reading udp failed")
					}
					go c.HandlePacket(ctx, buf)
				}
				return errors.New("udp connection must not be nil")
			}
		}
	}
}

// Close the connection for graceful shutdown or reconnect
func (c *Connection) Close(ctx context.Context) error {
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
func (c *Connection) Write(ctx context.Context, cmd string) (rcon.Transmission, error) {
	if c.UDP == nil {
		return nil, errors.New("udp connection must not be nil")
	}
	seq := c.AddSequence()
	trm := NewTransmission(cmd)
	c.AddTransmission(seq, trm)
	_, err := c.UDP.Write(c.Protocol.BuildCmdPacket([]byte(trm.Request()), seq))
	if err != nil {
		c.DeleteTransmission(seq)
		return nil, errors.Wrap(err, "writing udp failed")
	}
	return trm, nil
}

// Listen for events on the connection.
func (c *Connection) Listen(ctx context.Context, to chan *rcon.Event) {
	c.listenersMutex.Lock()
	defer c.listenersMutex.Unlock()
	c.listeners = append(c.listeners, to)
}
