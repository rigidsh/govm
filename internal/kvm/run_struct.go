package kvm

import "unsafe"

type Run struct {
	RequestInterruptWindow     uint8
	_                          [7]byte
	ExitReason                 uint32 // Причина выхода
	ReadyForInterruptInjection uint8
	IfFlag                     uint8
	_                          [18]byte
	Other                      [4096]uint8
}

func (run *Run) GetIO() *IO {
	return (*IO)(unsafe.Pointer(&run.Other[0]))
}

func (run *Run) GetMMIO() *MMIO {
	return (*MMIO)(unsafe.Pointer(&run.Other[0]))
}

func (run *Run) Get(offset uint64) uint8 {
	return run.Other[offset-32]
}

func (run *Run) Read(offset uint64, size uint8) []byte {
	return run.Other[offset-32 : offset-32+uint64(size)]
}

func (run *Run) Write(offset uint64, data []byte) {
	copy(run.Other[offset-32:], data)
}

type IO struct { // Для KVM_EXIT_IO
	Direction  uint8
	Size       uint8
	Port       uint16
	Count      uint32
	DataOffset uint64
}

type MMIO struct {
	PhysAddr uint64
	Data     [8]uint8
	Len      uint32
	IsWrite  bool
}
