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
	OP_I32_NEG
	OP_I32_ADD
	OP_I32_SUB
	OP_I32_MUL
	OP_I32_DIV
	OP_I32_LT
	OP_I32_GT
	OP_I32_LEQ
	OP_I32_GEQ
	OP_I32_EQ

	OP_BOOL_LOAD
	OP_BOOL_AND
	OP_BOOL_OR
	OP_BOOL_NOT
	OP_BOOL_XOR

	OP_JUMP
	OP_JUMPC
	OP_JUMPNC

	OP_BLOCK_START
	OP_BLOCK_END
	OP_BLOCK_BR
	OP_BLOCK_LOAD
	OP_BLOCK_SAVE

	OP_LOAD
	OP_SAVE
	OP_NEW

	OP_FUNC_CALL
	OP_FUNC_RET

	OP_LOCAL_LOAD
	OP_LOCAL_SAVE
)

var instrKindString = map[byte]string{
	OP_NULL: "null",
	OP_HALT: "halt",

	OP_I32_LOAD: "i32.load",
	OP_I32_NEG:  "i32.neg",
	OP_I32_ADD:  "i32.add",
	OP_I32_SUB:  "i32.sub",
	OP_I32_MUL:  "i32.mul",
	OP_I32_DIV:  "i32.div",
	OP_I32_LT:   "i32.lt",
	OP_I32_GT:   "i32.gt",
	OP_I32_LEQ:  "i32.leq",
	OP_I32_GEQ:  "i32.geq",
	OP_I32_EQ:   "i32.eq",

	OP_BOOL_LOAD: "bool.load",
	OP_BOOL_NOT:  "bool.not",
	OP_BOOL_AND:  "bool.and",
	OP_BOOL_OR:   "bool.or",
	OP_BOOL_XOR:  "bool.xor",

	OP_JUMP:   "jump",
	OP_JUMPC:  "jumpc",
	OP_JUMPNC: "jumpnc",

	OP_BLOCK_START: "block.block",
	OP_BLOCK_END:   "block.end",
	OP_BLOCK_BR:    "block.br",
	OP_BLOCK_LOAD:  "block.load",
	OP_BLOCK_SAVE:  "block.save",

	OP_LOAD: "load",
	OP_NEW:  "new",
	OP_SAVE: "save",

	OP_FUNC_CALL: "func.call",
	OP_FUNC_RET:  "func.ret",

	OP_LOCAL_LOAD: "local.load",
	OP_LOCAL_SAVE: "local.save",
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
	fmt.Fprintf(&buf, "%-12s", instrKindString[i.Kind])
	switch i.Kind {
	case OP_I32_LOAD:
		obj, err := CreateI32Object(i.Operands)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(&buf, " %s", obj.String())
	case OP_BOOL_LOAD:
		obj, err := CreateBool(i.Operands)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(&buf, " %s", obj.String())
	case OP_JUMP, OP_JUMPC, OP_JUMPNC, OP_BLOCK_START, OP_FUNC_CALL:
		obj, err := CreateI32Object(i.Operands[:4])
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(&buf, " [%d]", obj.ToI32())
	case OP_LOAD, OP_BLOCK_LOAD, OP_LOCAL_LOAD, OP_SAVE, OP_BLOCK_SAVE, OP_LOCAL_SAVE:
		obj, err := CreateI32Object(i.Operands)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(&buf, " $%d", obj.ToI32())
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
func BoolLoad(x bool) Instruction {
	buf := make([]byte, 0, 2)
	buf = append(buf, TAG_BOOL)
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
func BlockStart(x uint32) Instruction {
	buf := make([]byte, 0, 5)
	buf = append(buf, TAG_I32)
	buf = binary.LittleEndian.AppendUint32(buf, uint32(x))
	return Instruction{Kind: OP_BLOCK_START, Operands: buf}
}
func BlockEnd() Instruction {
	return Instruction{Kind: OP_BLOCK_END}
}
func BlockBr() Instruction {
	return Instruction{Kind: OP_BLOCK_BR}
}
func BlockLoad(x uint32) Instruction {
	buf := make([]byte, 0, 5)
	buf = append(buf, TAG_I32)
	buf = binary.LittleEndian.AppendUint32(buf, uint32(x))
	return Instruction{Kind: OP_BLOCK_LOAD, Operands: buf}
}
func BlockSave(x uint32) Instruction {
	buf := make([]byte, 0, 5)
	buf = append(buf, TAG_I32)
	buf = binary.LittleEndian.AppendUint32(buf, uint32(x))
	return Instruction{Kind: OP_BLOCK_SAVE, Operands: buf}
}
func LocalLoad(x uint32) Instruction {
	buf := make([]byte, 0, 5)
	buf = append(buf, TAG_I32)
	buf = binary.LittleEndian.AppendUint32(buf, uint32(x))
	return Instruction{Kind: OP_LOCAL_LOAD, Operands: buf}
}
func LocalSave(x uint32) Instruction {
	buf := make([]byte, 0, 5)
	buf = append(buf, TAG_I32)
	buf = binary.LittleEndian.AppendUint32(buf, uint32(x))
	return Instruction{Kind: OP_LOCAL_SAVE, Operands: buf}
}
func Load(x uint32) Instruction {
	buf := make([]byte, 0, 5)
	buf = append(buf, TAG_I32)
	buf = binary.LittleEndian.AppendUint32(buf, uint32(x))
	return Instruction{Kind: OP_LOAD, Operands: buf}
}
func Save(x uint32) Instruction {
	buf := make([]byte, 0, 5)
	buf = append(buf, TAG_I32)
	buf = binary.LittleEndian.AppendUint32(buf, uint32(x))
	return Instruction{Kind: OP_SAVE, Operands: buf}
}
func New() Instruction {
	return Instruction{Kind: OP_NEW}
}
func FuncCall(addr uint32, args uint32) Instruction {
	buf := make([]byte, 0, 10)
	buf = append(buf, TAG_I32)
	buf = binary.LittleEndian.AppendUint32(buf, addr)
	buf = append(buf, TAG_I32)
	buf = binary.LittleEndian.AppendUint32(buf, args)
	return Instruction{Kind: OP_FUNC_CALL, Operands: buf}
}
func FuncRet(x uint32) Instruction {
	buf := make([]byte, 0, 5)
	buf = append(buf, TAG_I32)
	buf = binary.LittleEndian.AppendUint32(buf, x)
	return Instruction{Kind: OP_FUNC_RET, Operands: buf}
}
