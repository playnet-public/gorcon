package battleye

import (
	"encoding/binary"
)

func buildHeader(checksum uint32) []byte {
	check := make([]byte, 4)
	binary.LittleEndian.PutUint32(check, checksum)
	return append([]byte{}, 'B', 'E', check[0], check[1], check[2], check[3])
}

func stripHeader(data []byte) ([]byte, error) {
	if len(data) < 7 {
		return []byte{}, ErrInvalidSizeNoHeader
	}
	return data[6:], nil
}

// BuildLoginResponse for testing
func BuildLoginResponse(t Type) Packet {
	data := []byte{0xFF, byte(t), byte(t)}
	return Packet(append(buildHeader(makeChecksum(data)), data...))
}
