package fdc

import (
	"fmt"
	"io"

	"github.com/rigidsh/govm/internal/dma"
	"github.com/rigidsh/govm/internal/kvm"
)

type DriveNumber uint8

type FDC struct {
	vm           *kvm.VM
	dreqLine     *dma.Line
	dmaConnector *dma.ChannelConnector
	drives       [4]*diskDrive

	busy         bool
	currentRead  io.Reader
	currentWrite io.Writer
}

func CreateFDC(vm *kvm.VM, dmaController *dma.DMA) *FDC {
	fdc := &FDC{
		vm: vm,
	}

	for i := 0; i < 4; i++ {
		fdc.drives[i] = &diskDrive{}
	}

	fdc.dmaConnector = &dma.ChannelConnector{
		Write: func(buf []byte) uint16 {
			if fdc.currentWrite != nil {
				result, _ := fdc.currentWrite.Write(buf)
				fmt.Println("Write to FDC")
				return uint16(result)
			}
			return 0
		},
		Read: func(buf []byte) uint16 {
			if fdc.currentRead != nil {
				result, _ := fdc.currentRead.Read(buf)
				fmt.Println("Read from FDC")
				return uint16(result)
			}
			return 0
		},
		TC: dma.NewLine(),
	}

	fdc.dreqLine = dmaController.ConnectChannel(2, fdc.dmaConnector)

	fdc.setupDORPort(0x3F2)
	fdc.setupMSRPort(0x3F4)
	fdc.setupDataPort(0x3F5)

	return fdc
}

func (fdc *FDC) enable(value bool) {

}

func (fdc *FDC) dma(value bool) {

}

func (fdc *FDC) readData(drive, head, cylinder, sector, sectorSize, endOfTrack, gapLength, dataLength uint8, callback func()) {
	fmt.Println("Do read FDC")
	fdc.busy = true

	selectedDrive := fdc.drives[drive-1]
	selectedDrive.seek(cylinder)
	selectedDrive.setSettings(sectorSize, gapLength, dataLength)
	reader, _ := selectedDrive.sectorReader(sector, endOfTrack, head, head)
	//TODO: process error
	fdc.currentRead = reader
	fdc.dreqLine.Set(true)

	go func() {
		<-fdc.dmaConnector.TC.PosEdge()
		fdc.busy = false
		callback()
	}()
	fdc.dreqLine.Set(true)
}

func (fdc *FDC) InsertDisk(drive uint8, disk Disk) {
	fdc.drives[drive].disk = disk
}

func (fdc *FDC) EjectDisk(drive uint8) {
	fdc.drives[drive].disk = nil
}
