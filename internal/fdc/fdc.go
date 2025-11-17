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
	dreqLine     dma.Line
	tcLine       *dma.ObservableLine
	dmaConnector *dma.BuffReadConnector
	drives       [4]*diskDrive
	dataPort     *kvm.CommandPort

	busy bool
}

func CreateFDC(vm *kvm.VM, dmaController *dma.DMA) *FDC {
	fdc := &FDC{
		vm:     vm,
		tcLine: dma.NewLine(),
	}

	for i := 0; i < 4; i++ {
		fdc.drives[i] = &diskDrive{}
	}

	fdc.dmaConnector = dma.NewBuffReadConnector(dmaController.DREQ(2), 512)
	dmaController.ConnectChannel(2, fdc.dmaConnector, fdc.tcLine)
	fdc.dreqLine = dmaController.DREQ(2)

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

	go func() {
		selectedDrive.setSettings(sectorSize, gapLength, dataLength)
		selectedDrive.seek(cylinder)
		selectedDrive.head(head)
		for sector != endOfTrack {
			if fdc.tcLine.Get() {
				fdc.busy = false
				return
			}

			err := fdc.readSector(selectedDrive, sector)
			if err != nil {
				return
			}
			sector = sector + 1
		}
	}()
}

func (fdc *FDC) readSector(drive *diskDrive, sector uint8) error {
	reader, _ := drive.sectorReader(sector)
	for {
		if fdc.tcLine.Get() {
			fdc.busy = false
			return nil
		}
		err := fdc.dmaConnector.ReadFrom(reader)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

func (fdc *FDC) InsertDisk(drive uint8, disk Disk) {
	fdc.drives[drive].disk = disk
}

func (fdc *FDC) EjectDisk(drive uint8) {
	fdc.drives[drive].disk = nil
}
