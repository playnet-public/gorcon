package watcher

import "time"

// Event describes a log event emitted by the process
// TODO(kwiesmueller): rework this to a new Event interface used by all broker dependents
type Event struct {
	timestamp time.Time
	kind      string
	payload   string
}

// Timestamp when the event occurred
func (e Event) Timestamp() time.Time {
	return e.timestamp
}

// Kind of the event
func (e Event) Kind() string {
	return e.kind
}

// Data of the real event
func (e Event) Data() string { return e.payload }

