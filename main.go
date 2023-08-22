package main

import (
	"context"
	"errors"
	"fmt"
)

type Frame struct {
	Start    int
	ReturnIP uint32
}

func main() {
	instrs := []Instruction{
		Halt(),
	}
	for i, instr := range instrs {
		fmt.Printf("%d: %s\n", i, instr.String())
	}
	fmt.Println()
	vm := CVM{Stack: make([]CVMObject, 0, 1024)}
	err := vm.Execute(context.TODO(), instrs)
	if err != nil {
		if !errors.Is(err, ErrReachHalt) {
			panic(err)
		}
	}
	fmt.Println(vm.StackTrace())
}
