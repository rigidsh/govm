package kvm_testing

import (
	"testing"

	"github.com/rigidsh/govm/internal/asm"
	"github.com/rigidsh/govm/internal/kvm"
)

func TestKVMTesting_SimpleTest(t *testing.T) {
	result, err := Do(
		func(vm *kvm.VM) {

		},
		asm.Compile(
			asm.MOV(asm.AL(), asm.Impl8(0x42)),
			WRITE_TEST_RESULT_FROM_AL(),
			STOP_TEST(0),
		),
	)
	if err != nil {
		t.Error(err)
		return
	}
	if result[0] != 0x42 {
		t.Errorf("Incorrect text result, expected %X, get %X", 0x42, result[0])
	}
}
