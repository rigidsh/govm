package kvm

import (
	"bytes"
	"os"
)

func CreateDebugPort(vm *VM) (*ComPort, error) {
	return CreateComPort(vm, 0x402, os.Stdout, bytes.NewBuffer([]byte{0xE9}))
}
