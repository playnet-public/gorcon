package battleye

import (
	"sync/atomic"
)

// Sequence gets the current sequence using atomic
func (c *Connection) Sequence() uint32 {
	return atomic.LoadUint32(&c.seq)
}

// AddSequence increments the sequence
func (c *Connection) AddSequence() {
	atomic.AddUint32(&c.seq, 1)
}

// Pingback gets the current pingbackCount using atomic
func (c *Connection) Pingback() int64 {
	return atomic.LoadInt64(&c.pingbackCount)
}

// AddPingback increments the pingbackCount
func (c *Connection) AddPingback() {
	atomic.AddInt64(&c.pingbackCount, 1)
}

// ResetPingback increments the pingbackCount
func (c *Connection) ResetPingback() {
	atomic.SwapInt64(&c.pingbackCount, 0)
}

// KeepAlive gets the current keepAliveCount using atomic
func (c *Connection) KeepAlive() int64 {
	return atomic.LoadInt64(&c.keepAliveCount)
}

// AddKeepAlive increments the keepAliveCount
func (c *Connection) AddKeepAlive() {
	atomic.AddInt64(&c.keepAliveCount, 1)
}

// ResetKeepAlive increments the keepAliveCount
func (c *Connection) ResetKeepAlive() {
	atomic.SwapInt64(&c.keepAliveCount, 0)
}
