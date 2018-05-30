package watcher

import "time"

// Event describes a log event emitted by the process
type Event struct {
	Timestamp time.Time
	Type      string
	Payload   string
}
