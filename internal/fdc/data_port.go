package fdc

import "github.com/rigidsh/govm/internal/kvm"

const readDataCommand byte = 0x06
const writeDataCommand byte = 0x45

type driveSelectMessage byte

func (message driveSelectMessage) drive() DriveNumber { return DriveNumber(message & 0b11) }
func (message driveSelectMessage) head() uint8        { return uint8((message >> 3) & 1) }

type readDataMessage [8]byte

func (message readDataMessage) selector() driveSelectMessage { return driveSelectMessage(message[0]) }
func (message readDataMessage) cylinder() uint8              { return message[1] }
func (message readDataMessage) head() uint8                  { return message[2] }
func (message readDataMessage) sector() uint8                { return message[3] }
func (message readDataMessage) sectorSize() uint8            { return message[4] }
func (message readDataMessage) endOfTrack() uint8            { return message[5] }
func (message readDataMessage) gapLength() uint8             { return message[6] }
func (message readDataMessage) dataLength() uint8            { return message[7] }

func (fdc *FDC) setupDataPort(port uint16) {

	dataPort := kvm.NewCommandPort()

	dataPort.RegisterCommand(readDataCommand, kvm.NewCommandDefinition(8,
		func(argument []byte, resultCallback func([]byte)) {
			message := readDataMessage(argument)
			fdc.readData(
				uint8(message.selector().drive()),
				message.head(),
				message.cylinder(),
				message.sector(),
				message.sectorSize(),
				message.endOfTrack(),
				message.gapLength(),
				message.dataLength(),
				func() {
					resultCallback(make([]byte, 7))
				},
			)
		}),
	)
	dataPort.RegisterCommand(writeDataCommand, kvm.NewCommandDefinition(8,
		func(argument []byte, resultCallback func([]byte)) {

		}),
	)

	fdc.vm.RegisterPortHandler(port, dataPort)
}
