package dma

import (
	"github.com/rigidsh/govm/internal/kvm"
)

type controlRegister uint8

func (register controlRegister) readBool(index uint8) bool {
	return register.read(index) != 0
}

func (register controlRegister) read(index uint8) uint8 {
	return uint8((register >> index) & 1)
}

func (register controlRegister) mmtEnabled() bool {
	return register.readBool(0)
}

func (register controlRegister) channel0AddressHold() bool {
	return register.readBool(1)
}

func (register controlRegister) enabled() bool {
	return register.readBool(2)
}

func (register controlRegister) timingSelect() uint8 {
	return register.read(3)
}

func (register controlRegister) priority(value uint8) uint8 {
	return register.read(4)
}
func (register controlRegister) extendedWrite() bool {
	return register.readBool(5)
}
func (register controlRegister) dreq() uint8 {
	return register.read(6)
}
func (register controlRegister) dack() uint8 {
	return register.read(7)
}

type statusRegister uint8

func (register *statusRegister) writeBool(index uint8, value bool) {
	intValue := uint8(0)
	if value {
		intValue = 1
	}

	intValue = intValue >> index

	*register = statusRegister(uint8(*register) | intValue)
}

func (register *statusRegister) tcStatus(channel uint8, value bool) {
	register.writeBool(channel, value)
}

func (register *statusRegister) requestStatus(channel uint8, value bool) {
	register.writeBool(channel+4, value)
}

func (dma *DMA) setupCommandPort(commandPort uint16) {

	//TODO: rewrite
	dma.vm.RegisterPortHandler(commandPort,
		kvm.CompositePort(
			kvm.WOPort(
				kvm.CallbackPort(func(write bool, data []byte) []byte {
					controlRegister := controlRegister(data[0])

					dma.updateMmt(controlRegister.mmtEnabled())
					return data
				}),
			),
			kvm.ROPort(
				kvm.CallbackPort(func(write bool, data []byte) []byte {
					result := statusRegister(0)
					for i := uint8(0); i < 4; i++ {
						result.tcStatus(i, dma.channels[i].readAndClearTC())
						//result.requestStatus(i, dma.channels[i].inProgressFlag)
					}

					return []byte{byte(result)}
				}),
			),
		),
	)
}
