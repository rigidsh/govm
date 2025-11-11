package asm

func INC(variable Argument) Command {
	result := make([]byte, 0)
	switch variable.(type) {
	case bl:
		result = append(result, 0xFE, 0xC3)
	case si:
		result = append(result, 0x46)
	default:
		panic("Not supported operation")
	}

	result = append(result, variable.bytes()...)

	return result
}
