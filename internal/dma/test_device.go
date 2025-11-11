package dma

import "github.com/rigidsh/govm/internal/kvm"

func createReadableTestDevice(vm *kvm.VM, dma *DMA, channelNumber uint8, port uint16, data []byte) {
	readPosition := 0
	completed := false

	dma.ConnectChannel(channelNumber, &ChannelConnector{
		Read: func(buf []byte) uint16 {
			copy(buf, data[readPosition:])
			if len(buf) > len(data)-readPosition {
				return uint16(len(buf))
			}
			return uint16(len(data) - readPosition)
		},
		TCCallback: func() {
			completed = true
		},
	})

	vm.RegisterPortHandler(port, kvm.CompositePort(
		kvm.CallbackPort(func(write bool, data []byte) []byte {
			dma.DREQ(channelNumber, true)
			return nil
		}),
		kvm.CallbackPort(func(write bool, data []byte) []byte {
			result := []byte{0}
			if completed {
				result = []byte{1}
			}
			completed = false

			return result
		}),
	))

}

type testDeviceBuffer struct {
	data []byte
}

func createWritableTestDevice(vm *kvm.VM, dma *DMA, channelNumber uint8, port uint16) *testDeviceBuffer {
	buffer := &testDeviceBuffer{
		data: make([]byte, 0),
	}
	completed := false

	dma.ConnectChannel(channelNumber, &ChannelConnector{
		Write: func(buf []byte) uint16 {
			buffer.data = append(buffer.data, buf...)
			return uint16(len(buf))
		},
		TCCallback: func() {
			completed = true
		},
	})

	vm.RegisterPortHandler(port, kvm.CompositePort(
		kvm.CallbackPort(func(write bool, data []byte) []byte {
			dma.DREQ(channelNumber, true)
			return nil
		}),
		kvm.CallbackPort(func(write bool, data []byte) []byte {
			result := []byte{0}
			if completed {
				result = []byte{1}
			}
			completed = false

			return result
		}),
	))

	return buffer
}
