package instruction

import (
	"cvm/object"
	"encoding/binary"
)

func ListLength() Instruction {
	return Instruction{Kind: OP_LIST_LENGTH}
}
func ListNew(it byte) Instruction {
	buf := make([]byte, 0, 6)
	buf = append(buf, it)
	buf = append(buf, object.TAG_I32)
	buf = binary.LittleEndian.AppendUint32(buf, 0)
	return Instruction{Kind: OP_LIST_NEW, Operands: buf}
}
func ListGet() Instruction {
	return Instruction{Kind: OP_LIST_GET}
}
func ListInsert() Instruction {
	return Instruction{Kind: OP_LIST_INSERT}
}
func ListRemove() Instruction {
	return Instruction{Kind: OP_LIST_REMOVE}
}
func ListReplace() Instruction {
	return Instruction{Kind: OP_LIST_REPLACE}
}
