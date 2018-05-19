package battleye_test

import (
	"testing"

	"github.com/playnet-public/battleye/battleye"
)

func TestVerify(t *testing.T) {
	p := battleye.New()
	tests := []struct {
		name    string
		packet  battleye.Packet
		wantErr bool
	}{
		{
			"ok",
			[]byte{66, 69, 49, 101, 26, 11, 255, 1, 0, 116, 101, 115, 99, 109, 100},
			false,
		},
		{
			"invalid_header",
			[]byte{66, 69, 0, 0, 0, 0},
			true,
		},
		{
			"invalid_packet",
			[]byte{0, 69, 49, 0, 26, 11, 255, 1, 0, 116, 101, 115, 99, 109, 100},
			true,
		},
		{
			"invalid_checksum",
			[]byte{66, 69, 49, 0, 26, 11, 255, 1, 0, 116, 101, 115, 99, 109, 100},
			true,
		},
		{
			"invalid_seq",
			[]byte{66, 69, 27, 223, 250, 165, 255, 1},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := p.Verify(tt.packet)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyPacket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestSequence(t *testing.T) {
	p := battleye.New()
	tests := []struct {
		name    string
		data    battleye.Packet
		want    battleye.Sequence
		wantErr bool
	}{
		{
			"ok",
			[]byte{'B', 'E', 1, 1, 1, 1, 0, 1, 255, 10, 5, 2, 82},
			255,
			false,
		},
		{
			"error",
			[]byte{'B', 'E', 1, 1, 1, 1},
			0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.Sequence(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSequence() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetSequence() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestType(t *testing.T) {
	p := battleye.New()
	type args struct {
	}
	tests := []struct {
		name    string
		data    battleye.Packet
		want    battleye.Type
		wantErr bool
	}{
		{
			"ok",
			[]byte{'B', 'E', 1, 1, 1, 1, 0, 1, 255, 10, 5, 2, 82},
			1,
			false,
		},
		{
			"error",
			[]byte{'B', 'E', 1, 1, 1, 1, 0},
			0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.Type(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResponseType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ResponseType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestData(t *testing.T) {
	p := battleye.New()
	tests := []struct {
		name     string
		packet   battleye.Packet
		wantData string
		wantErr  bool
	}{
		{
			"ok",
			[]byte{66, 69, 49, 101, 26, 11, 255, 1, 0, 116, 101, 115, 99, 109, 100},
			string([]byte{255, 1, 0, 116, 101, 115, 99, 109, 100}),
			false,
		},
		{
			"invalid_header",
			[]byte{66, 69, 0, 0, 0, 0},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, err := p.Data(tt.packet)
			if (err != nil) != tt.wantErr {
				t.Errorf("Data() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(gotData) != tt.wantData {
				t.Errorf("Data() gotData = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}

func TestVerifyLogin(t *testing.T) {
	p := battleye.New()
	tests := []struct {
		name    string
		packet  battleye.Packet
		wantErr bool
	}{
		{
			"ok",
			[]byte{66, 69, 40, 236, 197, 47, 255, 1, 0x01},
			false,
		},
		{
			"fail",
			[]byte{66, 69, 190, 220, 194, 88, 255, 1, 0x00},
			true,
		},
		{
			"invalid_type",
			battleye.BuildLoginResponse(0x02),
			true,
		},
		{
			"error",
			[]byte{0, 69, 49, 0, 26, 11, 255, 1, 0, 116, 101, 115, 99, 109, 100},
			true,
		},
		{
			"checksum_error",
			[]byte{66, 69, 0, 220, 194, 88, 255, 1, 0x00},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := p.VerifyLogin(tt.packet)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyLogin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMulti(t *testing.T) {
	p := battleye.New()
	tests := []struct {
		name  string
		data  battleye.Packet
		want  byte
		want1 byte
		want2 bool
	}{
		{
			"basic",
			[]byte{1, 0, 116, 101, 115, 99, 109, 100},
			0,
			0,
			true,
		},
		{
			"invalid_data",
			[]byte{1, 0},
			0,
			0,
			true,
		},
		{
			"multi_packet",
			[]byte{1, 1, 0, 4, 3, 99, 109, 100},
			4,
			3,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := p.Multi(tt.data)
			if got != tt.want {
				t.Errorf("CheckMultiPacketResponse() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CheckMultiPacketResponse() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("CheckMultiPacketResponse() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
