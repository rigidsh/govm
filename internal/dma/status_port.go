package dma

import (
	"fmt"

	"github.com/rigidsh/govm/internal/kvm"
)

func (dma *DMA) setupStatusPort(maskPort uint16) {
	dma.vm.RegisterPortHandler(maskPort,
		kvm.CallbackPort(
			func(write bool, data []byte) []byte {
				if !write {
					fmt.Println("Status port")
					result := statusMessage(0)

					for i := uint8(0); i < 4; i++ {
						result.setTCValue(i, dma.channels[i].readAndClearTC())
						result.setDREQValie(i, dma.channels[i].dreq.value)
					}
					return []byte{uint8(result)}
				}
				return nil
			},
		),
	)
}

type statusMessage uint8

func (message *statusMessage) setTCValue(channel uint8, value bool) {
	intValue := uint8(0)
	if value {
		intValue = 1
	}
	*message = statusMessage(uint8(*message) | (intValue << (4 + channel)))
}

func (message *statusMessage) setDREQValie(channel uint8, value bool) {
	intValue := uint8(0)
	if value {
		intValue = 1
	}
	*message = statusMessage(uint8(*message) | (intValue << channel))
}
