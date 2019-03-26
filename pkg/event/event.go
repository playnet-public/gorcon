package event

import "time"

// Event is the generic interface for events handled by the broker
type Event interface {
	Timestamp() time.Time
	Kind() string
	Data() string
}
