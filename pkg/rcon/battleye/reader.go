package battleye

import (
	be "github.com/playnet-public/battleye/battleye"
)

// HandlePacket received from UDP connection
func (c *Connection) HandlePacket(data be.Packet) {
	c.Protocol.Verify(data)
}
