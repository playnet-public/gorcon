package battleye_test

import (
	"reflect"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/playnet-public/battleye/battleye"
)

func TestBattleye(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Battleye Suite")
}

var _ = Describe("Battleye", func() {
	Describe("New", func() {
		It("does not return nil", func() {
			Expect(battleye.New()).NotTo(BeNil())
		})
	})
})

// Old Legacy Tests go here
func TestBuildPackets(t *testing.T) {
	p := battleye.New()
	tests := []struct {
		name string
		seq  battleye.Sequence
		data []byte
		t    battleye.Type
		want battleye.Packet
	}{
		{
			"login",
			battleye.Sequence(0),
			[]byte("pw"),
			battleye.Login,
			battleye.Packet([]byte{66, 69, 132, 68, 31, 30, 255, 0, 112, 119}),
		},
		{
			"keep_alive",
			battleye.Sequence(0),
			[]byte{},
			battleye.Command,
			battleye.Packet([]byte{66, 69, 190, 220, 194, 88, 255, 1, 0}),
		},
		{
			"cmd",
			battleye.Sequence(0),
			[]byte("xxx"),
			battleye.Command,
			battleye.Packet([]byte{66, 69, 199, 16, 188, 139, 255, 1, 0, 120, 120, 120}),
		},
		{
			"ack",
			battleye.Sequence(0),
			[]byte{},
			battleye.ServerMessage,
			battleye.Packet([]byte{66, 69, 125, 143, 239, 115, 255, 2, 0}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.t {
			case battleye.Login:
				if got := p.BuildLoginPacket(string(tt.data)); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("BuildLoginPacket() = %v, want %v", got, tt.want)
				}
			case battleye.Command:
				if len(tt.data) < 1 {
					if got := p.BuildKeepAlivePacket(tt.seq); !reflect.DeepEqual(got, tt.want) {
						t.Errorf("BuildKeepAlivePacket() = %v, want %v", got, tt.want)
					}
				} else {
					if got := p.BuildCmdPacket(tt.data, tt.seq); !reflect.DeepEqual(got, tt.want) {
						t.Errorf("BuildCmdPacket() = %v, want %v", got, tt.want)
					}
				}
			case battleye.ServerMessage:
				if got := p.BuildMsgAckPacket(tt.seq); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("BuildMsgAckPacket() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
