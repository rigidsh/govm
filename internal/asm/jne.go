package asm

func JNE(addr Argument) Command {
	result := make([]byte, 0)
	switch addr.(type) {
	case relAddr:
		result = append(result, 0x75)
	default:
		panic("Not supported operation")
	}

	result = append(result, addr.bytes()...)

	return result
}
