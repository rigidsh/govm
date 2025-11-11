package kvm

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type RAM []byte

func (ram RAM) CopyTo(baseAddress int, data []byte) {
	copy(ram[baseAddress:len(data)+baseAddress], data)
}

type MMIORegion struct {
	Size         uint64
	WriteHandler func(relAddress uint64, data []byte)
	ReadHandler  func(relAddress uint64, size int8) []byte
}

func newMemory(vm *VM) *Memory {
	return &Memory{
		vm:          vm,
		ramRegions:  make(map[uint64]RAM),
		mmioRegions: make(map[uint64]MMIORegion),
	}
}

type Memory struct {
	vm          *VM
	ramRegions  map[uint64]RAM
	mmioRegions map[uint64]MMIORegion
}

func (memory *Memory) getRAMRegion(address uint64) (uint64, RAM, error) {
	for ramBaseAddress, ram := range memory.ramRegions {
		if address >= ramBaseAddress && address < ramBaseAddress+uint64(len(ram)) {
			return ramBaseAddress, ram, nil
		}
	}

	return 0, nil, errors.New("address not mapped")
}

func (memory *Memory) checkIntersectRegion(baseAddress, endAddress uint64) bool {
	//for ramBaseAddress, ram := range memory.ramRegions {
	//	if baseAddress >= ramBaseAddress && baseAddress < ramBaseAddress+uint64(len(ram)) {
	//		return true
	//	}
	//
	//	if endAddress >= ramBaseAddress && endAddress < uint64(len(ram)) {
	//		return true
	//	}
	//}
	//
	//for mmioBaseAddress, mmioRegion := range memory.mmioRegions {
	//	if baseAddress >= mmioBaseAddress && baseAddress < mmioRegion.Size {
	//		return true
	//	}
	//
	//	if endAddress >= mmioBaseAddress && endAddress < mmioRegion.Size {
	//		return true
	//	}
	//}

	return false
}

func (memory *Memory) AddRAMRegion(baseAddress uint64, ram RAM) error {
	if memory.checkIntersectRegion(baseAddress, baseAddress+uint64(len(ram))-1) {
		return errors.New("memory region already mapped")
	}

	_, err := memory.vm.mapMemory(baseAddress, ram)
	if err != nil {
		return err
	}

	memory.ramRegions[baseAddress] = ram

	return nil
}

func (memory *Memory) AddMMIORegion(baseAddress, size uint64) error {
	if memory.checkIntersectRegion(baseAddress, baseAddress+size-1) {
		return errors.New("memory region already mapped")
	}

	_, err := memory.vm.mapMMIO(baseAddress, size)
	if err != nil {
		return err
	}

	//TODO:
	//memory.ramRegions[baseAddress] = ram

	return nil
}

func (memory *Memory) Read(address uint64, buf []byte) {
	//TODO: write faster version

	for i := 0; i < len(buf); i++ {
		buf[i], _ = memory.Read8(address + uint64(i))
	}
}

func (memory *Memory) Read8(address uint64) (uint8, error) {
	baseAddress, ram, err := memory.getRAMRegion(address)
	if err != nil {
		return 0, err
	}

	return ram[address-baseAddress], nil
}

func (memory *Memory) Write32(address uint64, data uint32) error {
	baseAddress, ram, err := memory.getRAMRegion(address)
	if err != nil {
		return err
	}

	binary.LittleEndian.PutUint32(ram[address-baseAddress:], data)

	return nil
}

func (memory *Memory) Write16(address uint64, data uint16) error {
	baseAddress, ram, err := memory.getRAMRegion(address)
	if err != nil {
		return err
	}

	binary.LittleEndian.PutUint16(ram[address-baseAddress:], data)

	return nil
}

func (memory *Memory) Write8(address uint64, data uint8) error {
	baseAddress, ram, err := memory.getRAMRegion(address)
	if err != nil {
		return err
	}

	ram[address-baseAddress] = data

	return nil
}

func (memory *Memory) Write(address uint64, data []byte) error {
	baseAddress, ram, err := memory.getRAMRegion(address)
	if err != nil {
		return err
	}

	copy(ram[address-baseAddress:], data)

	return nil
}

func (memory *Memory) PrintMemoryRegion(baseAddress, size uint64) {
	for line := uint64(0); line < size/16; line++ {
		fmt.Printf("0x%08X:", baseAddress+line*16)

		for i := uint64(0); i < 16; i++ {
			if line*16+i < size {
				value, err := memory.Read8(baseAddress + line*16 + i)
				if err != nil {
					return
				}
				fmt.Printf("%02X ", value)
			}
		}

		fmt.Println()
	}
}
