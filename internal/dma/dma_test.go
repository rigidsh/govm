package dma

import (
	"encoding/binary"
	"testing"

	"github.com/rigidsh/govm/internal/asm"
	"github.com/rigidsh/govm/internal/kvm"
	"github.com/rigidsh/govm/internal/kvm_testing"
)

func TestDMA_WriteReadBaseAddress(t *testing.T) {
	result, err := kvm_testing.Do(
		func(vm *kvm.VM) {
			CreateDMA(vm, MasterPortConfig)
		},
		asm.Compile(
			//Reset DMA flip-flop
			asm.MOV(asm.AL(), asm.Impl8(0x00)),
			asm.OUT(asm.Impl8(uint8(MasterPortConfig.FlipFlopResetPort)), asm.AL()),
			//Write base address for DMA ch0 0x0201
			asm.MOV(asm.AL(), asm.Impl8(0x01)),
			asm.OUT(asm.Impl8(uint8(MasterPortConfig.BaseAddressPort[0])), asm.AL()),
			asm.MOV(asm.AL(), asm.Impl8(0x02)),
			asm.OUT(asm.Impl8(uint8(MasterPortConfig.BaseAddressPort[0])), asm.AL()),
			// Read current address for DMA ch0
			asm.IN(asm.AL(), asm.Impl8(uint8(MasterPortConfig.BaseAddressPort[0]))),
			kvm_testing.WRITE_TEST_RESULT_FROM_AL(),
			asm.IN(asm.AL(), asm.Impl8(uint8(MasterPortConfig.BaseAddressPort[0]))),
			kvm_testing.WRITE_TEST_RESULT_FROM_AL(),

			kvm_testing.STOP_TEST(0),
		),
	)
	if err != nil {
		t.Error(err)
		return
	}

	if binary.LittleEndian.Uint16(result) != 0x0201 {
		t.Errorf("Incorrect ch0 current address, expected %X, get %X", 0x0201, binary.LittleEndian.Uint16(result))
	}
}

func TestDMA_WriteReadBaseCounter(t *testing.T) {
	result, err := kvm_testing.Do(
		func(vm *kvm.VM) {
			CreateDMA(vm, MasterPortConfig)
		},
		asm.Compile(
			//Reset DMA flip-flop
			asm.MOV(asm.AL(), asm.Impl8(0x00)),
			asm.OUT(asm.Impl8(uint8(MasterPortConfig.FlipFlopResetPort)), asm.AL()),
			//Write base counter for DMA ch0 0x0201
			asm.MOV(asm.AL(), asm.Impl8(0x01)),
			asm.OUT(asm.Impl8(uint8(MasterPortConfig.BaseCounterPort[0])), asm.AL()),
			asm.MOV(asm.AL(), asm.Impl8(0x02)),
			asm.OUT(asm.Impl8(uint8(MasterPortConfig.BaseCounterPort[0])), asm.AL()),
			// Read current counter for DMA ch0
			asm.IN(asm.AL(), asm.Impl8(uint8(MasterPortConfig.BaseCounterPort[0]))),
			kvm_testing.WRITE_TEST_RESULT_FROM_AL(),
			asm.IN(asm.AL(), asm.Impl8(uint8(MasterPortConfig.BaseCounterPort[0]))),
			kvm_testing.WRITE_TEST_RESULT_FROM_AL(),

			kvm_testing.STOP_TEST(0),
		),
	)
	if err != nil {
		t.Error(err)
		return
	}

	if binary.LittleEndian.Uint16(result) != 0x0201 {
		t.Errorf("Incorrect ch0 current counter, expected %X, get %X", 0x0201, binary.LittleEndian.Uint16(result))
	}
}

func TestDMA_MMTProcess(t *testing.T) {
	result, err := kvm_testing.Do(
		func(vm *kvm.VM) {
			CreateDMA(vm, MasterPortConfig)
			testData := make([]byte, 0x100)
			for i := 0; i < 0x100; i++ {
				testData[i] = 0xFF
			}
			vm.Memory().Write(0xFE00, testData)
		},
		asm.Compile(
			//Reset DMA flip-flop
			asm.MOV(asm.AL(), asm.Impl8(0x00)),
			asm.OUT(asm.Impl8(uint8(MasterPortConfig.FlipFlopResetPort)), asm.AL()),
			//Write base address for DMA ch0 0xFE00
			asm.MOV(asm.AL(), asm.Impl8(0x00)),
			asm.OUT(asm.Impl8(uint8(MasterPortConfig.BaseAddressPort[0])), asm.AL()),
			asm.MOV(asm.AL(), asm.Impl8(0xFE)),
			asm.OUT(asm.Impl8(uint8(MasterPortConfig.BaseAddressPort[0])), asm.AL()),
			//Write base address for DMA ch1 0xFF00
			asm.MOV(asm.AL(), asm.Impl8(0x00)),
			asm.OUT(asm.Impl8(uint8(MasterPortConfig.BaseAddressPort[1])), asm.AL()),
			asm.MOV(asm.AL(), asm.Impl8(0xFF)),
			asm.OUT(asm.Impl8(uint8(MasterPortConfig.BaseAddressPort[1])), asm.AL()),
			//Write base counter for DMA ch0 0x0100-1=0x00FF
			asm.MOV(asm.AL(), asm.Impl8(0xFF)),
			asm.OUT(asm.Impl8(uint8(MasterPortConfig.BaseCounterPort[0])), asm.AL()),
			asm.MOV(asm.AL(), asm.Impl8(0x00)),
			asm.OUT(asm.Impl8(uint8(MasterPortConfig.BaseCounterPort[0])), asm.AL()),
			// Write 0b00001000 to mode port(Set Write mode(0b10) for ch0)
			asm.MOV(asm.AL(), asm.Impl8(0b00001000)),
			asm.OUT(asm.Impl8(uint8(MasterPortConfig.ModePort)), asm.AL()),
			//Write 0b00000001 to control port(enable MMMT)
			asm.MOV(asm.AL(), asm.Impl8(0b00000001)),
			asm.OUT(asm.Impl8(uint8(MasterPortConfig.CommandPort)), asm.AL()),
			//Write 0b00000001 to request port(dreq for ch0)
			asm.MOV(asm.AL(), asm.Impl8(0b00000100)),
			asm.OUT(asm.Impl8(uint8(MasterPortConfig.RequestPort)), asm.AL()),
			//Read tc for chan 0 in circle
			asm.IN(asm.AL(), asm.Impl8(uint8(MasterPortConfig.StatusPort))),
			asm.CMP(asm.AX(), asm.Impl16(0b00010000)),
			asm.JNE(asm.RelAddr(-7)),
			//Copy result to port
			kvm_testing.WRITE_TEST_RESULT_REGION(0xFF00, 256),

			kvm_testing.STOP_TEST(0),
		),
	)
	if err != nil {
		t.Error(err)
		return
	}

	if len(result) != 256 {
		t.Errorf("Incorrect result size, expected %d, get %d", 256, len(result))
	}

	for i := 0; i < 256; i++ {
		if result[i] != 0xFF {
			t.Errorf("Incorrect copied value for byte %d, expected %X, get %X", i, 0xFF, result[i])
		}
	}
}

func TestDMA_ReadFromDevice(t *testing.T) {
	result, err := kvm_testing.Do(
		func(vm *kvm.VM) {
			dma := CreateDMA(vm, MasterPortConfig)

			testData := make([]byte, 0x100)
			for i := 0; i < 0x100; i++ {
				testData[i] = 0xAA
			}

			createReadableTestDevice(
				vm,
				dma,
				0,
				0xA0,
				testData,
			)
		},
		asm.Compile(
			//Reset DMA flip-flop
			asm.MOV(asm.AL(), asm.Impl8(0x00)),
			asm.OUT(asm.Impl8(uint8(MasterPortConfig.FlipFlopResetPort)), asm.AL()),
			//Write base address for DMA ch0 0xFF00
			asm.MOV(asm.AL(), asm.Impl8(0x00)),
			asm.OUT(asm.Impl8(uint8(MasterPortConfig.BaseAddressPort[0])), asm.AL()),
			asm.MOV(asm.AL(), asm.Impl8(0xFF)),
			asm.OUT(asm.Impl8(uint8(MasterPortConfig.BaseAddressPort[0])), asm.AL()),
			//Write base counter for DMA ch0 0x0100-1=0x00FF
			asm.MOV(asm.AL(), asm.Impl8(0xFF)),
			asm.OUT(asm.Impl8(uint8(MasterPortConfig.BaseCounterPort[0])), asm.AL()),
			asm.MOV(asm.AL(), asm.Impl8(0x00)),
			asm.OUT(asm.Impl8(uint8(MasterPortConfig.BaseCounterPort[0])), asm.AL()),
			// Write 0b00001000 to mode port(Set ReadMode mode(0b11) for ch0)
			asm.MOV(asm.AL(), asm.Impl8(0b00001100)),
			asm.OUT(asm.Impl8(uint8(MasterPortConfig.ModePort)), asm.AL()),
			//Trigger test device
			asm.OUT(asm.Impl8(uint8(0xA0)), asm.AL()),
			//Read tc for chan 0 in circle
			asm.IN(asm.AL(), asm.Impl8(uint8(MasterPortConfig.StatusPort))),
			asm.CMP(asm.AX(), asm.Impl16(0b00010000)),
			asm.JNE(asm.RelAddr(-7)),
			//Copy result to port
			kvm_testing.WRITE_TEST_RESULT_REGION(0xFF00, 256),

			kvm_testing.STOP_TEST(0),
		),
	)
	if err != nil {
		t.Error(err)
		return
	}

	if len(result) != 256 {
		t.Errorf("Incorrect result size, expected %d, get %d", 256, len(result))
	}

	for i := 0; i < 256; i++ {
		if result[i] != 0xAA {
			t.Errorf("Incorrect copied value for byte %d, expected %X, get %X", i, 0xFF, result[i])
		}
	}
}

func TestDMA_WriteToDevice(t *testing.T) {
	var testDevice *testDeviceBuffer
	_, err := kvm_testing.Do(
		func(vm *kvm.VM) {
			dma := CreateDMA(vm, MasterPortConfig)

			testData := make([]byte, 0x100)
			for i := 0; i < 0x100; i++ {
				testData[i] = 0xAA
			}
			vm.Memory().Write(0xFF00, testData)

			testDevice = createWritableTestDevice(
				vm,
				dma,
				0,
				0xA0,
			)
		},
		asm.Compile(
			//Reset DMA flip-flop
			asm.MOV(asm.AL(), asm.Impl8(0x00)),
			asm.OUT(asm.Impl8(uint8(MasterPortConfig.FlipFlopResetPort)), asm.AL()),
			//Write base address for DMA ch0 0xFF00
			asm.MOV(asm.AL(), asm.Impl8(0x00)),
			asm.OUT(asm.Impl8(uint8(MasterPortConfig.BaseAddressPort[0])), asm.AL()),
			asm.MOV(asm.AL(), asm.Impl8(0xFF)),
			asm.OUT(asm.Impl8(uint8(MasterPortConfig.BaseAddressPort[0])), asm.AL()),
			//Write base counter for DMA ch0 0x0100-1=0x00FF
			asm.MOV(asm.AL(), asm.Impl8(0xFF)),
			asm.OUT(asm.Impl8(uint8(MasterPortConfig.BaseCounterPort[0])), asm.AL()),
			asm.MOV(asm.AL(), asm.Impl8(0x00)),
			asm.OUT(asm.Impl8(uint8(MasterPortConfig.BaseCounterPort[0])), asm.AL()),
			// Write 0b00001000 to mode port(Set WriteMode mode(0b10) for ch0)
			asm.MOV(asm.AL(), asm.Impl8(0b00001000)),
			asm.OUT(asm.Impl8(uint8(MasterPortConfig.ModePort)), asm.AL()),
			//Trigger test device
			asm.OUT(asm.Impl8(uint8(0xA0)), asm.AL()),
			//Read tc for chan 0 in circle
			asm.IN(asm.AL(), asm.Impl8(uint8(MasterPortConfig.StatusPort))),
			asm.CMP(asm.AX(), asm.Impl16(0b00010000)),
			asm.JNE(asm.RelAddr(-7)),

			kvm_testing.STOP_TEST(0),
		),
	)
	if err != nil {
		t.Error(err)
		return
	}

	if len(testDevice.data) != 256 {
		t.Errorf("Incorrect result size, expected %d, get %d", 256, len(testDevice.data))
	}

	for i := 0; i < 256; i++ {
		if testDevice.data[i] != 0xAA {
			t.Errorf("Incorrect copied value for byte %d, expected %X, get %X", i, 0xFF, testDevice.data[i])
		}
	}

}
