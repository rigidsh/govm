package asm

func IN(to, port Argument) Command {
	result := make([]byte, 0)
	switch to.(type) {
	case al:
		switch port.(type) {
		case impl8:
			result = append(result, 0xE4)
		case dx:
			result = append(result, 0xEC)
		default:
			panic("Not supported operation")
		}
	default:
		panic("Not supported operation")
	}

	result = append(result, port.bytes()...)
	result = append(result, to.bytes()...)

	return result
}
