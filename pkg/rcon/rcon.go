// Package rcon offers generics for working with rcon connections. Game specific connections reside in their respective sub-packages
package rcon

import (
	"sync"
	"time"

	"github.com/pkg/errors"
)

// Rcon is the wrapper around the rcon connection interface
type Rcon struct {
	Client Client

	m   sync.Mutex
	Con Connection
}

// Connection is the parent interface for working with different rcon connections
//go:generate counterfeiter -o ../mocks/rcon_connection.go --fake-name RconConnection . Connection
type Connection interface {
	// Open the connection
	Open() error
	// Close the connection for graceful shutdown or reconnect
	Close() error
	// Write a command to the connection and return the resulting transmission
	Write(string) (Transmission, error)
	// Listen for events on the connection. This is a blocking call sending on the passed in channel and returning once an error occurs
	Listen(chan<- Event) error
}

// Client is the interface for specific rcon implementations which provides connections or acts as connection pool
//go:generate counterfeiter -o ../mocks/rcon_client.go --fake-name RconClient . Client
type Client interface {
	NewConnection() Connection
}

// Transmission is the interface describing rcon commands and their respective response.
// It offers the use of specific types internally while still enabling interoperability between the general rcon interface
//go:generate counterfeiter -o ../mocks/rcon_transmission.go --fake-name RconTransmission . Transmission
type Transmission interface {
	Key() uint32
	Request() string
	Done() <-chan bool
	Response() string
}

// Event describes an rcon event happening on the server and being received by the connection
type Event struct {
	Timestamp time.Time
	Message   string
}

// Connect to rcon server
func (r *Rcon) Connect() error {
	if r.Client == nil {
		return errors.New("client must not be nil")
	}
	if r.Con != nil {
		return errors.New("connection already present")
	}
	r.Con = r.Client.NewConnection()
	if r.Con == nil {
		return errors.New("client returned nil connection")
	}
	r.m.Lock()
	defer r.m.Unlock()
	return r.Con.Open()
}

// Reconnect to rcon server. This tries to gracefully close the current connection and then replace it with a new one
// A failing close will not stop the reconnection process for now
func (r *Rcon) Reconnect() error {
	if r.Client == nil {
		return errors.New("client must not be nil")
	}
	r.Con.Close()
	r.Con = r.Client.NewConnection()
	if r.Con == nil {
		return errors.New("client returned nil connection")
	}
	r.m.Lock()
	defer r.m.Unlock()
	return r.Con.Open()
}

// Disconnect from rcon. This tries to gracefully close the current connection and resets the local Connection internally
// A failing close will result in an error
func (r *Rcon) Disconnect() error {
	if r.Con == nil {
		return errors.New("connection already nil")
	}
	err := r.Con.Close()
	if err != nil {
		return errors.Wrap(err, "failed to close current connection")
	}
	r.m.Lock()
	defer r.m.Unlock()
	r.Con = nil
	return nil
}
