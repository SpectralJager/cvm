package instruction

import "cvm/object"

func StructNew(tags ...byte) Instruction {
	lO, err := object.CreateI32(int32(len(tags)))
	if err != nil {
		panic(err)
	}
	buf := []byte{object.TAG_STRUCT}
	buf = append(buf, object.Bytes(lO)...)
	buf = append(buf, tags...)
	for _, tag := range tags {
		obj, err := object.CreateDefault(tag)
		if err != nil {
			panic(err)
		}
		buf = append(buf, object.Bytes(obj)...)
	}
	return Instruction{Kind: OP_STRUCT_NEW, Operands: buf}
}

func StructGet() Instruction {
	return Instruction{Kind: OP_STRUCT_GET}
}

func StructSet() Instruction {
	return Instruction{Kind: OP_STRUCT_SET}
}
