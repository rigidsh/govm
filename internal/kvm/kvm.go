package kvm

import (
	"golang.org/x/sys/unix"
)

func OpenKVM() (*KVM, error) {
	fd, err := unix.Open("/dev/kvm", unix.O_RDWR|unix.O_CLOEXEC, 0)
	if err != nil {
		return nil, err
	}

	return &KVM{
		fd: fd,
	}, nil
}

type KVM struct {
	fd int
}

func (kvm *KVM) getVCpuRunStructSize() (uintptr, error) {
	runSize, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(kvm.fd), KVM_GET_VCPU_MMAP_SIZE, 0)
	if errno != 0 {
		return 0, errno
	}

	return runSize, nil
}

func (kvm *KVM) AllocateRAM(size int) (RAM, error) {
	data, err := unix.Mmap(
		-1, 0, size,
		unix.PROT_READ|unix.PROT_WRITE, unix.MAP_SHARED|unix.MAP_ANON)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (kvm *KVM) CreateVM() (*VM, error) {
	fd, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(kvm.fd), KVM_CREATE_VM, 0)
	if errno != 0 {
		return nil, errno
	}

	vm := &VM{
		fd:    int(fd),
		kvm:   kvm,
		ports: make(map[uint16]IOPort),
	}

	vm.memory = newMemory(vm)

	return vm, nil
}
