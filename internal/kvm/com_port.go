package kvm

import "io"

type ComPort struct {
	vm *VM
}

func CreateComPort(vm *VM, portNumber uint16, out io.Writer, in io.Reader) (*ComPort, error) {
	comPort := &ComPort{
		vm: vm,
	}

	vm.RegisterPortHandler(portNumber, CreateIOPort(in, out))

	return comPort, nil
}
