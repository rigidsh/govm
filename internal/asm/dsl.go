package asm

type Command []byte

// func OUT_AL uint8 = 0xE6
// func MOV_AL uint8 = 0xB0
// func MOV_AH uint8 = 0xB4
// func MOV_AX uint8 = 0xA1
// func CMP_AX uint8 = MOV_AX
//func INT(number uint8) Command     { return []byte{0xCD, byte(number)} }
//func HLT() Command                 { return []byte{0xF4} }
//func JNE_REL(address int8) Command { return []byte{0x75, byte(address)} }
//func JMP() Command                 { return []byte{0xEA} }
//func IRET() Command                { return []byte{0xCF} }

func Compile(commands ...Command) []byte {
	result := make([]uint8, 0)

	for _, command := range commands {
		result = append(result, command...)
	}

	return result
}
