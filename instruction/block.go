package instruction

import (
	"cvm/object"
	"encoding/binary"
)

func BlockStart(x uint32) Instruction {
	buf := make([]byte, 0, 5)
	buf = append(buf, object.TAG_I32)
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
	buf = append(buf, object.TAG_I32)
	buf = binary.LittleEndian.AppendUint32(buf, uint32(x))
	return Instruction{Kind: OP_BLOCK_LOAD, Operands: buf}
}
func BlockSave(x uint32) Instruction {
	buf := make([]byte, 0, 5)
	buf = append(buf, object.TAG_I32)
	buf = binary.LittleEndian.AppendUint32(buf, uint32(x))
	return Instruction{Kind: OP_BLOCK_SAVE, Operands: buf}
}
