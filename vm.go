package main

import (
	"bytes"
	"context"
	"fmt"
)

type CVM struct {
	Stack      []CVMObject
	C          bool
	StackFrame []Frame
}

func (vm *CVM) Push(ctx context.Context, obj CVMObject) error {
	vm.Stack = append(vm.Stack, obj)
	return nil
}
func (vm *CVM) Pop(ctx context.Context) (CVMObject, error) {
	if len(vm.Stack) == 0 {
		return CVMObject{}, fmt.Errorf("cant pop, stack is empty")
	}
	obj := vm.Stack[len(vm.Stack)-1]
	vm.Stack = vm.Stack[:len(vm.Stack)-1]
	return obj, nil
}
func (vm *CVM) StackTrace() string {
	var buf bytes.Buffer
	fmt.Fprint(&buf, "=== Stack:\n")
	for i, obj := range vm.Stack {
		fmt.Fprintf(&buf, "\t$%03d -> %s\n", i, obj.String())
	}
	return buf.String()
}

func (vm *CVM) Execute(ctx context.Context, instrs []Instruction) error {
	for ip := uint32(0); ip < uint32(len(instrs)); {
		instr := instrs[ip]
		switch instr.Kind {
		case OP_NULL:
			ip++
		case OP_HALT:
			return ErrReachHalt
		case OP_I32_LOAD:
			ip++
			obj, err := CreateI32Object(instr.Operands)
			if err != nil {
				return err
			}
			vm.Push(ctx, obj)
		case OP_I32_ADD:
			ip++
			obj2, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			obj1, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			resVal := obj1.ToI32() + obj2.ToI32()
			resObj, err := CreateI32Object[int32](resVal)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case OP_I32_SUB:
			ip++
			obj2, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			obj1, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			resVal := obj1.ToI32() - obj2.ToI32()
			resObj, err := CreateI32Object[int32](resVal)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case OP_I32_MUL:
			ip++
			obj2, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			obj1, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			resVal := obj1.ToI32() * obj2.ToI32()
			resObj, err := CreateI32Object[int32](resVal)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case OP_I32_DIV:
			ip++
			obj2, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			obj1, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			if obj2.ToI32() == 0 {
				return fmt.Errorf("division by zero")
			}
			resVal := obj1.ToI32() / obj2.ToI32()
			resObj, err := CreateI32Object[int32](resVal)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case OP_I32_LT:
			ip++
			obj2, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			obj1, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			if obj2.ToI32() == 0 {
				return fmt.Errorf("division by zero")
			}
			res := obj1.ToI32() < obj2.ToI32()
			resObj, err := CreateBool[bool](res)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case OP_I32_GT:
			ip++
			obj2, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			obj1, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			if obj2.ToI32() == 0 {
				return fmt.Errorf("division by zero")
			}
			res := obj1.ToI32() > obj2.ToI32()
			resObj, err := CreateBool[bool](res)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case OP_I32_LEQ:
			ip++
			obj2, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			obj1, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			if obj2.ToI32() == 0 {
				return fmt.Errorf("division by zero")
			}
			res := obj1.ToI32() <= obj2.ToI32()
			resObj, err := CreateBool[bool](res)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case OP_I32_GEQ:
			ip++
			obj2, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			obj1, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			if obj2.ToI32() == 0 {
				return fmt.Errorf("division by zero")
			}
			res := obj1.ToI32() >= obj2.ToI32()
			resObj, err := CreateBool[bool](res)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case OP_I32_EQ:
			ip++
			obj2, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			obj1, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			if obj2.ToI32() == 0 {
				return fmt.Errorf("division by zero")
			}
			res := obj1.ToI32() == obj2.ToI32()
			resObj, err := CreateBool[bool](res)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case OP_JUMP:
			obj, err := CreateI32Object(instr.Operands)
			if err != nil {
				return err
			}
			ip = uint32(obj.ToI32())
		case OP_JUMPC:
			obj, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			val := obj.ToBool()
			if !val {
				ip++
				continue
			}
			addr, err := CreateI32Object(instr.Operands)
			if err != nil {
				return err
			}
			ip = uint32(addr.ToI32())
		case OP_JUMPNC:
			obj, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			val := obj.ToBool()
			if val {
				ip++
				continue
			}
			addr, err := CreateI32Object(instr.Operands)
			if err != nil {
				return err
			}
			ip = uint32(addr.ToI32())
		default:
			return fmt.Errorf("unknown instruction of kind 0x%02x", instr.Kind)
		}
	}
	return nil
}
