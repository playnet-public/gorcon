package protocol

import (
	"io"
	"time"
)

// Transmission containing a full bercon conversation
type Transmission struct {
	Packet      []byte
	Command     []byte
	Sequence    uint32
	Response    []byte
	Timestamp   time.Time
	WriteCloser io.WriteCloser
}

//PacketType contains all possible types a packet could have
var PacketType = struct {
	Login         byte
	Command       byte
	MultiCommand  byte
	ServerMessage byte
}{
	Login:         0x00,
	Command:       0x01,
	MultiCommand:  0x00,
	ServerMessage: 0x02,
}

//PacketResponse contains all types a server could respond with
var PacketResponse = struct {
	LoginOk      byte
	LoginFail    byte
	MultiCommand byte
}{
	LoginOk:      0x01,
	LoginFail:    0x00,
	MultiCommand: 0x00,
}

//ResponseType gets the PacketResponse from a packet
func ResponseType(data []byte) (byte, error) {
	if len(data) < 8 {
		return 0, ErrInvalidSize
	}
	return data[7], nil
}

//GetSequence extracts the seq number from a packet
func GetSequence(data []byte) (uint32, error) {
	if len(data) < 9 {
		return 0, ErrInvalidSizeNoSequence
	}
	return uint32(data[8]), nil
}
