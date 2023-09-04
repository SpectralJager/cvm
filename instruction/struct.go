package instruction

import "cvm/object"

func StructLoad(tags ...byte) Instruction {
	lO, err := object.CreateI32(int32(len(tags)))
	if err != nil {
		panic(err)
	}
	buf := []byte{object.TAG_STRUCT}
	buf = append(buf, object.Bytes(lO)...)
	buf = append(buf, tags...)
	return Instruction{Kind: OP_STRUCT_LOAD, Operands: buf}
}
