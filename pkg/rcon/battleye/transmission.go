package battleye

// Transmission is the BattlEye implementation of rcon.Transmission
type Transmission struct {
	seq      uint32
	request  []byte
	done     bool
	response []byte
}

// NewTransmission containing request
func NewTransmission(request string) *Transmission {
	return &Transmission{
		request: []byte(request),
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

// Done returns true if the transmission was finished receiving
func (t *Transmission) Done() bool {
	return t.done
}

// Response returns the final response
// Checking if the transmission is done before retrieving is suggested
func (t *Transmission) Response() string {
	return string(t.response)
}