package dma

import "github.com/rigidsh/govm/internal/kvm"

func (dma *DMA) setupModePort(modePort uint16) {
	dma.vm.RegisterPortHandler(modePort,
		kvm.CallbackPort(
			func(write bool, data []byte) []byte {
				if write {
					message := modeMessage(data[0])
					dma.channels[message.channel()].transferType = message.transferType()
					dma.channels[message.channel()].autoInit = message.autoInit()
				}

				return nil
			},
		),
	)
}

type modeMessage uint8

func (message modeMessage) channel() uint8 {
	return uint8(message & 0b11)
}

func (message modeMessage) transferType() transferType {
	return transferType((message >> 2) & 0b11)
}

func (message modeMessage) autoInit() bool {
	return (message>>4)&1 == 1
}

func (message modeMessage) decreaseMode() bool {
	return (message>>5)&1 == 1
}
