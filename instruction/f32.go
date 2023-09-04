package instruction

import (
	"cvm/object"
	"encoding/binary"
	"math"
)

func F32Load(x float32) Instruction {
	buf := make([]byte, 0, 5)
	buf = append(buf, object.TAG_F32)
	buf = binary.LittleEndian.AppendUint32(buf, math.Float32bits(x))
	return Instruction{Kind: OP_F32_LOAD, Operands: buf}
}
func F32Neg() Instruction {
	return Instruction{Kind: OP_F32_NEG}
}
func F32Add() Instruction {
	return Instruction{Kind: OP_F32_ADD}
}
func F32Sub() Instruction {
	return Instruction{Kind: OP_F32_SUB}
}
func F32Mul() Instruction {
	return Instruction{Kind: OP_F32_MUL}
}
func F32Div() Instruction {
	return Instruction{Kind: OP_F32_DIV}
}
func F32Lt() Instruction {
	return Instruction{Kind: OP_F32_LT}
}
func F32Gt() Instruction {
	return Instruction{Kind: OP_F32_GT}
}
func F32Leq() Instruction {
	return Instruction{Kind: OP_F32_LEQ}
}
func F32Geq() Instruction {
	return Instruction{Kind: OP_F32_GEQ}
}
func F32Eq() Instruction {
	return Instruction{Kind: OP_F32_EQ}
}
func F32Neq() Instruction {
	return Instruction{Kind: OP_F32_NEQ}
}
