package instruction

import (
	"cvm/object"
	"encoding/binary"
)

func I32Load(x int32) Instruction {
	buf := make([]byte, 0, 5)
	buf = append(buf, object.TAG_I32)
	buf = binary.LittleEndian.AppendUint32(buf, uint32(x))
	return Instruction{Kind: OP_I32_LOAD, Operands: buf}
}
func I32Neg() Instruction {
	return Instruction{Kind: OP_I32_NEG}
}
func I32Add() Instruction {
	return Instruction{Kind: OP_I32_ADD}
}
func I32Sub() Instruction {
	return Instruction{Kind: OP_I32_SUB}
}
func I32Mul() Instruction {
	return Instruction{Kind: OP_I32_MUL}
}
func I32Div() Instruction {
	return Instruction{Kind: OP_I32_DIV}
}
func I32Lt() Instruction {
	return Instruction{Kind: OP_I32_LT}
}
func I32Gt() Instruction {
	return Instruction{Kind: OP_I32_GT}
}
func I32Leq() Instruction {
	return Instruction{Kind: OP_I32_LEQ}
}
func I32Geq() Instruction {
	return Instruction{Kind: OP_I32_GEQ}
}
func I32Eq() Instruction {
	return Instruction{Kind: OP_I32_EQ}
}
