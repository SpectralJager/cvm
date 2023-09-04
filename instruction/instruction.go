package instruction

import (
	"bytes"
	"cvm/object"
	"encoding/binary"
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
	OP_I32_NEQ

	OP_BOOL_LOAD
	OP_BOOL_AND
	OP_BOOL_OR
	OP_BOOL_NOT
	OP_BOOL_NAND
	OP_BOOL_NOR
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
	OP_F32_NEQ

	OP_LIST_NEW
	OP_LIST_LENGTH
	OP_LIST_GET
	OP_LIST_INSERT
	OP_LIST_REPLACE
	OP_LIST_REMOVE

	OP_STRING_LOAD
	OP_STRING_CONCAT
	OP_STRING_FORMAT
	OP_STRING_LENGTH
	OP_STRING_SPLIT

	OP_STRUCT_LOAD

	OP_TO_STRING
	OP_TO_I32
	OP_TO_F32
	OP_TO_BOOL

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

	OP_PRINT
	OP_PRINTF
	OP_PRINTLN
	OP_READ

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
	OP_I32_NEQ:  "i32.neq",

	OP_BOOL_LOAD: "bool.load",
	OP_BOOL_NOT:  "bool.not",
	OP_BOOL_AND:  "bool.and",
	OP_BOOL_OR:   "bool.or",
	OP_BOOL_NAND: "bool.nand",
	OP_BOOL_NOR:  "bool.nor",
	OP_BOOL_XOR:  "bool.xor",

	OP_F32_LOAD: "f32.load",
	OP_F32_NEG:  "f32.neg",
	OP_F32_ADD:  "f32.add",
	OP_F32_SUB:  "f32.sub",
	OP_F32_MUL:  "f32.mul",
	OP_F32_DIV:  "f32.div",
	OP_F32_LT:   "f32.lt",
	OP_F32_GT:   "f32.gt",
	OP_F32_LEQ:  "f32.leq",
	OP_F32_GEQ:  "f32.geq",
	OP_F32_EQ:   "f32.eq",
	OP_F32_NEQ:  "f32.neq",

	OP_LIST_NEW:     "list.new",
	OP_LIST_LENGTH:  "list.length",
	OP_LIST_GET:     "list.get",
	OP_LIST_REMOVE:  "list.remove",
	OP_LIST_INSERT:  "list.insert",
	OP_LIST_REPLACE: "list.replace",

	OP_STRING_LOAD:   "string.load",
	OP_STRING_CONCAT: "string.concat",
	OP_STRING_SPLIT:  "string.split",
	OP_STRING_FORMAT: "string.format",
	OP_STRING_LENGTH: "string.length",

	OP_STRUCT_LOAD: "struct.load",

	OP_TO_STRING: "to_string",
	OP_TO_BOOL:   "to_bool",
	OP_TO_I32:    "to_i32",
	OP_TO_F32:    "to_f32",

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

	OP_PRINT:   "print",
	OP_PRINTF:  "printf",
	OP_PRINTLN: "println",
	OP_READ:    "read",

	OP_FUNC_CALL: "func.call",
	OP_FUNC_RET:  "func.ret",

	OP_LOCAL_LOAD: "local.load",
	OP_LOCAL_SAVE: "local.save",
}

type Instruction struct {
	Kind     byte
	Operands []byte
}

func (i *Instruction) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%-12s", instrKindString[i.Kind])
	switch i.Kind {
	case OP_I32_LOAD, OP_F32_LOAD, OP_BOOL_LOAD, OP_STRING_LOAD, OP_STRUCT_LOAD:
		obj, err := object.CreateObject(i.Operands)
		if err != nil {
			panic(err)
		}
		str, err := object.String(obj)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(&buf, " %s", str)
	case OP_JUMP, OP_JUMPC, OP_JUMPNC, OP_BLOCK_START, OP_FUNC_CALL:
		obj, err := object.CreateObject(i.Operands[:4])
		if err != nil {
			panic(err)
		}
		val, err := object.ValueI32(obj)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(&buf, " [%d]", val)
	case OP_LOAD, OP_BLOCK_LOAD, OP_LOCAL_LOAD, OP_SAVE, OP_BLOCK_SAVE, OP_LOCAL_SAVE:
		obj, err := object.CreateObject(i.Operands)
		if err != nil {
			panic(err)
		}
		val, err := object.ValueI32(obj)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(&buf, " $%d", val)
	}
	return buf.String()
}

func Null() Instruction {
	return Instruction{Kind: OP_NULL}
}

func Halt() Instruction {
	return Instruction{Kind: OP_HALT}
}

func Jump(x uint32) Instruction {
	buf := make([]byte, 0, 5)
	buf = append(buf, object.TAG_I32)
	buf = binary.LittleEndian.AppendUint32(buf, uint32(x))
	return Instruction{Kind: OP_JUMP, Operands: buf}
}

func JumpC(x uint32) Instruction {
	buf := make([]byte, 0, 5)
	buf = append(buf, object.TAG_I32)
	buf = binary.LittleEndian.AppendUint32(buf, uint32(x))
	return Instruction{Kind: OP_JUMPC, Operands: buf}
}

func JumpNC(x uint32) Instruction {
	buf := make([]byte, 0, 5)
	buf = append(buf, object.TAG_I32)
	buf = binary.LittleEndian.AppendUint32(buf, uint32(x))
	return Instruction{Kind: OP_JUMPNC, Operands: buf}
}

func LocalLoad(x uint32) Instruction {
	buf := make([]byte, 0, 5)
	buf = append(buf, object.TAG_I32)
	buf = binary.LittleEndian.AppendUint32(buf, uint32(x))
	return Instruction{Kind: OP_LOCAL_LOAD, Operands: buf}
}

func LocalSave(x uint32) Instruction {
	buf := make([]byte, 0, 5)
	buf = append(buf, object.TAG_I32)
	buf = binary.LittleEndian.AppendUint32(buf, uint32(x))
	return Instruction{Kind: OP_LOCAL_SAVE, Operands: buf}
}

func Load(x uint32) Instruction {
	buf := make([]byte, 0, 5)
	buf = append(buf, object.TAG_I32)
	buf = binary.LittleEndian.AppendUint32(buf, uint32(x))
	return Instruction{Kind: OP_LOAD, Operands: buf}
}

func Save(x uint32) Instruction {
	buf := make([]byte, 0, 5)
	buf = append(buf, object.TAG_I32)
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
	buf = append(buf, object.TAG_I32)
	buf = binary.LittleEndian.AppendUint32(buf, uint32(x))
	return Instruction{Kind: OP_FREE, Operands: buf}
}

func ToString() Instruction {
	return Instruction{Kind: OP_TO_STRING}
}

func ToI32() Instruction {
	return Instruction{Kind: OP_TO_I32}
}

func ToF32() Instruction {
	return Instruction{Kind: OP_TO_F32}
}

func ToBool() Instruction {
	return Instruction{Kind: OP_TO_BOOL}
}

func Print() Instruction {
	return Instruction{Kind: OP_PRINT}
}

func Printf() Instruction {
	return Instruction{Kind: OP_PRINTF}
}

func Println() Instruction {
	return Instruction{Kind: OP_PRINTLN}
}

func Read() Instruction {
	return Instruction{Kind: OP_READ}
}
