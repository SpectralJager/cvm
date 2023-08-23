package main

import (
	"context"
	"errors"
	"fmt"
)

type Frame struct {
	StackOffset int
	HeapOffset  int
	FrameOffset int
	ReturnIP    uint32
}

func (f *Frame) String() string {
	return fmt.Sprintf("#stack: %d, #heap: %d,#return: %d", f.StackOffset, f.HeapOffset, f.ReturnIP)
}

func main() {
	instrs := []Instruction{
		I32Load(47),
		FuncCall(4, 1),
		New(),
		Halt(),
		New(), // fn fib
		BlockStart(13),
		LocalLoad(0),
		I32Load(2),
		I32Lt(),
		JumpNC(13),
		LocalLoad(0),
		FuncRet(1),
		BlockBr(),
		BlockEnd(),
		LocalLoad(0),
		I32Load(1),
		I32Sub(),
		FuncCall(4, 1),
		LocalLoad(0),
		I32Load(2),
		I32Sub(),
		FuncCall(4, 1),
		I32Add(),
		FuncRet(1),
	}
	for i, instr := range instrs {
		fmt.Printf("%04d: %s\n", i, instr.String())
	}
	fmt.Println()
	vm := CVM{
		Stack:      make([]CVMObject, 0, 1024),
		Heap:       make([]CVMObject, 0, 2048),
		StackFrame: make([]Frame, 0, 256),
	}
	// fl, err := os.Create("fib.prof")
	// if err != nil {
	// 	panic(err)
	// }
	// pprof.StartCPUProfile(fl)
	err := vm.Execute(context.TODO(), instrs)
	if err != nil {
		if !errors.Is(err, ErrReachHalt) {
			panic(err)
		}
	}
	// pprof.StopCPUProfile()
	fmt.Println(vm.Trace())
}
