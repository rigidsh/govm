package kvm

import (
	"log"
	"unsafe"

	"golang.org/x/sys/unix"
)

type VM struct {
	kvm            *KVM
	fd             int
	nextMemorySlot uint32
	memory         *Memory
	ports          map[uint16]IOPort
	cpuList        []*VCPU
}

func (vm *VM) Memory() *Memory {
	return vm.memory
}

func (vm *VM) CreateCPU() (*VCPU, error) {
	fd, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(vm.fd), KVM_CREATE_VCPU, 0)

	if errno != 0 {
		return nil, errno
	}

	runSize, err := vm.kvm.getVCpuRunStructSize()
	if err != nil {
		return nil, err
	}

	runMap, err := unix.Mmap(
		int(fd), 0, int(runSize),
		unix.PROT_READ|unix.PROT_WRITE, unix.MAP_SHARED)
	if err != nil {
		log.Fatalf("Mmap kvm_run failed: %v", err)
	}
	runStruct := (*Run)(unsafe.Pointer(&runMap[0]))

	cpu := &VCPU{
		fd:        int(fd),
		runStruct: runStruct,
		vm:        vm,
	}
	vm.cpuList = append(vm.cpuList, cpu)

	return cpu, nil
}

func (vm *VM) SetTSS(address uintptr) error {
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(vm.fd), KVM_SET_TSS_ADDR, address)
	if errno != 0 {
		return errno
	}
	return nil
}

func (vm *VM) CreateIRQChip() error {
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(vm.fd), KVM_CREATE_IRQCHIP, 0)
	if errno != 0 {
		return errno
	}
	return nil
}

func (vm *VM) CreatePIT2() error {
	config := KvmPitConfig{
		Flags: 1,
	}
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(vm.fd), KVM_CREATE_PIT2, uintptr(unsafe.Pointer(&config)))
	if errno != 0 {
		return errno
	}
	return nil
}

func (vm *VM) RegisterPortHandler(portNumber uint16, handler IOPort) {
	vm.ports[portNumber] = handler
}

func (vm *VM) mapMemory(baseAddress uint64, ram RAM) (uint32, error) {
	memRegion := KvmUserspaceMemoryRegion{
		Slot:          vm.nextMemorySlot,
		Flags:         0,
		GuestPhysAddr: baseAddress,
		MemorySize:    uint64(len(ram)),
		UserspaceAddr: uint64(uintptr(unsafe.Pointer(&ram[0]))),
	}

	_, _, errno := unix.Syscall(
		unix.SYS_IOCTL, uintptr(vm.fd), KVM_SET_USER_MEMORY_REGION,
		uintptr(unsafe.Pointer(&memRegion)))

	if errno != 0 {
		return 0, errno
	}

	vm.nextMemorySlot++

	return memRegion.Slot, nil
}

func (vm *VM) mapMMIO(baseAddress uint64, size uint64) (uint32, error) {
	memRegion := KvmUserspaceMemoryRegion{
		Slot:          vm.nextMemorySlot,
		Flags:         0,
		GuestPhysAddr: baseAddress,
		MemorySize:    size,
		UserspaceAddr: 0,
	}

	_, _, errno := unix.Syscall(
		unix.SYS_IOCTL, uintptr(vm.fd), KVM_SET_USER_MEMORY_REGION,
		uintptr(unsafe.Pointer(&memRegion)))

	if errno != 0 {
		return 0, errno
	}

	vm.nextMemorySlot++

	return memRegion.Slot, nil
}
