package instruction

import "cvm/object"

func BoolLoad(x bool) Instruction {
	buf := make([]byte, 0, 2)
	buf = append(buf, object.TAG_BOOL)
	if x {
		buf = append(buf, 1)
	} else {
		buf = append(buf, 0)
	}
	return Instruction{Kind: OP_BOOL_LOAD, Operands: buf}
}
func BoolNot() Instruction {
	return Instruction{Kind: OP_BOOL_NOT}
}
func BoolAnd() Instruction {
	return Instruction{Kind: OP_BOOL_AND}
}
func BoolOr() Instruction {
	return Instruction{Kind: OP_BOOL_OR}
}
func BoolXor() Instruction {
	return Instruction{Kind: OP_BOOL_XOR}
}
