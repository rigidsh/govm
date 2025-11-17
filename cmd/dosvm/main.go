package main

import (
	"fmt"
	"log"

	"github.com/rigidsh/govm/internal/bios"
	"github.com/rigidsh/govm/internal/dma"
	"github.com/rigidsh/govm/internal/fdc"
	"github.com/rigidsh/govm/internal/kvm"
)

func main() {

	kvmManager, err := kvm.OpenKVM()
	if err != nil {
		return
	}

	vm, err := kvmManager.CreateVM()
	if err != nil {
		return
	}

	dmaController := dma.CreateDMA(vm, dma.MasterPortConfig)
	floppyController := fdc.CreateFDC(vm, dmaController)

	dosDiskImage, err := fdc.OpenRaw144DiskImage("/home/alex/x86BOOT.img")
	if err != nil {
		fmt.Printf("Can't open disk image")
		return
	}

	floppyController.InsertDisk(0, dosDiskImage)

	fmt.Println("VM created.")

	err = vm.SetTSS(0xfffbd000)
	if err != nil {
		log.Fatalf("KVM_SET_TSS_ADDR failed: %v", err)
	}
	fmt.Println("TSS setup")

	err = vm.CreateIRQChip()
	if err != nil {
		log.Fatalf("Can't create IRQ chip: %v", err)
	}

	err = vm.CreatePIT2()
	if err != nil {
		log.Fatalf("Can't create PIT chip: %v", err)
	}

	lowMemory, err := kvmManager.AllocateRAM(0x100000 - 256*1024)
	if err != nil {
		log.Fatalf("Can't allocate VM memory: %v", err)
	}

	err = vm.Memory().AddRAMRegion(0, lowMemory)
	if err != nil {
		log.Fatalf("Can't map low ram", err)
	}

	kvm.SetupA20Register(vm)

	fmt.Println("Memory is ready.")

	_, err = kvm.CreateDebugPort(vm)
	if err != nil {
		log.Fatalf("Can't init debug port %t", err)
	}
	bios := bios.CreateSeaBIOS(vm)

	err = bios.Init(kvmManager)
	if err != nil {
		log.Fatalf("Can't init BIOS %t", err)
	}

	vm.Memory().PrintMemoryRegion(0x000, 128)

	highMemory, err := kvmManager.AllocateRAM(0x100000)
	if err != nil {
		log.Fatalf("Can't allocate VM memory: %v", err)
	}

	err = vm.Memory().AddRAMRegion(0x100000, highMemory)
	if err != nil {
		log.Fatalf("Can't map bios rom", err)
	}

	cpu, err := vm.CreateCPU()
	if err != nil {
		log.Fatalf("KVM_CREATE_VCPU failed: %v", err)
	}
	fmt.Println("vCPU created")

	sregs, err := cpu.GetSRegs()
	if err != nil {
		log.Fatalf("Can't get SRegs")
	}

	sregs.CS.Base = 0xF0000
	sregs.CS.Selector = 0xF000

	err = cpu.SetSRegs(sregs)
	if err != nil {
		log.Fatalf("Can't set SRegs %v", err)
	}

	regs, err := cpu.GetRegs()
	if err != nil {
		log.Fatalf("Can't get SRegs")
	}

	regs.RIP = 0xFFF0
	regs.RFlags = regs.RFlags | 0x200

	err = cpu.SetRegs(regs)
	if err != nil {
		log.Fatalf("Can't set Regs")
	}

	fmt.Println("Start vCPU...")

	vm.Memory().PrintMemoryRegion(0xF0000+0x0100, 16)

	cpu.Run(func(run *kvm.Run) bool {
		if run.ExitReason == 10 {
			return true
		}

		if run.ExitReason == 6 {
			mmio := run.GetMMIO()
			regs, _ = cpu.GetRegs()
			fmt.Printf("Address: %X, Data: %X, IsWrite: %t", mmio.PhysAddr, mmio.Data, mmio.IsWrite)
			return false
		}

		if run.ExitReason == 5 {

			return true
		}

		fmt.Printf("Exit code: %d\n", run.ExitReason)
		regs, _ = cpu.GetRegs()
		sregs, _ = cpu.GetSRegs()
		fmt.Printf("\nExecuted at: %X:%X\n", sregs.CS.Selector, regs.RIP)
		fmt.Printf("AX=%X\n", regs.RAX.X())
		fmt.Printf("BX=%X\n", regs.RBX.X())
		fmt.Printf("DS=%X\n", sregs.DS.Selector)
		fmt.Printf("ES=%X\n", sregs.ES.Selector)
		fmt.Printf("RFII=%X\n", run.ReadyForInterruptInjection)
		fmt.Printf("IfFlag=%X\n", run.IfFlag)
		vm.Memory().PrintMemoryRegion(sregs.CS.Base+regs.RIP, 16)

		return false
	})

}
