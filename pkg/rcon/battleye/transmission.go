package battleye

import (
	"sort"
)

// Transmission is the BattlEye implementation of rcon.Transmission
type Transmission struct {
	seq      uint32
	request  []byte
	done     chan bool
	response []byte

	// As we might receive multiple packets responding to a single transmission
	// we have to collect them by their respective id and return them after a final sort
	// Thanks UDP.
	multiBuffer map[int][]byte
}

// NewTransmission containing request
func NewTransmission(request string) *Transmission {
	return &Transmission{
		request:     []byte(request),
		done:        make(chan bool),
		multiBuffer: make(map[int][]byte),
	}
}

// Key retrieves the transmissions sequence for identifying and retrieving it further on in the process
func (t *Transmission) Key() uint32 {
	return t.seq
}

// Request retrieves a string representation of the command to send
func (t *Transmission) Request() string {
	return string(t.request)
}

// Done returns blocking channel indicating transmission status
func (t *Transmission) Done() <-chan bool {
	return t.done
}

// Response returns the final response
// Checking if the transmission is done before retrieving is suggested
// Otherwise this might render the transaction useless caused by the way multiResponsePackets get handled
// TODO: If this might as well be done for every call
func (t *Transmission) Response() string {
	if len(t.response) < 1 {
		// Sort the responses stored in buffer
		var keys []int
		for k := range t.multiBuffer {
			keys = append(keys, k)
		}
		sort.Ints(keys)

		// Build final response
		for _, k := range keys {
			t.response = append(t.response, t.multiBuffer[k]...)
		}
	}
	return string(t.response)
}
