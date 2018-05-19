package battleye

import (
	"fmt"

	"github.com/pkg/errors"
)

// Packet is the real data sent between server and client
type Packet []byte

// Verify if the packet is valid
func (p *protocol) Verify(d Packet) error {
	_, err := stripHeader(d)
	if err != nil {
		return err
	}
	checksum, err := getChecksum(d)
	if err != nil {
		return err
	}
	match := verifyChecksum(d[6:], checksum)
	if !match {
		err = ErrInvalidChecksum
		return err
	}
	_, err = p.Sequence(d)
	if err != nil {
		return err
	}
	_, err = p.Type(d)
	return nil
}

// Sequence extracts the seq number from a packet
func (p *protocol) Sequence(d Packet) (Sequence, error) {
	if len(d) < 9 {
		return 0, ErrInvalidSizeNoSequence
	}
	return Sequence(d[8]), nil
}

// Type determines the kind of response from a packet
func (p *protocol) Type(d Packet) (Type, error) {
	if len(d) < 8 {
		return 0, ErrInvalidSize
	}
	return Type(d[7]), nil
}

// Data returns the actual data inside the packet
func (p *protocol) Data(d Packet) ([]byte, error) {
	data, err := stripHeader(d)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// VerifyLogin returns nil on successful login
// and a respective error on failed login
func (p *protocol) VerifyLogin(d Packet) error {
	if len(d) != 9 {
		fmt.Println("Call:", d)
		return ErrInvalidLoginPacket
	}
	if match, err := verifyChecksumMatch(d); match == false || err != nil {
		return err
	}
	switch Type(d[8]) {
	case LoginOk:
		return nil
	case LoginFail:
		return ErrInvalidLogin
	}
	return errors.Wrap(ErrInvalidLoginPacket, "triggered default error")
}

// Multi checks whether a packet is part of a multiPacketResponse
// Returns: packetCount, currentPacket and isSingle
func (p *protocol) Multi(d Packet) (byte, byte, bool) {
	if len(d) < 3 {
		return 0, 0, true
	}
	if d[0] != 0x01 || d[2] != 0x00 {
		return 0, 0, true
	}
	return d[3], d[4], false
}
