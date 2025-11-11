package asm

func JMP(addr Argument) Command {
	result := make([]byte, 0)
	switch addr.(type) {
	case segAddr:
		result = append(result, 0xEA)
	default:
		panic("Not supported operation")
	}

	result = append(result, addr.bytes()...)

	return result
}
