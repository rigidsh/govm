package dma

import "github.com/rigidsh/govm/internal/kvm"

type maskStatusMessage uint8

func (message *maskStatusMessage) mask(channelNumber uint8, value bool) {
	intValue := uint8(0)
	if value {
		intValue = 1
	}

	intValue = intValue << channelNumber

	*message = maskStatusMessage(uint8(*message) | intValue)
}

func (dma *DMA) setupResetPort(resetPort uint16) {
	dma.vm.RegisterPortHandler(resetPort,
		kvm.CompositePort(
			kvm.CallbackPort(func(write bool, data []byte) []byte {
				dma.reset()
				return nil
			}),
			kvm.CallbackPort(func(write bool, data []byte) []byte {
				result := maskStatusMessage(0)
				for i := uint8(0); i < 4; i++ {
					result.mask(i, dma.channels[i].enabled)
				}

				return []byte{byte(result)}
			}),
		),
	)
}
