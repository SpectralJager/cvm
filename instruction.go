package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
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
	OP_I32_TO_F32
	OP_I32_TO_BOOL

	OP_BOOL_LOAD
	OP_BOOL_AND
	OP_BOOL_OR
	OP_BOOL_NOT
	OP_BOOL_XOR

	OP_F32_LOAD
	OP_F32_NEG
	OP_F32_ADD
	OP_F32_SUB
	OP_F32_MUL
	OP_F32_DIV
	OP_F32_LT
	OP_F32_GT
	OP_F32_LEQ
	OP_F32_GEQ
	OP_F32_EQ
	OP_F32_TO_I32
	OP_F32_TO_BOOL

	OP_LIST_NEW
	OP_LIST_APPEND
	OP_LIST_GET
	OP_LIST_POP
	OP_LIST_INSERT
	OP_LIST_REPLACE
	OP_LIST_REMOVE

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
	OP_FREE
	OP_POP

	OP_FUNC_CALL
	OP_FUNC_RET

	OP_LOCAL_LOAD
	OP_LOCAL_SAVE
)

var instrKindString = map[byte]string{
	OP_NULL: "null",
	OP_HALT: "halt",

	OP_I32_LOAD:    "i32.load",
	OP_I32_NEG:     "i32.neg",
	OP_I32_ADD:     "i32.add",
	OP_I32_SUB:     "i32.sub",
	OP_I32_MUL:     "i32.mul",
	OP_I32_DIV:     "i32.div",
	OP_I32_LT:      "i32.lt",
	OP_I32_GT:      "i32.gt",
	OP_I32_LEQ:     "i32.leq",
	OP_I32_GEQ:     "i32.geq",
	OP_I32_EQ:      "i32.eq",
	OP_I32_TO_F32:  "i32.to_f32",
	OP_I32_TO_BOOL: "i32.to_bool",

	OP_BOOL_LOAD: "bool.load",
	OP_BOOL_NOT:  "bool.not",
	OP_BOOL_AND:  "bool.and",
	OP_BOOL_OR:   "bool.or",
	OP_BOOL_XOR:  "bool.xor",

	OP_F32_LOAD:    "f32.load",
	OP_F32_NEG:     "f32.neg",
	OP_F32_ADD:     "f32.add",
	OP_F32_SUB:     "f32.sub",
	OP_F32_MUL:     "f32.mul",
	OP_F32_DIV:     "f32.div",
	OP_F32_LT:      "f32.lt",
	OP_F32_GT:      "f32.gt",
	OP_F32_LEQ:     "f32.leq",
	OP_F32_GEQ:     "f32.geq",
	OP_F32_EQ:      "f32.eq",
	OP_F32_TO_I32:  "f32.to_i32",
	OP_F32_TO_BOOL: "f32.to_bool",

	OP_LIST_NEW:     "list.new",
	OP_LIST_APPEND:  "list.append",
	OP_LIST_GET:     "list.get",
	OP_LIST_POP:     "list.pop",
	OP_LIST_REMOVE:  "list.remove",
	OP_LIST_INSERT:  "list.insert",
	OP_LIST_REPLACE: "list.replace",

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
	OP_POP:  "pop",
	OP_FREE: "free",

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
		obj, err := CreateObject(i.Operands)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(&buf, " %s", obj.String())
	case OP_BOOL_LOAD:
		obj, err := CreateObject(i.Operands)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(&buf, " %s", obj.String())
	case OP_F32_LOAD:
		obj, err := CreateObject(i.Operands)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(&buf, " %s", obj.String())
	case OP_JUMP, OP_JUMPC, OP_JUMPNC, OP_BLOCK_START, OP_FUNC_CALL, OP_LIST_GET:
		obj, err := CreateObject(i.Operands[:4])
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(&buf, " [%d]", obj.ToI32())
	case OP_LOAD, OP_BLOCK_LOAD, OP_LOCAL_LOAD, OP_SAVE, OP_BLOCK_SAVE, OP_LOCAL_SAVE:
		obj, err := CreateObject(i.Operands)
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
func I32ToF32() Instruction {
	return Instruction{Kind: OP_I32_TO_F32}
}
func I32ToBool() Instruction {
	return Instruction{Kind: OP_I32_TO_BOOL}
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
func F32Load(x float32) Instruction {
	buf := make([]byte, 0, 5)
	buf = append(buf, TAG_F32)
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
func F32ToI32() Instruction {
	return Instruction{Kind: OP_F32_TO_I32}
}
func F32ToBool() Instruction {
	return Instruction{Kind: OP_F32_TO_BOOL}
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
func Pop() Instruction {
	return Instruction{Kind: OP_POP}
}
func Free(x uint32) Instruction {
	buf := make([]byte, 0, 5)
	buf = append(buf, TAG_I32)
	buf = binary.LittleEndian.AppendUint32(buf, uint32(x))
	return Instruction{Kind: OP_FREE, Operands: buf}
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
func ListNew(it byte) Instruction {
	buf := make([]byte, 0, 6)
	buf = append(buf, it)
	buf = append(buf, TAG_I32)
	buf = binary.LittleEndian.AppendUint32(buf, 0)
	return Instruction{Kind: OP_LIST_NEW, Operands: buf}
}
func ListAppend() Instruction {
	return Instruction{Kind: OP_LIST_APPEND}
}
func ListGet(x uint32) Instruction {
	buf := make([]byte, 0, 5)
	buf = append(buf, TAG_I32)
	buf = binary.LittleEndian.AppendUint32(buf, x)
	return Instruction{Kind: OP_LIST_GET, Operands: buf}
}
func ListPop() Instruction {
	return Instruction{Kind: OP_LIST_POP}
}
func ListInsert(index uint32) Instruction {
	buf := make([]byte, 0, 10)
	buf = append(buf, TAG_I32)
	buf = binary.LittleEndian.AppendUint32(buf, index)
	return Instruction{Kind: OP_LIST_INSERT, Operands: buf}
}
func ListRemove(index uint32) Instruction {
	buf := make([]byte, 0, 10)
	buf = append(buf, TAG_I32)
	buf = binary.LittleEndian.AppendUint32(buf, index)
	return Instruction{Kind: OP_LIST_REMOVE, Operands: buf}
}
func ListReplace(index uint32) Instruction {
	buf := make([]byte, 0, 10)
	buf = append(buf, TAG_I32)
	buf = binary.LittleEndian.AppendUint32(buf, index)
	return Instruction{Kind: OP_LIST_REPLACE, Operands: buf}
}
