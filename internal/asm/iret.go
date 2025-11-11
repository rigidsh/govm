package asm

func IRET() Command {
	return []byte{0xCF}
}
