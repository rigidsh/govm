package fdc

import (
	"testing"

	"github.com/rigidsh/govm/internal/asm"
	"github.com/rigidsh/govm/internal/dma"
	"github.com/rigidsh/govm/internal/kvm"
	"github.com/rigidsh/govm/internal/kvm_testing"
)

func readFromDisk(disk, head, cylinder, sector, sectorSize, endOfTrack, gapLength, dataLength uint8) asm.Command {
	return asm.Compile(
		asm.MOV(asm.DX(), asm.Impl16(0x3F5)),
		//Read command
		asm.MOV(asm.AL(), asm.Impl8(0x06)),
		asm.OUT(asm.DX(), asm.AL()),
		//Select Disk and Head
		asm.MOV(asm.AL(), asm.Impl8(disk+(head)<<2)),
		asm.OUT(asm.DX(), asm.AL()),
		//Set Cylinder
		asm.MOV(asm.AL(), asm.Impl8(cylinder)),
		asm.OUT(asm.DX(), asm.AL()),
		//Set Head
		asm.MOV(asm.AL(), asm.Impl8(head)),
		asm.OUT(asm.DX(), asm.AL()),
		//Set Sector
		asm.MOV(asm.AL(), asm.Impl8(sector)),
		asm.OUT(asm.DX(), asm.AL()),
		//Set SectorSize
		asm.MOV(asm.AL(), asm.Impl8(sectorSize)),
		asm.OUT(asm.DX(), asm.AL()),
		//Set EndOfTrack
		asm.MOV(asm.AL(), asm.Impl8(endOfTrack)),
		asm.OUT(asm.DX(), asm.AL()),
		//Set Gap Length
		asm.MOV(asm.AL(), asm.Impl8(gapLength)),
		asm.OUT(asm.DX(), asm.AL()),
		//Set Data Length
		asm.MOV(asm.AL(), asm.Impl8(dataLength)),
		asm.OUT(asm.DX(), asm.AL()),
	)
}

func waitFDCNotBusy() asm.Command {
	return asm.Compile(
		//Read MSR status in circle
		asm.MOV(asm.DX(), asm.Impl16(0x3F4)),
		asm.IN(asm.AL(), asm.DX()),
		asm.CMP(asm.AX(), asm.Impl16(0)),
		asm.JNE(asm.RelAddr(-6)))
}

func TestFDC_ReadData(t *testing.T) {
	result, err := kvm_testing.Do(
		func(vm *kvm.VM) {
			dmaController := dma.CreateDMA(vm, dma.MasterPortConfig)
			fdc := CreateFDC(vm, dmaController)
			disk := InMemoryRawDisk(make([]byte, 512))
			for i := 0; i < 512; i++ {
				disk[i] = 0xFF
			}
			fdc.InsertDisk(0, disk)
		},
		asm.Compile(
			//Reset DMA flip-flop
			asm.MOV(asm.AL(), asm.Impl8(0x00)),
			asm.OUT(asm.Impl8(uint8(dma.MasterPortConfig.FlipFlopResetPort)), asm.AL()),
			//Write base address for DMA ch2 0xF000
			asm.MOV(asm.AL(), asm.Impl8(0x00)),
			asm.OUT(asm.Impl8(uint8(dma.MasterPortConfig.BaseAddressPort[2])), asm.AL()),
			asm.MOV(asm.AL(), asm.Impl8(0xF0)),
			asm.OUT(asm.Impl8(uint8(dma.MasterPortConfig.BaseAddressPort[2])), asm.AL()),
			//Write base counter for DMA ch2 0x200-1=0x01FF
			asm.MOV(asm.AL(), asm.Impl8(0xFF)),
			asm.OUT(asm.Impl8(uint8(dma.MasterPortConfig.BaseCounterPort[2])), asm.AL()),
			asm.MOV(asm.AL(), asm.Impl8(0x01)),
			asm.OUT(asm.Impl8(uint8(dma.MasterPortConfig.BaseCounterPort[2])), asm.AL()),
			// Write 0b00001110 to mode port(Set ReadMode mode(0b11) for ch2)
			asm.MOV(asm.AL(), asm.Impl8(0b00001110)),
			asm.OUT(asm.Impl8(uint8(dma.MasterPortConfig.ModePort)), asm.AL()),
			//Read from FDC, Disk 1, Head 0, Cylinder 0, Sector 1
			readFromDisk(1, 0, 0, 1, 0x02, 18, 0x1B, 0xFF),
			waitFDCNotBusy(),

			kvm_testing.WRITE_TEST_RESULT_REGION(0xF000, 512),
			kvm_testing.STOP_TEST(0),
		),
	)
	if err != nil {
		t.Error(err)
		return
	}

	if len(result) != 512 {
		t.Errorf("Incorrect result size. Expected: %d, Actual: %d", 512, len(result))
	}

	for i := 0; i < 512; i++ {
		if result[i] != 0xFF {
			t.Errorf("Incorrect read value for byte %d, expected %X, get %X", i, 0xFF, result[i])
		}
	}
}
