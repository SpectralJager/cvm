package cvm

import (
	"bytes"
	"context"
	i "cvm/instruction"
	"cvm/object"
	"testing"
)

func obj(obj object.CVMObject, err error) object.CVMObject {
	return obj
}

func TestI32(t *testing.T) {
	testCases := []struct {
		desc   string
		instrs []i.Instruction
		result object.CVMObject
	}{
		{
			desc: "test i32 addition #1",
			instrs: []i.Instruction{
				i.I32Load(10),
				i.I32Load(20),
				i.I32Add(),
			},
			result: obj(object.CreateI32(30)),
		},
		{
			desc: "test i32 addition #2",
			instrs: []i.Instruction{
				i.I32Load(10),
				i.I32Load(-20),
				i.I32Add(),
			},
			result: obj(object.CreateI32(-10)),
		},
		{
			desc: "test i32 subtraction #1",
			instrs: []i.Instruction{
				i.I32Load(10),
				i.I32Load(20),
				i.I32Sub(),
			},
			result: obj(object.CreateI32(-10)),
		},
		{
			desc: "test i32 subtraction #2",
			instrs: []i.Instruction{
				i.I32Load(10),
				i.I32Load(-20),
				i.I32Sub(),
			},
			result: obj(object.CreateI32(30)),
		},
		{
			desc: "test i32 multiplication #1",
			instrs: []i.Instruction{
				i.I32Load(10),
				i.I32Load(20),
				i.I32Mul(),
			},
			result: obj(object.CreateI32(200)),
		},
		{
			desc: "test i32 multiplication #2",
			instrs: []i.Instruction{
				i.I32Load(10),
				i.I32Load(-20),
				i.I32Mul(),
			},
			result: obj(object.CreateI32(-200)),
		},
		{
			desc: "test i32 division #1",
			instrs: []i.Instruction{
				i.I32Load(200),
				i.I32Load(10),
				i.I32Div(),
			},
			result: obj(object.CreateI32(20)),
		},
		{
			desc: "test i32 division #2",
			instrs: []i.Instruction{
				i.I32Load(200),
				i.I32Load(-10),
				i.I32Div(),
			},
			result: obj(object.CreateI32(-20)),
		},
		{
			desc: "test i32 lt #1",
			instrs: []i.Instruction{
				i.I32Load(10),
				i.I32Load(20),
				i.I32Lt(),
			},
			result: obj(object.CreateBool(true)),
		},
		{
			desc: "test i32 lt #2",
			instrs: []i.Instruction{
				i.I32Load(20),
				i.I32Load(10),
				i.I32Lt(),
			},
			result: obj(object.CreateBool(false)),
		},
		{
			desc: "test i32 gt #1",
			instrs: []i.Instruction{
				i.I32Load(10),
				i.I32Load(20),
				i.I32Gt(),
			},
			result: obj(object.CreateBool(false)),
		},
		{
			desc: "test i32 gt #2",
			instrs: []i.Instruction{
				i.I32Load(20),
				i.I32Load(10),
				i.I32Gt(),
			},
			result: obj(object.CreateBool(true)),
		},
		{
			desc: "test i32 leq #1",
			instrs: []i.Instruction{
				i.I32Load(10),
				i.I32Load(20),
				i.I32Leq(),
			},
			result: obj(object.CreateBool(true)),
		},
		{
			desc: "test i32 leq #2",
			instrs: []i.Instruction{
				i.I32Load(20),
				i.I32Load(10),
				i.I32Leq(),
			},
			result: obj(object.CreateBool(false)),
		},
		{
			desc: "test i32 leq #3",
			instrs: []i.Instruction{
				i.I32Load(20),
				i.I32Load(20),
				i.I32Leq(),
			},
			result: obj(object.CreateBool(true)),
		},
		{
			desc: "test i32 geq #1",
			instrs: []i.Instruction{
				i.I32Load(10),
				i.I32Load(20),
				i.I32Geq(),
			},
			result: obj(object.CreateBool(false)),
		},
		{
			desc: "test i32 geq #2",
			instrs: []i.Instruction{
				i.I32Load(20),
				i.I32Load(10),
				i.I32Geq(),
			},
			result: obj(object.CreateBool(true)),
		},
		{
			desc: "test i32 geq #3",
			instrs: []i.Instruction{
				i.I32Load(20),
				i.I32Load(20),
				i.I32Geq(),
			},
			result: obj(object.CreateBool(true)),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			vm := CVM{}
			err := vm.Execute(context.TODO(), tC.instrs)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(object.Bytes(vm.Stack[0]), object.Bytes(tC.result)) {
				t.Fatalf("%v != %v", vm.Stack[0], tC.result)
			}
		})
	}
}

func TestF32(t *testing.T) {
	testCases := []struct {
		desc   string
		instrs []i.Instruction
		result object.CVMObject
	}{
		{
			desc: "test f32 addition #1",
			instrs: []i.Instruction{
				i.F32Load(10),
				i.F32Load(20),
				i.F32Add(),
			},
			result: obj(object.CreateF32(30)),
		},
		{
			desc: "test f32 addition #2",
			instrs: []i.Instruction{
				i.F32Load(10),
				i.F32Load(-20),
				i.F32Add(),
			},
			result: obj(object.CreateF32(-10)),
		},
		{
			desc: "test f32 subtraction #1",
			instrs: []i.Instruction{
				i.F32Load(10),
				i.F32Load(20),
				i.F32Sub(),
			},
			result: obj(object.CreateF32(-10)),
		},
		{
			desc: "test f32 subtraction #2",
			instrs: []i.Instruction{
				i.F32Load(10),
				i.F32Load(-20),
				i.F32Sub(),
			},
			result: obj(object.CreateF32(30)),
		},
		{
			desc: "test f32 multiplication #1",
			instrs: []i.Instruction{
				i.F32Load(10),
				i.F32Load(20),
				i.F32Mul(),
			},
			result: obj(object.CreateF32(200)),
		},
		{
			desc: "test f32 multiplication #2",
			instrs: []i.Instruction{
				i.F32Load(10),
				i.F32Load(-20),
				i.F32Mul(),
			},
			result: obj(object.CreateF32(-200)),
		},
		{
			desc: "test f32 division #1",
			instrs: []i.Instruction{
				i.F32Load(200),
				i.F32Load(10),
				i.F32Div(),
			},
			result: obj(object.CreateF32(20)),
		},
		{
			desc: "test f32 division #2",
			instrs: []i.Instruction{
				i.F32Load(200),
				i.F32Load(-10),
				i.F32Div(),
			},
			result: obj(object.CreateF32(-20)),
		},
		{
			desc: "test f32 lt #1",
			instrs: []i.Instruction{
				i.F32Load(10),
				i.F32Load(20),
				i.F32Lt(),
			},
			result: obj(object.CreateBool(true)),
		},
		{
			desc: "test f32 lt #2",
			instrs: []i.Instruction{
				i.F32Load(20),
				i.F32Load(10),
				i.F32Lt(),
			},
			result: obj(object.CreateBool(false)),
		},
		{
			desc: "test f32 gt #1",
			instrs: []i.Instruction{
				i.F32Load(10),
				i.F32Load(20),
				i.F32Gt(),
			},
			result: obj(object.CreateBool(false)),
		},
		{
			desc: "test f32 gt #2",
			instrs: []i.Instruction{
				i.F32Load(20),
				i.F32Load(10),
				i.F32Gt(),
			},
			result: obj(object.CreateBool(true)),
		},
		{
			desc: "test f32 leq #1",
			instrs: []i.Instruction{
				i.F32Load(10),
				i.F32Load(20),
				i.F32Leq(),
			},
			result: obj(object.CreateBool(true)),
		},
		{
			desc: "test f32 leq #2",
			instrs: []i.Instruction{
				i.F32Load(20),
				i.F32Load(10),
				i.F32Leq(),
			},
			result: obj(object.CreateBool(false)),
		},
		{
			desc: "test f32 leq #3",
			instrs: []i.Instruction{
				i.F32Load(20),
				i.F32Load(20),
				i.F32Leq(),
			},
			result: obj(object.CreateBool(true)),
		},
		{
			desc: "test f32 geq #1",
			instrs: []i.Instruction{
				i.F32Load(10),
				i.F32Load(20),
				i.F32Geq(),
			},
			result: obj(object.CreateBool(false)),
		},
		{
			desc: "test f32 geq #2",
			instrs: []i.Instruction{
				i.F32Load(20),
				i.F32Load(10),
				i.F32Geq(),
			},
			result: obj(object.CreateBool(true)),
		},
		{
			desc: "test f32 geq #3",
			instrs: []i.Instruction{
				i.F32Load(20),
				i.F32Load(20),
				i.F32Geq(),
			},
			result: obj(object.CreateBool(true)),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			vm := CVM{}
			err := vm.Execute(context.TODO(), tC.instrs)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(object.Bytes(vm.Stack[0]), object.Bytes(tC.result)) {
				t.Fatalf("%v != %v", vm.Stack[0], tC.result)
			}
		})
	}
}
