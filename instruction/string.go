package instruction

import "cvm/object"

func StringLoad(str string) Instruction {
	obj, err := object.CreateString(str)
	if err != nil {
		panic(err)
	}
	return Instruction{Kind: OP_STRING_LOAD, Operands: object.Bytes(obj)}
}

func StringConcat() Instruction {
	return Instruction{Kind: OP_STRING_CONCAT}
}

func StringSplit() Instruction {
	return Instruction{Kind: OP_STRING_SPLIT}
}

func StringFormat() Instruction {
	return Instruction{Kind: OP_STRING_FORMAT}
}

func StringPrint() Instruction {
	return Instruction{Kind: OP_STRING_PRINT}
}

func StringPrintf() Instruction {
	return Instruction{Kind: OP_STRING_PRINTF}
}

func StringPrintln() Instruction {
	return Instruction{Kind: OP_STRING_PRINTLN}
}
