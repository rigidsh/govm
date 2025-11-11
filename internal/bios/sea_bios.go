package bios

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"io"

	"github.com/rigidsh/govm/internal/kvm"
)

//go:embed bios.gz
var seaBiosROMCompressed []byte

func CreateSeaBIOS(vm *kvm.VM) *SeaBIOS {
	return &SeaBIOS{
		vm: vm,
	}
}

type SeaBIOS struct {
	vm *kvm.VM
}

func (bios *SeaBIOS) Init(kvmManager *kvm.KVM) error {

	biosReader, err := gzip.NewReader(bytes.NewReader(seaBiosROMCompressed))
	if err != nil {
		return err
	}

	biosImage, err := io.ReadAll(biosReader)
	if err != nil {
		return err
	}

	rom, err := kvmManager.AllocateRAM(256 * 1024)
	if err != nil {
		return err
	}

	copy(rom, biosImage)

	err = bios.vm.Memory().AddRAMRegion(0x100000-256*1024, rom)

	if err != nil {
		return err
	}

	return nil
}
