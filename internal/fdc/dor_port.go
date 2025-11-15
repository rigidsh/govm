package fdc

import "github.com/rigidsh/govm/internal/kvm"

type dorMessage uint8

func (message dorMessage) motor(drive DriveNumber) bool {
	return message>>(4+drive)&1 == 1
}

func (message dorMessage) drive() DriveNumber {
	return DriveNumber(message & 0b11)
}

func (message dorMessage) enabled() bool {
	return message&0b100 != 0
}

func (message dorMessage) dma() bool {
	return message&0b1000 != 0
}

func (fdc *FDC) setupDORPort(port uint16) {
	fdc.vm.RegisterPortHandler(port, kvm.CompositePort(
		kvm.CallbackPort(func(write bool, data []byte) []byte {
			message := dorMessage(data[0])

			fdc.enable(message.enabled())

			if !message.enabled() {
				return nil
			}

			fdc.dma(message.dma())

			//for drive := DriveNumber(0); drive < 4; drive++ {
			//	fdc.drives[drive].motor(message.motor(drive))
			//}
			//
			//fdc.selectDrive(message.drive())

			return nil
		}),
		kvm.NopPort(),
	))
}
