package main

import (
	"context"
	"cvm"
	i "cvm/instruction"
	"cvm/object"
	"fmt"
)

func main() {
	instrs := []i.Instruction{
		i.ListNew(object.TAG_I32),
		i.New(),
		i.Load(0),
		i.Load(0),
		i.ListLength(),
		i.I32Load(6),
		i.ListInsert(),
		i.Save(0),
		i.Load(0),
		i.I32Load(0),
		i.I32Load(3),
		i.ListInsert(),
		i.Save(0),
		i.Load(0),
		i.I32Load(0),
		i.I32Load(2),
		i.ListReplace(),
		i.I32Load(1),
		i.I32Load(4),
		i.ListReplace(),
		i.New(),
		i.Load(1),
		i.I32Load(0),
		i.ListGet(),
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
