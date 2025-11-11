package bios

import (
	"github.com/rigidsh/govm/internal/asm"
	"github.com/rigidsh/govm/internal/kvm"
)

var BIOS_ROM = asm.Compile(
	asm.MOV(asm.AH(), asm.Impl8(0x0E)),
	asm.MOV(asm.AL(), asm.Impl8('L')), asm.INT(asm.Impl8(0x10)),
	asm.MOV(asm.AL(), asm.Impl8('o')), asm.INT(asm.Impl8(0x10)),
	asm.MOV(asm.AL(), asm.Impl8('a')), asm.INT(asm.Impl8(0x10)),
	asm.MOV(asm.AL(), asm.Impl8('d')), asm.INT(asm.Impl8(0x10)),
	asm.MOV(asm.AL(), asm.Impl8('i')), asm.INT(asm.Impl8(0x10)),
	asm.MOV(asm.AL(), asm.Impl8('n')), asm.INT(asm.Impl8(0x10)),
	asm.MOV(asm.AL(), asm.Impl8('g')), asm.INT(asm.Impl8(0x10)),
	asm.MOV(asm.AL(), asm.Impl8('\n')), asm.INT(asm.Impl8(0x10)),

	asm.XOR(asm.AX(), asm.AX()),
	asm.MOV(asm.ES(), asm.AX()),
	asm.MOV(asm.BX(), asm.Impl16(0x7C00)),
	asm.MOV(asm.AH(), asm.Impl8(0x02)),
	asm.MOV(asm.CH(), asm.Impl8(0x00)),
	asm.MOV(asm.CL(), asm.Impl8(0x01)),
	asm.MOV(asm.DH(), asm.Impl8(0x00)),
	asm.MOV(asm.AL(), asm.Impl8(0x01)),
	asm.INT(asm.Impl8(0x13)),

	asm.MOV(asm.AX(), asm.Impl16(0x7DFE)),

	asm.CMP(asm.AX(), asm.Impl16(0xAA55)),
	asm.JNE(asm.RelAddr(0x05)),
	asm.JMP(asm.SegAddr(0x0000, 0x7C00)),
	asm.JMP(asm.SegAddr(0xFFF0, 0xF000)),
)

var BIOS_INT_HANDLERS_ROM = asm.Compile(
	//int 10h handler
	asm.OUT(asm.Impl8(0xFA), asm.AL()),
	asm.IRET(),
	//int 13h handler
	asm.OUT(asm.Impl8(0xFB), asm.AL()),
	asm.IRET(),
	//int 15h handler
	asm.OUT(asm.Impl8(0xFC), asm.AL()),
	asm.IRET(),
	//int 16h handler
	asm.OUT(asm.Impl8(0xFc), asm.AL()),
	asm.IRET(),
)

func CreateBIOS(vm *kvm.VM) *BIOS {
	return &BIOS{
		vm: vm,
	}
}

type BIOS struct {
	vm *kvm.VM
}

func (bios *BIOS) Init(kvmManager *kvm.KVM) error {
	rom, err := bios.rom(kvmManager)
	if err != nil {
		return err
	}

	err = bios.vm.Memory().AddRAMRegion(0xF0000, rom)

	if err != nil {
		return err
	}

	return nil
}

func (bios *BIOS) fillBDA() {
	bdaBaseAddress := uint64(0x0400)
	bios.vm.Memory().Write16(bdaBaseAddress+0x10, 0b0000000000100101)
	bios.vm.Memory().Write16(bdaBaseAddress+0x13, 640)
	// keyboard bufer
	bios.vm.Memory().Write16(bdaBaseAddress+0x1A, 0x0041E)
	bios.vm.Memory().Write16(bdaBaseAddress+0x1C, 0x0041E)
	bios.vm.Memory().Write8(bdaBaseAddress+0x49, 0x03)
	bios.vm.Memory().Write8(bdaBaseAddress+0x4A, 80)
	bios.vm.Memory().Write8(bdaBaseAddress+0x4C, 24)
	bios.vm.Memory().Write16(bdaBaseAddress+0x63, 0x3D4)
	//bios.vm.Memory().Write16(bdaBaseAddress+0x2A, 0x)
}

func (bios *BIOS) rom(kvmManager *kvm.KVM) (kvm.RAM, error) {
	rom, err := kvmManager.AllocateRAM(0x10000)
	if err != nil {
		return nil, err
	}

	copy(rom, BIOS_INT_HANDLERS_ROM)
	copy(rom[0x0100:], BIOS_ROM)

	//copy(rom[0xFFF0:], a.Compile(a.INT(a.Impl8(0x16))))
	//copy(rom[0xFFF0:], a.Compile(
	//	a.HLT(),
	//	a.JMP(a.SegAddr(0xF000, 0x0100))))
	copy(rom[0xFFF0:], asm.Compile(asm.JMP(asm.SegAddr(0xF000, 0x0100))))

	return rom, nil
}
