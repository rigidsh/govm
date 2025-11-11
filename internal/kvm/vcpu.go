package kvm

import (
	"fmt"
	"runtime"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

type VCPU struct {
	fd          int
	vm          *VM
	runStruct   *Run
	osThreadPID int
}

func (cpu *VCPU) GetSRegs() (*KvmSRegs, error) {
	var result KvmSRegs
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(cpu.fd), KVM_GET_SREGS, uintptr(unsafe.Pointer(&result)))
	if errno != 0 {
		return nil, errno
	}
	return &result, nil
}

func (cpu *VCPU) GetRegs() (*KvmRegs, error) {
	var result KvmRegs
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(cpu.fd), KVM_GET_REGS, uintptr(unsafe.Pointer(&result)))
	if errno != 0 {
		return nil, errno
	}
	return &result, nil
}

func (cpu *VCPU) SetSRegs(sregs *KvmSRegs) error {
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(cpu.fd), KVM_SET_SREGS, uintptr(unsafe.Pointer(sregs)))
	if errno != 0 {
		return errno
	}
	return nil
}

func (cpu *VCPU) SetRegs(regs *KvmRegs) error {
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(cpu.fd), KVM_SET_REGS, uintptr(unsafe.Pointer(regs)))
	if errno != 0 {
		return errno
	}
	return nil
}

func (cpu *VCPU) Interrupt(irq uint32) error {
	interruptStruct := &KvmInterrupt{
		Irq: irq,
	}
	cpu.runStruct.RequestInterruptWindow = 1
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(cpu.fd), KVM_INTERRUPT, uintptr(unsafe.Pointer(interruptStruct)))
	if errno != 0 {
		return errno
	}
	return nil

}

const SYS_TGKILL = 234

func (cpu *VCPU) Stop() error {
	fmt.Println(syscall.Getpid())
	fmt.Println(cpu.osThreadPID)
	_, _, errno := syscall.Syscall(
		SYS_TGKILL,
		uintptr(syscall.Getpid()),
		uintptr(cpu.osThreadPID),
		uintptr(syscall.SIGUSR1))
	if errno != 0 {
		return errno
	}

	return nil
}

func (cpu *VCPU) Run(callback func(run *Run) bool) error {
	runtime.LockOSThread()
	cpu.osThreadPID = syscall.Gettid()
	defer runtime.UnlockOSThread()
	for {
		_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(cpu.fd), KVM_RUN, 0)
		if errno != 0 && errno != unix.EINTR {
			fmt.Println("!!!!!!!")
			return errno
		}

		if cpu.runStruct.ExitReason == 2 {
			io := cpu.runStruct.GetIO()
			handler, ok := cpu.vm.ports[io.Port]
			if ok {
				if io.Direction == 1 {
					handler.OnWrite(cpu.runStruct.Read(io.DataOffset, io.Size))
				} else {
					cpu.runStruct.Write(io.DataOffset, handler.OnRead(io.Size))
				}
				continue
			}
			//fmt.Printf("IO %t: %X\n", io.Direction == 1, io.Port)
		}

		if !callback(cpu.runStruct) {
			return nil
		}
	}
}
