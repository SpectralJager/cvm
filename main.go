package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime/pprof"
)

func main() {
	instrs := []Instruction{
		F32Load(12.2),
		F32ToI32(),
		I32ToF32(),
		Halt(),
	}
	for i, instr := range instrs {
		fmt.Printf("%04d: %s\n", i, instr.String())
	}
	fmt.Println()
	vm := CVM{}
	fl, err := os.Create("fib.prof")
	if err != nil {
		panic(err)
	}
	pprof.StartCPUProfile(fl)
	err = vm.Execute(context.TODO(), instrs)
	if err != nil {
		if !errors.Is(err, ErrReachHalt) {
			panic(err)
		}
	}
	pprof.StopCPUProfile()
	fmt.Println(vm.Trace())
}
