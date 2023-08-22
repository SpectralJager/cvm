package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	OP_NULL byte = iota
	OP_HALT

	OP_I32_LOAD
	OP_I32_ADD
	OP_I32_SUB
	OP_I32_MUL
	OP_I32_DIV
	OP_I32_LT
	OP_I32_GT
	OP_I32_LEQ
	OP_I32_GEQ
	OP_I32_EQ

	OP_JUMP
	OP_JUMPC
	OP_JUMPNC

	OP_BLOCK
	OP_END
)

var instrKindString = map[byte]string{
	OP_NULL: "null",
	OP_HALT: "halt",

	OP_I32_LOAD: "i32.load",
	OP_I32_ADD:  "i32.add",
	OP_I32_SUB:  "i32.sub",
	OP_I32_MUL:  "i32.mul",
	OP_I32_DIV:  "i32.div",
	OP_I32_LT:   "i32.lt",
	OP_I32_GT:   "i32.gt",
	OP_I32_LEQ:  "i32.leq",
	OP_I32_GEQ:  "i32.geq",
	OP_I32_EQ:   "i32.eq",

	OP_JUMP:   "jump",
	OP_JUMPC:  "jumpc",
	OP_JUMPNC: "jumpnc",

	OP_BLOCK: "block",
	OP_END:   "end",
}

var (
	ErrReachHalt = errors.New("reach halt instruction")
)

type Instruction struct {
	Kind     byte
	Operands []byte
}

func (i *Instruction) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s", instrKindString[i.Kind])
	if len(i.Operands) > 0 {
		fmt.Fprintf(&buf, " %v", i.Operands)
	}
	return buf.String()
}

func Null() Instruction {
	return Instruction{Kind: OP_NULL}
}
func Halt() Instruction {
	return Instruction{Kind: OP_HALT}
}
func I32Load(x int32) Instruction {
	buf := make([]byte, 0, 5)
	buf = append(buf, TAG_I32)
	buf = binary.LittleEndian.AppendUint32(buf, uint32(x))
	return Instruction{Kind: OP_I32_LOAD, Operands: buf}
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
func Jump(x uint32) Instruction {
	buf := make([]byte, 0, 5)
	buf = append(buf, TAG_I32)
	buf = binary.LittleEndian.AppendUint32(buf, uint32(x))
	return Instruction{Kind: OP_JUMP, Operands: buf}
}
func JumpC(x uint32) Instruction {
	buf := make([]byte, 0, 5)
	buf = append(buf, TAG_I32)
	buf = binary.LittleEndian.AppendUint32(buf, uint32(x))
	return Instruction{Kind: OP_JUMPC, Operands: buf}
}
func JumpNC(x uint32) Instruction {
	buf := make([]byte, 0, 5)
	buf = append(buf, TAG_I32)
	buf = binary.LittleEndian.AppendUint32(buf, uint32(x))
	return Instruction{Kind: OP_JUMPNC, Operands: buf}
}
func Block() Instruction {
	return Instruction{Kind: OP_BLOCK}
}
func End() Instruction {
	return Instruction{Kind: OP_BLOCK}
}
