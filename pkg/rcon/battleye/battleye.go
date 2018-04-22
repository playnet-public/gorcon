package battleye

import (
	"net"
	"time"

	"github.com/pkg/errors"

	"github.com/playnet-public/gorcon/pkg/rcon"
)

// Client is a BattlEye specific implementation of rcon.Client to create new BattlEye rcon connections
type Client struct {
}

// NewConnection from the current client's configuration
func (c *Client) NewConnection() rcon.Connection {
	return &Connection{}
}

// Connection is a BattlEye specific implementation of rcon.Connection offering all required rcon generics
type Connection struct {
	Addr     *net.UDPAddr
	Password string
	Dialer   udpDialer

	UDP UDPConnection
}

//go:generate counterfeiter -o ../../mocks/udp_dialer.go --fake-name UDPDialer . udpDialer
type udpDialer interface {
	DialUDP(string, *net.UDPAddr, *net.UDPAddr) (UDPConnection, error)
}

// UDPConnection interface defines all udp functions required and is used primarily for mocking
//go:generate counterfeiter -o ../../mocks/udp_connection.go --fake-name UDPConnection . UDPConnection
type UDPConnection interface {
	Close() error

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
	return nil
}

// Close the connection for graceful shutdown or reconnect
func (c *Connection) Close() error {
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
