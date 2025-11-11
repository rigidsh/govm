package asm

func HLT() Command {
	return []byte{0xF4}
}
