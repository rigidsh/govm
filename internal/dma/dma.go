package dma

import (
	"fmt"

	"github.com/rigidsh/govm/internal/kvm"
)

type transferType uint8

var verifyTransferType transferType = 0b01
var writeTransferType transferType = 0b10
var readTransferType transferType = 0b11

type ChannelConnector interface {
	Read(buf []byte) uint16
	Write(buf []byte) uint16
}

func CreateDMA(vm *kvm.VM, config PortConfig) *DMA {
	dma := &DMA{
		vm:       vm,
		channels: [4]*channel{},
		flipFlop: kvm.NewFlipFlop(),
	}

	for channelNumber := 0; channelNumber < 4; channelNumber++ {
		channel := newChannel(
			dma,
			config.BaseAddressPort[channelNumber],
			config.PagePort[channelNumber],
			config.BaseCounterPort[channelNumber],
		)
		dma.channels[channelNumber] = channel
	}

	dma.setupCommandPort(config.CommandPort)
	dma.setupModePort(config.ModePort)
	dma.setupRequestPort(config.RequestPort)
	dma.setupMaskPort(config.SingleMaskPort)
	dma.setupStatusPort(config.StatusPort)

	vm.RegisterPortHandler(config.FlipFlopResetPort,
		kvm.WOPort(
			kvm.CallbackPort(
				func(write bool, _ []byte) []byte {
					if write {
						fmt.Println("Reset FlipFlop")
						dma.flipFlop.Reset()
					}
					return nil
				},
			),
		),
	)

	dma.run()

	return dma
}

type DMA struct {
	vm       *kvm.VM
	flipFlop *kvm.FlipFlop
	channels [4]*channel
	stopChan chan interface{}
	mmt      bool
}

func (dma *DMA) run() {
	go func() {
		for dma.waitDREQ() {

			for i := 0; i < 4; i++ {
				if dma.channels[i].dreq.value {
					dma.channels[i].doIteration()
				}
			}
		}
		fmt.Println("DMA stop :(")
	}()
}

func (dma *DMA) waitDREQ() bool {
	if dma.channels[0].dreq.value ||
		dma.channels[1].dreq.value ||
		dma.channels[2].dreq.value ||
		dma.channels[3].dreq.value {
		return true
	}

	select {
	case <-dma.channels[0].dreq.PosEdge():
	case <-dma.channels[1].dreq.PosEdge():
	case <-dma.channels[2].dreq.PosEdge():
	case <-dma.channels[3].dreq.PosEdge():
	case <-dma.stopChan:
		return false
	}

	return true
}

func (dma *DMA) updateMmt(value bool) {
	if value {
		dma.ConnectChannel(0, &mmtBuffConnector{dma: dma}, nil)
	}
	dma.mmt = value
}

func (dma *DMA) reset() {
	dma.flipFlop.Reset()
	for i := 0; i < 4; i++ {
		dma.channels[i].mask()
	}
}

func (dma *DMA) DREQ(channel uint8) Line {
	return dma.channels[channel].dreq
}

func (dma *DMA) ConnectChannel(channelNumber uint8, connector ChannelConnector, tc *ObservableLine) {
	dma.channels[channelNumber].connector = connector
	dma.channels[channelNumber].tcLine = tc
}
