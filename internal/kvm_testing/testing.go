package kvm_testing

import (
	"errors"
	"fmt"

	"github.com/rigidsh/govm/internal/asm"
	"github.com/rigidsh/govm/internal/kvm"
)

func STOP_TEST(exitCode byte) asm.Command {
	return []byte{
		0xB0, exitCode,
		0xE6, 0xFF,
	}
}

func WRITE_TEST_RESULT_FROM_AL() asm.Command {
	return []byte{
		0xE6, 0xFE,
	}
}

func WRITE_TEST_RESULT_REGION(fromAddress, count uint16) asm.Command {
	return asm.Compile(
		asm.MOV(asm.BX(), asm.Impl16(fromAddress)),
		asm.MOV(asm.SI(), asm.Impl16(0)),
		asm.MOV(asm.AL(), asm.BXRelSIAddr()),
		asm.OUT(asm.Impl8(0xFE), asm.AL()),
		asm.INC(asm.SI()),
		asm.CMP(asm.SI(), asm.Impl16(count)),
		asm.JNE(asm.RelAddr(-9)),
	)

}

func Do(setupCallback func(*kvm.VM), code []byte) ([]byte, error) {
	kvmManager, err := kvm.OpenKVM()
	if err != nil {
		return nil, err
	}

	vm, err := kvmManager.CreateVM()
	if err != nil {
		return nil, err
	}

	ram, err := kvmManager.AllocateRAM(0x100000)
	if err != nil {
		return nil, err
	}

	copy(ram[0x500:], code)

	err = vm.Memory().AddRAMRegion(0x0000, ram)
	if err != nil {
		return nil, err
	}

	cpu, err := vm.CreateCPU()
	if err != nil {
		return nil, err
	}

	sregs, err := cpu.GetSRegs()
	if err != nil {
		return nil, err
	}

	sregs.CS.Base = 0x0
	sregs.CS.Selector = 0x0

	err = cpu.SetSRegs(sregs)
	if err != nil {
		return nil, err
	}

	regs, err := cpu.GetRegs()
	if err != nil {
		return nil, err
	}

	regs.RIP = 0x500
	regs.RFlags = regs.RFlags | 0x200

	err = cpu.SetRegs(regs)
	if err != nil {
		return nil, err
	}

	setupCallback(vm)

	testResult := make([]byte, 0)

	vm.RegisterPortHandler(0xFE, kvm.CallbackPort(
		func(write bool, data []byte) []byte {
			if write {
				testResult = append(testResult, data...)
			}
			return data
		}))

	exitCode := byte(0)

	err = cpu.Run(func(run *kvm.Run) bool {
		if run.ExitReason == 2 && run.GetIO().Port == 0xFF {
			fmt.Println("Write test result")
			exitCode = run.Read(run.GetIO().DataOffset, run.GetIO().Size)[0]
			return false
		}

		vm.Memory().PrintMemoryRegion(sregs.CS.Base+regs.RIP, 16)

		return true
	})

	if err != nil {
		return nil, err
	}

	if exitCode != 0 {
		return nil, errors.New(fmt.Sprintf("exit with code %X", exitCode))
	}

	vm.Memory().PrintMemoryRegion(0x526, 128)
	//TODO: destroy vm and free resources

	return testResult, nil
}
