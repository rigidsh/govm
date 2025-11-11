package dma

import (
	"github.com/rigidsh/govm/internal/kvm"
)

type requestMessage uint8

func (message requestMessage) channel() uint8 {
	return uint8(message & 0b11)
}

type requestAction uint8

var dmaRequestActionReset requestAction = 0
var dmaRequestActionSet requestAction = 1

func (message requestMessage) action() requestAction {
	return requestAction((message & 0b100) >> 2)
}

func (dma *DMA) setupRequestPort(requestPort uint16) {
	dma.vm.RegisterPortHandler(requestPort,
		kvm.WOPort(
			kvm.CallbackPort(
				func(write bool, data []byte) []byte {
					if write {
						message := requestMessage(data[0])
						dma.DREQ(message.channel(), message.action() == dmaRequestActionSet)

					}
					return nil
				},
			),
		),
	)
}
