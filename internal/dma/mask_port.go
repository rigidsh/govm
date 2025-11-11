package dma

import "github.com/rigidsh/govm/internal/kvm"

func (dma *DMA) setupMaskPort(maskPort uint16) {
	dma.vm.RegisterPortHandler(maskPort,
		kvm.CallbackPort(
			func(write bool, data []byte) []byte {
				if write {
					message := maskMessage(data[0])
					if message.mask() {
						dma.channels[message.channel()].mask()
					} else {
						dma.channels[message.channel()].unmask()
					}
				}
				return nil
			},
		),
	)
}

type maskMessage uint8

func (message maskMessage) channel() uint8 {
	return uint8(message & 0b11)
}

func (message maskMessage) mask() bool {
	return (message>>2)&1 == 1
}
