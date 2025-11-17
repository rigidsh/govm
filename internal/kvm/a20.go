package kvm

func SetupA20Register(vm *VM) {
	var a20Register = uint8(0x00)
	a20RegisterPort := CreateRegister8Port(&a20Register)
	vm.RegisterPortHandler(0x92, a20RegisterPort)
}
