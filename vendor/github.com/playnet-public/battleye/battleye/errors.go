package battleye

import "errors"

var (
	//ErrUnknownPacketType .
	ErrUnknownPacketType = errors.New("received unknown packet type")
	//ErrInvalidLoginPacket .
	ErrInvalidLoginPacket = errors.New("received invalid login packet")
	//ErrInvalidLogin .
	ErrInvalidLogin = errors.New("server refused login")
	//ErrInvalidChecksum .
	ErrInvalidChecksum = errors.New("received invalid packet checksum")
	//ErrInvalidSizeNoHeader .
	ErrInvalidSizeNoHeader = errors.New("invalid packet size, no header found")
	//ErrInvalidSizeNoSequence .
	ErrInvalidSizeNoSequence = errors.New("invalid packet size, no sequence found")
	//ErrInvalidHeaderSize .
	ErrInvalidHeaderSize = errors.New("invalid packet header size")
	//ErrInvalidHeaderSyntax .
	ErrInvalidHeaderSyntax = errors.New("invalid packet header syntax")
	//ErrInvalidHeaderEnd .
	ErrInvalidHeaderEnd = errors.New("invalid packet header end")
	//ErrInvalidSize .
	ErrInvalidSize = errors.New("packet size too small")
	//ErrUnknownEventType .
	ErrUnknownEventType = errors.New("unknown event type")
	//ErrUnableToParse .
	ErrUnableToParse = errors.New("unable to parse")
)
