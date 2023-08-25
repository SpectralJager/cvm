package main

import (
	"context"
	"errors"
	"fmt"
)

func main() {
	instrs := []Instruction{
		ListNew(TAG_I32),
		I32Load(12),
		ListAppend(),
		I32Load(24),
		ListAppend(),
		I32Load(32),
		ListAppend(),
		I32Load(48),
		ListAppend(),
		I32Load(100),
		ListAppend(),
		New(),
		Load(0),
		ListPop(),
		Pop(),
		Save(0),
		Halt(),
	}
	for i, instr := range instrs {
		fmt.Printf("%04d: %s\n", i, instr.String())
	}
	fmt.Println()
	vm := CVM{}
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
