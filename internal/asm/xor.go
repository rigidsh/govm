package asm

func XOR(to, from Argument) Command {
	result := make([]byte, 0)
	switch to.(type) {
	case ax:
		switch from.(type) {
		case ax:
			result = append(result, 0x31, 0xC0)
		default:
			panic("Not supported operation")
		}
	}

	result = append(result, to.bytes()...)
	result = append(result, from.bytes()...)

	return result
}
