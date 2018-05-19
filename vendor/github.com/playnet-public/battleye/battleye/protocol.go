package battleye

// Protocol offers an interface representation for BattlEye communications
//go:generate counterfeiter -o ../mocks/protocol.go --fake-name Protocol . Protocol
type Protocol interface {
	BuildPacket([]byte, Type) Packet
	BuildLoginPacket(string) Packet
	BuildCmdPacket([]byte, Sequence) Packet
	BuildKeepAlivePacket(Sequence) Packet
	BuildMsgAckPacket(Sequence) Packet

	Verify(Packet) error
	Sequence(Packet) (Sequence, error)
	Type(Packet) (Type, error)
	Data(Packet) ([]byte, error)
	VerifyLogin(d Packet) error
	Multi(Packet) (byte, byte, bool)
}

type protocol struct{}

// New provides a real and tested implementation of Protocol
func New() Protocol {
	return &protocol{}
}

// Sequence is used by BattlEye to keep track of transactions
type Sequence uint32

//BuildPacket creates a new packet with data and type
func (p *protocol) BuildPacket(data []byte, t Type) Packet {
	data = append([]byte{0xFF, byte(t)}, data...)
	checksum := makeChecksum(data)
	header := buildHeader(checksum)

	return append(header, data...)
}

//BuildLoginPacket creates a login packet with password
func (p *protocol) BuildLoginPacket(pw string) Packet {
	return p.BuildPacket([]byte(pw), Login)
}

//BuildCmdPacket creates a packet with cmd and seq
func (p *protocol) BuildCmdPacket(cmd []byte, seq Sequence) Packet {
	return p.BuildPacket(append([]byte{byte(seq)}, cmd...), Command)
}

//BuildKeepAlivePacket creates a keepAlivePacket with seq
func (p *protocol) BuildKeepAlivePacket(seq Sequence) Packet {
	return p.BuildPacket([]byte{byte(seq)}, Command)
}

//BuildMsgAckPacket creates a server message packet with seq
func (p *protocol) BuildMsgAckPacket(seq Sequence) Packet {
	return p.BuildPacket([]byte{byte(seq)}, ServerMessage)
}
