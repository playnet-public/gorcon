package rcon

import (
	"errors"
	"sync"
	"time"
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
	// Write a command to the connection
	Write(string) error
	// Listen for events on the connection. This is a blocking call sending on the passed in channel and returning once an error occurs
	Listen(chan<- Event) error
}

// Client is the interface for specific rcon implementations which provides connections or acts as connection pool
//go:generate counterfeiter -o ../mocks/rcon_client.go --fake-name RconClient . Client
type Client interface {
	NewConnection() Connection
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
	}
	con := r.Client.NewConnection()
	if con == nil {
		return errors.New("client returned nil connection")
	}
	r.m.Lock()
	defer r.m.Unlock()
	r.Con = con
	return r.Con.Open()
}
