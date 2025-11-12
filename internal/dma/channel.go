package dma

import (
	"encoding/binary"
	"fmt"

	"github.com/rigidsh/govm/internal/kvm"
)

type channel struct {
	dma *DMA

	dreq *Line

	direction uint8

	connector *ChannelConnector

	transferType transferType
	autoInit     bool

	enabled bool

	tcFlag bool

	baseAddress uint16
	baseCounter uint16

	currentAddress uint16
	currentCounter uint16

	page uint8
}

func newChannel(dma *DMA, baseAddressPort, pagePort, baseCounterPort uint16) *channel {
	channel := &channel{
		enabled: false,
		dma:     dma,
		dreq:    NewLine(),
	}

	dma.vm.RegisterPortHandler(baseAddressPort,
		kvm.CompositePort(
			kvm.OneByteOperationsPort(
				kvm.FlipFlopPort(
					kvm.CallbackPort(
						func(write bool, data []byte) []byte {
							channel.setBaseAddress(binary.LittleEndian.Uint16(data))
							return data
						},
					),
					dma.flipFlop,
				),
			),
			kvm.OneByteOperationsPort(
				kvm.FlipFlopPort(
					kvm.CreateRegister16Port(&channel.currentAddress),
					dma.flipFlop,
				),
			),
		),
	)

	dma.vm.RegisterPortHandler(pagePort,
		kvm.WOPort(
			kvm.CreateRegister8Port(&channel.page),
		),
	)

	dma.vm.RegisterPortHandler(baseCounterPort,
		kvm.CompositePort(
			kvm.OneByteOperationsPort(
				kvm.FlipFlopPort(
					kvm.CallbackPort(
						func(write bool, data []byte) []byte {
							channel.setBaseCounter(binary.LittleEndian.Uint16(data))
							return data
						},
					),
					dma.flipFlop,
				),
			),
			kvm.OneByteOperationsPort(
				kvm.FlipFlopPort(
					kvm.CreateRegister16Port(&channel.currentCounter),
					dma.flipFlop,
				),
			),
		),
	)

	return channel
}

func (channel *channel) doIteration() {
	if channel.transferType == readTransferType {
		requiredToRead := channel.currentCounter + 1
		buf := make([]byte, requiredToRead)
		readFromSource := channel.connector.Read(buf)
		channel.dma.vm.Memory().Write(channel.phyAddress(), buf[:readFromSource])
		channel.currentAddress = channel.currentAddress + readFromSource
		if requiredToRead-readFromSource == 0 {
			channel.currentCounter = 0

			channel.tc()

			return
		}

		channel.currentCounter = requiredToRead - readFromSource - 1
	} else if channel.transferType == writeTransferType {
		requiredToWrite := channel.currentCounter + 1
		buf := make([]byte, requiredToWrite)
		channel.dma.vm.Memory().Read(channel.phyAddress(), buf)
		wroteFromSource := channel.connector.Write(buf)
		channel.currentAddress = channel.currentAddress + wroteFromSource
		if requiredToWrite-wroteFromSource == 0 {
			channel.currentCounter = 0
			channel.tc()

			return
		}

		channel.currentCounter = requiredToWrite - wroteFromSource - 1
	}
}

func (channel *channel) phyAddress() uint64 {
	return (uint64(channel.page) << 16) + uint64(channel.currentAddress)
}

func (channel *channel) mask() {
	if channel.enabled {
		channel.enabled = false
	}
}

func (channel *channel) unmask() {
	channel.enabled = true
}

func (channel *channel) setBaseAddress(value uint16) {
	channel.baseAddress = value
	channel.currentAddress = value
}

func (channel *channel) setBaseCounter(value uint16) {
	channel.baseCounter = value
	channel.currentCounter = value
}

func (channel *channel) init() {
	channel.currentCounter = channel.baseCounter
	channel.currentAddress = channel.baseAddress
}

func (channel *channel) readAndClearTC() bool {
	result := channel.tcFlag

	channel.tcFlag = false

	return result
}

func (channel *channel) tc() {
	channel.tcFlag = true
	channel.dreq.Set(false)
	if channel.connector.TC != nil {
		channel.connector.TC.Set(true)
	}

	if channel.autoInit {
		channel.init()
	}
	fmt.Println("Done")
}
