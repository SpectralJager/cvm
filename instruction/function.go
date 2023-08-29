package instruction

import (
	"cvm/object"
	"encoding/binary"
)

func FuncCall(addr uint32, args uint32) Instruction {
	buf := make([]byte, 0, 10)
	buf = append(buf, object.TAG_I32)
	buf = binary.LittleEndian.AppendUint32(buf, addr)
	buf = append(buf, object.TAG_I32)
	buf = binary.LittleEndian.AppendUint32(buf, args)
	return Instruction{Kind: OP_FUNC_CALL, Operands: buf}
}
func FuncRet(x uint32) Instruction {
	buf := make([]byte, 0, 5)
	buf = append(buf, object.TAG_I32)
	buf = binary.LittleEndian.AppendUint32(buf, x)
	return Instruction{Kind: OP_FUNC_RET, Operands: buf}
}
