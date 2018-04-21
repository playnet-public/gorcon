package battleye

import "github.com/playnet-public/gorcon/pkg/rcon"

// Client is a BattlEye specific implementation of rcon.Client to create new BattlEye rcon connections
type Client struct {
}

// NewConnection from the current client's configuration
func (c *Client) NewConnection() rcon.Connection {
	return &Connection{}
}

// Connection is a BattlEye specific implementation of rcon.Connection offering all required rcon generics
type Connection struct {
}

// Open the connection
func (c *Connection) Open() error {
	return nil
}

// Close the connection for graceful shutdown or reconnect
func (c *Connection) Close() error {
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
