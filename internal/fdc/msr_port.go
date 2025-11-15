package fdc

import "github.com/rigidsh/govm/internal/kvm"

type msrMessage byte

func (message *msrMessage) busy(value bool) {
	intValue := byte(0)
	if value {
		intValue = 1
	}

	*message = msrMessage(intValue << 5)
}

func (fdc *FDC) setupMSRPort(port uint16) {
	fdc.vm.RegisterPortHandler(port, kvm.CompositePort(
		kvm.NopPort(),
		kvm.CallbackPort(func(write bool, data []byte) []byte {
			message := msrMessage(0)

			message.busy(fdc.busy)

			return []byte{byte(message)}
		}),
	))
}
