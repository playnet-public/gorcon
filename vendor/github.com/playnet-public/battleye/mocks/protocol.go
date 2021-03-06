// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"sync"

	"github.com/playnet-public/battleye/battleye"
)

type Protocol struct {
	BuildPacketStub        func([]byte, battleye.Type) battleye.Packet
	buildPacketMutex       sync.RWMutex
	buildPacketArgsForCall []struct {
		arg1 []byte
		arg2 battleye.Type
	}
	buildPacketReturns struct {
		result1 battleye.Packet
	}
	buildPacketReturnsOnCall map[int]struct {
		result1 battleye.Packet
	}
	BuildLoginPacketStub        func(string) battleye.Packet
	buildLoginPacketMutex       sync.RWMutex
	buildLoginPacketArgsForCall []struct {
		arg1 string
	}
	buildLoginPacketReturns struct {
		result1 battleye.Packet
	}
	buildLoginPacketReturnsOnCall map[int]struct {
		result1 battleye.Packet
	}
	BuildCmdPacketStub        func([]byte, battleye.Sequence) battleye.Packet
	buildCmdPacketMutex       sync.RWMutex
	buildCmdPacketArgsForCall []struct {
		arg1 []byte
		arg2 battleye.Sequence
	}
	buildCmdPacketReturns struct {
		result1 battleye.Packet
	}
	buildCmdPacketReturnsOnCall map[int]struct {
		result1 battleye.Packet
	}
	BuildKeepAlivePacketStub        func(battleye.Sequence) battleye.Packet
	buildKeepAlivePacketMutex       sync.RWMutex
	buildKeepAlivePacketArgsForCall []struct {
		arg1 battleye.Sequence
	}
	buildKeepAlivePacketReturns struct {
		result1 battleye.Packet
	}
	buildKeepAlivePacketReturnsOnCall map[int]struct {
		result1 battleye.Packet
	}
	BuildMsgAckPacketStub        func(battleye.Sequence) battleye.Packet
	buildMsgAckPacketMutex       sync.RWMutex
	buildMsgAckPacketArgsForCall []struct {
		arg1 battleye.Sequence
	}
	buildMsgAckPacketReturns struct {
		result1 battleye.Packet
	}
	buildMsgAckPacketReturnsOnCall map[int]struct {
		result1 battleye.Packet
	}
	VerifyStub        func(battleye.Packet) error
	verifyMutex       sync.RWMutex
	verifyArgsForCall []struct {
		arg1 battleye.Packet
	}
	verifyReturns struct {
		result1 error
	}
	verifyReturnsOnCall map[int]struct {
		result1 error
	}
	SequenceStub        func(battleye.Packet) (battleye.Sequence, error)
	sequenceMutex       sync.RWMutex
	sequenceArgsForCall []struct {
		arg1 battleye.Packet
	}
	sequenceReturns struct {
		result1 battleye.Sequence
		result2 error
	}
	sequenceReturnsOnCall map[int]struct {
		result1 battleye.Sequence
		result2 error
	}
	TypeStub        func(battleye.Packet) (battleye.Type, error)
	typeMutex       sync.RWMutex
	typeArgsForCall []struct {
		arg1 battleye.Packet
	}
	typeReturns struct {
		result1 battleye.Type
		result2 error
	}
	typeReturnsOnCall map[int]struct {
		result1 battleye.Type
		result2 error
	}
	DataStub        func(battleye.Packet) ([]byte, error)
	dataMutex       sync.RWMutex
	dataArgsForCall []struct {
		arg1 battleye.Packet
	}
	dataReturns struct {
		result1 []byte
		result2 error
	}
	dataReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	VerifyLoginStub        func(d battleye.Packet) error
	verifyLoginMutex       sync.RWMutex
	verifyLoginArgsForCall []struct {
		d battleye.Packet
	}
	verifyLoginReturns struct {
		result1 error
	}
	verifyLoginReturnsOnCall map[int]struct {
		result1 error
	}
	MultiStub        func(battleye.Packet) (byte, byte, bool)
	multiMutex       sync.RWMutex
	multiArgsForCall []struct {
		arg1 battleye.Packet
	}
	multiReturns struct {
		result1 byte
		result2 byte
		result3 bool
	}
	multiReturnsOnCall map[int]struct {
		result1 byte
		result2 byte
		result3 bool
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Protocol) BuildPacket(arg1 []byte, arg2 battleye.Type) battleye.Packet {
	var arg1Copy []byte
	if arg1 != nil {
		arg1Copy = make([]byte, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.buildPacketMutex.Lock()
	ret, specificReturn := fake.buildPacketReturnsOnCall[len(fake.buildPacketArgsForCall)]
	fake.buildPacketArgsForCall = append(fake.buildPacketArgsForCall, struct {
		arg1 []byte
		arg2 battleye.Type
	}{arg1Copy, arg2})
	fake.recordInvocation("BuildPacket", []interface{}{arg1Copy, arg2})
	fake.buildPacketMutex.Unlock()
	if fake.BuildPacketStub != nil {
		return fake.BuildPacketStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.buildPacketReturns.result1
}

func (fake *Protocol) BuildPacketCallCount() int {
	fake.buildPacketMutex.RLock()
	defer fake.buildPacketMutex.RUnlock()
	return len(fake.buildPacketArgsForCall)
}

func (fake *Protocol) BuildPacketArgsForCall(i int) ([]byte, battleye.Type) {
	fake.buildPacketMutex.RLock()
	defer fake.buildPacketMutex.RUnlock()
	return fake.buildPacketArgsForCall[i].arg1, fake.buildPacketArgsForCall[i].arg2
}

func (fake *Protocol) BuildPacketReturns(result1 battleye.Packet) {
	fake.BuildPacketStub = nil
	fake.buildPacketReturns = struct {
		result1 battleye.Packet
	}{result1}
}

func (fake *Protocol) BuildPacketReturnsOnCall(i int, result1 battleye.Packet) {
	fake.BuildPacketStub = nil
	if fake.buildPacketReturnsOnCall == nil {
		fake.buildPacketReturnsOnCall = make(map[int]struct {
			result1 battleye.Packet
		})
	}
	fake.buildPacketReturnsOnCall[i] = struct {
		result1 battleye.Packet
	}{result1}
}

func (fake *Protocol) BuildLoginPacket(arg1 string) battleye.Packet {
	fake.buildLoginPacketMutex.Lock()
	ret, specificReturn := fake.buildLoginPacketReturnsOnCall[len(fake.buildLoginPacketArgsForCall)]
	fake.buildLoginPacketArgsForCall = append(fake.buildLoginPacketArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("BuildLoginPacket", []interface{}{arg1})
	fake.buildLoginPacketMutex.Unlock()
	if fake.BuildLoginPacketStub != nil {
		return fake.BuildLoginPacketStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.buildLoginPacketReturns.result1
}

func (fake *Protocol) BuildLoginPacketCallCount() int {
	fake.buildLoginPacketMutex.RLock()
	defer fake.buildLoginPacketMutex.RUnlock()
	return len(fake.buildLoginPacketArgsForCall)
}

func (fake *Protocol) BuildLoginPacketArgsForCall(i int) string {
	fake.buildLoginPacketMutex.RLock()
	defer fake.buildLoginPacketMutex.RUnlock()
	return fake.buildLoginPacketArgsForCall[i].arg1
}

func (fake *Protocol) BuildLoginPacketReturns(result1 battleye.Packet) {
	fake.BuildLoginPacketStub = nil
	fake.buildLoginPacketReturns = struct {
		result1 battleye.Packet
	}{result1}
}

func (fake *Protocol) BuildLoginPacketReturnsOnCall(i int, result1 battleye.Packet) {
	fake.BuildLoginPacketStub = nil
	if fake.buildLoginPacketReturnsOnCall == nil {
		fake.buildLoginPacketReturnsOnCall = make(map[int]struct {
			result1 battleye.Packet
		})
	}
	fake.buildLoginPacketReturnsOnCall[i] = struct {
		result1 battleye.Packet
	}{result1}
}

func (fake *Protocol) BuildCmdPacket(arg1 []byte, arg2 battleye.Sequence) battleye.Packet {
	var arg1Copy []byte
	if arg1 != nil {
		arg1Copy = make([]byte, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.buildCmdPacketMutex.Lock()
	ret, specificReturn := fake.buildCmdPacketReturnsOnCall[len(fake.buildCmdPacketArgsForCall)]
	fake.buildCmdPacketArgsForCall = append(fake.buildCmdPacketArgsForCall, struct {
		arg1 []byte
		arg2 battleye.Sequence
	}{arg1Copy, arg2})
	fake.recordInvocation("BuildCmdPacket", []interface{}{arg1Copy, arg2})
	fake.buildCmdPacketMutex.Unlock()
	if fake.BuildCmdPacketStub != nil {
		return fake.BuildCmdPacketStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.buildCmdPacketReturns.result1
}

func (fake *Protocol) BuildCmdPacketCallCount() int {
	fake.buildCmdPacketMutex.RLock()
	defer fake.buildCmdPacketMutex.RUnlock()
	return len(fake.buildCmdPacketArgsForCall)
}

func (fake *Protocol) BuildCmdPacketArgsForCall(i int) ([]byte, battleye.Sequence) {
	fake.buildCmdPacketMutex.RLock()
	defer fake.buildCmdPacketMutex.RUnlock()
	return fake.buildCmdPacketArgsForCall[i].arg1, fake.buildCmdPacketArgsForCall[i].arg2
}

func (fake *Protocol) BuildCmdPacketReturns(result1 battleye.Packet) {
	fake.BuildCmdPacketStub = nil
	fake.buildCmdPacketReturns = struct {
		result1 battleye.Packet
	}{result1}
}

func (fake *Protocol) BuildCmdPacketReturnsOnCall(i int, result1 battleye.Packet) {
	fake.BuildCmdPacketStub = nil
	if fake.buildCmdPacketReturnsOnCall == nil {
		fake.buildCmdPacketReturnsOnCall = make(map[int]struct {
			result1 battleye.Packet
		})
	}
	fake.buildCmdPacketReturnsOnCall[i] = struct {
		result1 battleye.Packet
	}{result1}
}

func (fake *Protocol) BuildKeepAlivePacket(arg1 battleye.Sequence) battleye.Packet {
	fake.buildKeepAlivePacketMutex.Lock()
	ret, specificReturn := fake.buildKeepAlivePacketReturnsOnCall[len(fake.buildKeepAlivePacketArgsForCall)]
	fake.buildKeepAlivePacketArgsForCall = append(fake.buildKeepAlivePacketArgsForCall, struct {
		arg1 battleye.Sequence
	}{arg1})
	fake.recordInvocation("BuildKeepAlivePacket", []interface{}{arg1})
	fake.buildKeepAlivePacketMutex.Unlock()
	if fake.BuildKeepAlivePacketStub != nil {
		return fake.BuildKeepAlivePacketStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.buildKeepAlivePacketReturns.result1
}

func (fake *Protocol) BuildKeepAlivePacketCallCount() int {
	fake.buildKeepAlivePacketMutex.RLock()
	defer fake.buildKeepAlivePacketMutex.RUnlock()
	return len(fake.buildKeepAlivePacketArgsForCall)
}

func (fake *Protocol) BuildKeepAlivePacketArgsForCall(i int) battleye.Sequence {
	fake.buildKeepAlivePacketMutex.RLock()
	defer fake.buildKeepAlivePacketMutex.RUnlock()
	return fake.buildKeepAlivePacketArgsForCall[i].arg1
}

func (fake *Protocol) BuildKeepAlivePacketReturns(result1 battleye.Packet) {
	fake.BuildKeepAlivePacketStub = nil
	fake.buildKeepAlivePacketReturns = struct {
		result1 battleye.Packet
	}{result1}
}

func (fake *Protocol) BuildKeepAlivePacketReturnsOnCall(i int, result1 battleye.Packet) {
	fake.BuildKeepAlivePacketStub = nil
	if fake.buildKeepAlivePacketReturnsOnCall == nil {
		fake.buildKeepAlivePacketReturnsOnCall = make(map[int]struct {
			result1 battleye.Packet
		})
	}
	fake.buildKeepAlivePacketReturnsOnCall[i] = struct {
		result1 battleye.Packet
	}{result1}
}

func (fake *Protocol) BuildMsgAckPacket(arg1 battleye.Sequence) battleye.Packet {
	fake.buildMsgAckPacketMutex.Lock()
	ret, specificReturn := fake.buildMsgAckPacketReturnsOnCall[len(fake.buildMsgAckPacketArgsForCall)]
	fake.buildMsgAckPacketArgsForCall = append(fake.buildMsgAckPacketArgsForCall, struct {
		arg1 battleye.Sequence
	}{arg1})
	fake.recordInvocation("BuildMsgAckPacket", []interface{}{arg1})
	fake.buildMsgAckPacketMutex.Unlock()
	if fake.BuildMsgAckPacketStub != nil {
		return fake.BuildMsgAckPacketStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.buildMsgAckPacketReturns.result1
}

func (fake *Protocol) BuildMsgAckPacketCallCount() int {
	fake.buildMsgAckPacketMutex.RLock()
	defer fake.buildMsgAckPacketMutex.RUnlock()
	return len(fake.buildMsgAckPacketArgsForCall)
}

func (fake *Protocol) BuildMsgAckPacketArgsForCall(i int) battleye.Sequence {
	fake.buildMsgAckPacketMutex.RLock()
	defer fake.buildMsgAckPacketMutex.RUnlock()
	return fake.buildMsgAckPacketArgsForCall[i].arg1
}

func (fake *Protocol) BuildMsgAckPacketReturns(result1 battleye.Packet) {
	fake.BuildMsgAckPacketStub = nil
	fake.buildMsgAckPacketReturns = struct {
		result1 battleye.Packet
	}{result1}
}

func (fake *Protocol) BuildMsgAckPacketReturnsOnCall(i int, result1 battleye.Packet) {
	fake.BuildMsgAckPacketStub = nil
	if fake.buildMsgAckPacketReturnsOnCall == nil {
		fake.buildMsgAckPacketReturnsOnCall = make(map[int]struct {
			result1 battleye.Packet
		})
	}
	fake.buildMsgAckPacketReturnsOnCall[i] = struct {
		result1 battleye.Packet
	}{result1}
}

func (fake *Protocol) Verify(arg1 battleye.Packet) error {
	fake.verifyMutex.Lock()
	ret, specificReturn := fake.verifyReturnsOnCall[len(fake.verifyArgsForCall)]
	fake.verifyArgsForCall = append(fake.verifyArgsForCall, struct {
		arg1 battleye.Packet
	}{arg1})
	fake.recordInvocation("Verify", []interface{}{arg1})
	fake.verifyMutex.Unlock()
	if fake.VerifyStub != nil {
		return fake.VerifyStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.verifyReturns.result1
}

func (fake *Protocol) VerifyCallCount() int {
	fake.verifyMutex.RLock()
	defer fake.verifyMutex.RUnlock()
	return len(fake.verifyArgsForCall)
}

func (fake *Protocol) VerifyArgsForCall(i int) battleye.Packet {
	fake.verifyMutex.RLock()
	defer fake.verifyMutex.RUnlock()
	return fake.verifyArgsForCall[i].arg1
}

func (fake *Protocol) VerifyReturns(result1 error) {
	fake.VerifyStub = nil
	fake.verifyReturns = struct {
		result1 error
	}{result1}
}

func (fake *Protocol) VerifyReturnsOnCall(i int, result1 error) {
	fake.VerifyStub = nil
	if fake.verifyReturnsOnCall == nil {
		fake.verifyReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.verifyReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *Protocol) Sequence(arg1 battleye.Packet) (battleye.Sequence, error) {
	fake.sequenceMutex.Lock()
	ret, specificReturn := fake.sequenceReturnsOnCall[len(fake.sequenceArgsForCall)]
	fake.sequenceArgsForCall = append(fake.sequenceArgsForCall, struct {
		arg1 battleye.Packet
	}{arg1})
	fake.recordInvocation("Sequence", []interface{}{arg1})
	fake.sequenceMutex.Unlock()
	if fake.SequenceStub != nil {
		return fake.SequenceStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.sequenceReturns.result1, fake.sequenceReturns.result2
}

func (fake *Protocol) SequenceCallCount() int {
	fake.sequenceMutex.RLock()
	defer fake.sequenceMutex.RUnlock()
	return len(fake.sequenceArgsForCall)
}

func (fake *Protocol) SequenceArgsForCall(i int) battleye.Packet {
	fake.sequenceMutex.RLock()
	defer fake.sequenceMutex.RUnlock()
	return fake.sequenceArgsForCall[i].arg1
}

func (fake *Protocol) SequenceReturns(result1 battleye.Sequence, result2 error) {
	fake.SequenceStub = nil
	fake.sequenceReturns = struct {
		result1 battleye.Sequence
		result2 error
	}{result1, result2}
}

func (fake *Protocol) SequenceReturnsOnCall(i int, result1 battleye.Sequence, result2 error) {
	fake.SequenceStub = nil
	if fake.sequenceReturnsOnCall == nil {
		fake.sequenceReturnsOnCall = make(map[int]struct {
			result1 battleye.Sequence
			result2 error
		})
	}
	fake.sequenceReturnsOnCall[i] = struct {
		result1 battleye.Sequence
		result2 error
	}{result1, result2}
}

func (fake *Protocol) Type(arg1 battleye.Packet) (battleye.Type, error) {
	fake.typeMutex.Lock()
	ret, specificReturn := fake.typeReturnsOnCall[len(fake.typeArgsForCall)]
	fake.typeArgsForCall = append(fake.typeArgsForCall, struct {
		arg1 battleye.Packet
	}{arg1})
	fake.recordInvocation("Type", []interface{}{arg1})
	fake.typeMutex.Unlock()
	if fake.TypeStub != nil {
		return fake.TypeStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.typeReturns.result1, fake.typeReturns.result2
}

func (fake *Protocol) TypeCallCount() int {
	fake.typeMutex.RLock()
	defer fake.typeMutex.RUnlock()
	return len(fake.typeArgsForCall)
}

func (fake *Protocol) TypeArgsForCall(i int) battleye.Packet {
	fake.typeMutex.RLock()
	defer fake.typeMutex.RUnlock()
	return fake.typeArgsForCall[i].arg1
}

func (fake *Protocol) TypeReturns(result1 battleye.Type, result2 error) {
	fake.TypeStub = nil
	fake.typeReturns = struct {
		result1 battleye.Type
		result2 error
	}{result1, result2}
}

func (fake *Protocol) TypeReturnsOnCall(i int, result1 battleye.Type, result2 error) {
	fake.TypeStub = nil
	if fake.typeReturnsOnCall == nil {
		fake.typeReturnsOnCall = make(map[int]struct {
			result1 battleye.Type
			result2 error
		})
	}
	fake.typeReturnsOnCall[i] = struct {
		result1 battleye.Type
		result2 error
	}{result1, result2}
}

func (fake *Protocol) Data(arg1 battleye.Packet) ([]byte, error) {
	fake.dataMutex.Lock()
	ret, specificReturn := fake.dataReturnsOnCall[len(fake.dataArgsForCall)]
	fake.dataArgsForCall = append(fake.dataArgsForCall, struct {
		arg1 battleye.Packet
	}{arg1})
	fake.recordInvocation("Data", []interface{}{arg1})
	fake.dataMutex.Unlock()
	if fake.DataStub != nil {
		return fake.DataStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.dataReturns.result1, fake.dataReturns.result2
}

func (fake *Protocol) DataCallCount() int {
	fake.dataMutex.RLock()
	defer fake.dataMutex.RUnlock()
	return len(fake.dataArgsForCall)
}

func (fake *Protocol) DataArgsForCall(i int) battleye.Packet {
	fake.dataMutex.RLock()
	defer fake.dataMutex.RUnlock()
	return fake.dataArgsForCall[i].arg1
}

func (fake *Protocol) DataReturns(result1 []byte, result2 error) {
	fake.DataStub = nil
	fake.dataReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *Protocol) DataReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.DataStub = nil
	if fake.dataReturnsOnCall == nil {
		fake.dataReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.dataReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *Protocol) VerifyLogin(d battleye.Packet) error {
	fake.verifyLoginMutex.Lock()
	ret, specificReturn := fake.verifyLoginReturnsOnCall[len(fake.verifyLoginArgsForCall)]
	fake.verifyLoginArgsForCall = append(fake.verifyLoginArgsForCall, struct {
		d battleye.Packet
	}{d})
	fake.recordInvocation("VerifyLogin", []interface{}{d})
	fake.verifyLoginMutex.Unlock()
	if fake.VerifyLoginStub != nil {
		return fake.VerifyLoginStub(d)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.verifyLoginReturns.result1
}

func (fake *Protocol) VerifyLoginCallCount() int {
	fake.verifyLoginMutex.RLock()
	defer fake.verifyLoginMutex.RUnlock()
	return len(fake.verifyLoginArgsForCall)
}

func (fake *Protocol) VerifyLoginArgsForCall(i int) battleye.Packet {
	fake.verifyLoginMutex.RLock()
	defer fake.verifyLoginMutex.RUnlock()
	return fake.verifyLoginArgsForCall[i].d
}

func (fake *Protocol) VerifyLoginReturns(result1 error) {
	fake.VerifyLoginStub = nil
	fake.verifyLoginReturns = struct {
		result1 error
	}{result1}
}

func (fake *Protocol) VerifyLoginReturnsOnCall(i int, result1 error) {
	fake.VerifyLoginStub = nil
	if fake.verifyLoginReturnsOnCall == nil {
		fake.verifyLoginReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.verifyLoginReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *Protocol) Multi(arg1 battleye.Packet) (byte, byte, bool) {
	fake.multiMutex.Lock()
	ret, specificReturn := fake.multiReturnsOnCall[len(fake.multiArgsForCall)]
	fake.multiArgsForCall = append(fake.multiArgsForCall, struct {
		arg1 battleye.Packet
	}{arg1})
	fake.recordInvocation("Multi", []interface{}{arg1})
	fake.multiMutex.Unlock()
	if fake.MultiStub != nil {
		return fake.MultiStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fake.multiReturns.result1, fake.multiReturns.result2, fake.multiReturns.result3
}

func (fake *Protocol) MultiCallCount() int {
	fake.multiMutex.RLock()
	defer fake.multiMutex.RUnlock()
	return len(fake.multiArgsForCall)
}

func (fake *Protocol) MultiArgsForCall(i int) battleye.Packet {
	fake.multiMutex.RLock()
	defer fake.multiMutex.RUnlock()
	return fake.multiArgsForCall[i].arg1
}

func (fake *Protocol) MultiReturns(result1 byte, result2 byte, result3 bool) {
	fake.MultiStub = nil
	fake.multiReturns = struct {
		result1 byte
		result2 byte
		result3 bool
	}{result1, result2, result3}
}

func (fake *Protocol) MultiReturnsOnCall(i int, result1 byte, result2 byte, result3 bool) {
	fake.MultiStub = nil
	if fake.multiReturnsOnCall == nil {
		fake.multiReturnsOnCall = make(map[int]struct {
			result1 byte
			result2 byte
			result3 bool
		})
	}
	fake.multiReturnsOnCall[i] = struct {
		result1 byte
		result2 byte
		result3 bool
	}{result1, result2, result3}
}

func (fake *Protocol) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.buildPacketMutex.RLock()
	defer fake.buildPacketMutex.RUnlock()
	fake.buildLoginPacketMutex.RLock()
	defer fake.buildLoginPacketMutex.RUnlock()
	fake.buildCmdPacketMutex.RLock()
	defer fake.buildCmdPacketMutex.RUnlock()
	fake.buildKeepAlivePacketMutex.RLock()
	defer fake.buildKeepAlivePacketMutex.RUnlock()
	fake.buildMsgAckPacketMutex.RLock()
	defer fake.buildMsgAckPacketMutex.RUnlock()
	fake.verifyMutex.RLock()
	defer fake.verifyMutex.RUnlock()
	fake.sequenceMutex.RLock()
	defer fake.sequenceMutex.RUnlock()
	fake.typeMutex.RLock()
	defer fake.typeMutex.RUnlock()
	fake.dataMutex.RLock()
	defer fake.dataMutex.RUnlock()
	fake.verifyLoginMutex.RLock()
	defer fake.verifyLoginMutex.RUnlock()
	fake.multiMutex.RLock()
	defer fake.multiMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Protocol) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ battleye.Protocol = new(Protocol)
