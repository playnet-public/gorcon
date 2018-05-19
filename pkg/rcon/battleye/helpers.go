package battleye

import (
	"sync/atomic"

	be "github.com/playnet-public/battleye/battleye"
)

// Sequence gets the current sequence using atomic
func (c *Connection) Sequence() be.Sequence {
	return be.Sequence(atomic.LoadUint32(&c.seq))
}

// AddSequence increments the sequence
func (c *Connection) AddSequence() be.Sequence {
	return be.Sequence(atomic.AddUint32(&c.seq, 1))
}

// ResetSequence to zero
func (c *Connection) ResetSequence() {
	atomic.SwapUint32(&c.seq, 0)
}

// Pingback gets the current pingbackCount using atomic
func (c *Connection) Pingback() int64 {
	return atomic.LoadInt64(&c.pingbackCount)
}

// AddPingback increments the pingbackCount
func (c *Connection) AddPingback() int64 {
	return atomic.AddInt64(&c.pingbackCount, 1)
}

// ResetPingback to zero
func (c *Connection) ResetPingback() {
	atomic.SwapInt64(&c.pingbackCount, 0)
}

// KeepAlive gets the current keepAliveCount using atomic
func (c *Connection) KeepAlive() int64 {
	return atomic.LoadInt64(&c.keepAliveCount)
}

// AddKeepAlive increments the keepAliveCount
func (c *Connection) AddKeepAlive() int64 {
	return atomic.AddInt64(&c.keepAliveCount, 1)
}

// ResetKeepAlive to zero
func (c *Connection) ResetKeepAlive() {
	atomic.SwapInt64(&c.keepAliveCount, 0)
}

// AddTransmission for sequence to the connection
func (c *Connection) AddTransmission(seq be.Sequence, t *Transmission) {
	c.transmissionsMutext.Lock()
	defer c.transmissionsMutext.Unlock()
	c.transmissions[seq] = t
}

// GetTransmission for sequence from the connection
func (c *Connection) GetTransmission(seq be.Sequence) *Transmission {
	c.transmissionsMutext.RLock()
	defer c.transmissionsMutext.RUnlock()
	return c.transmissions[seq]
}

// DeleteTransmission for sequence from the connection
func (c *Connection) DeleteTransmission(seq be.Sequence) {
	c.transmissionsMutext.RLock()
	defer c.transmissionsMutext.RUnlock()
	delete(c.transmissions, seq)
}
