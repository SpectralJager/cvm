package main

import (
	"context"
	"cvm"
	i "cvm/instruction"
	"fmt"
)

func main() {
	instrs := []i.Instruction{
		i.Halt(),
	}
	for i, inst := range instrs {
		fmt.Printf("%04d: %s\n", i, inst.String())
	}
	fmt.Println()
	vm := cvm.CVM{}
	// fl, err := os.Create("prof.prof")
	// if err != nil {
	// 	panic(err)
	// }
	// pprof.StartCPUProfile(fl)
	err := vm.Execute(context.TODO(), instrs)
	if err != nil {
		panic(err)
	}
	// pprof.StopCPUProfile()
	fmt.Println(vm.Trace())
}
