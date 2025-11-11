package kvm

type KvmUserspaceMemoryRegion struct {
	Slot          uint32
	Flags         uint32
	GuestPhysAddr uint64
	MemorySize    uint64
	UserspaceAddr uint64
}

type KvmSeg struct {
	Base     uint64
	Limit    uint32
	Selector uint16
	Type     uint8
	_        [8]uint8
	_        uint8
}

type KvmPitConfig struct {
	Flags uint32
	_     [15]uint32
}

type KvmDTable struct {
	Base  uint64
	Limit uint16
	_     [3]uint16
}

type KvmInterrupt struct {
	Irq uint32
}

type KvmIrqRouting struct {
	Nr    uint32
	Flags uint32
	//Entries
}

type KvmSRegs struct {
	CS              KvmSeg
	DS              KvmSeg
	ES              KvmSeg
	FS              KvmSeg
	GS              KvmSeg
	SS              KvmSeg
	TR              KvmSeg
	LDT             KvmSeg
	GDT             KvmDTable
	IDT             KvmDTable
	CR0             uint64
	CR2             uint64
	CR3             uint64
	CR4             uint64
	CR8             uint64
	Efer            uint64
	ApicBase        uint64
	InterruptBitmap [(256 + 63) / 64]uint64
}

type KvmRegs struct {
	RAX    Register64
	RBX    Register64
	RCX    Register64
	RDX    Register64
	RSI    uint64
	RDI    uint64
	RSP    uint64
	RBP    uint64
	R8     uint64
	R9     uint64
	R10    uint64
	R11    uint64
	R12    uint64
	R13    uint64
	R14    uint64
	R15    uint64
	RIP    uint64
	RFlags uint64
}

type Register64 uint64

func (r Register64) H() uint8 {
	return uint8(r & 0xFF00 >> 8)
}
func (r Register64) L() uint8 {
	return uint8(r & 0xFF)
}

func (r Register64) X() uint16 {
	return uint16(r & 0xFFFF)
}
