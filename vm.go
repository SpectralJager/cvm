package main

import (
	"bytes"
	"context"
	"fmt"
)

type CVM struct {
	Stack      []CVMObject
	Heap       []CVMObject
	StackFrame []Frame
}

func (vm *CVM) New(ctx context.Context, obj CVMObject) (uint32, error) {
	vm.Heap = append(vm.Heap, obj)
	return uint32(len(vm.Heap) - 1), nil
}
func (vm *CVM) Load(ctx context.Context, ind uint32) (CVMObject, error) {
	if uint32(len(vm.Heap)) <= ind {
		return CVMObject{}, fmt.Errorf("symbol with index %d not found", ind)
	}
	return vm.Heap[ind], nil
}
func (vm *CVM) Save(ctx context.Context, ind uint32, obj CVMObject) error {
	if uint32(len(vm.Heap)) <= ind {
		return fmt.Errorf("symbol with index %d not found", ind)
	}
	if vm.Heap[ind].Tag != obj.Tag {
		return fmt.Errorf("unexpected tag %d, want %d", obj.Tag, vm.Heap[ind].Tag)
	}
	vm.Heap[ind] = obj
	return nil
}

func (vm *CVM) LastFuncFrame(ctx context.Context) (Frame, error) {
	var fr Frame
	for i := len(vm.StackFrame) - 1; i >= 0; i-- {
		fr = vm.StackFrame[i]
		if fr.FrameOffset != -1 {
			break
		}
	}
	if fr.FrameOffset == -1 {
		return fr, fmt.Errorf("cant find function frame")
	}
	return fr, nil
}
func (vm *CVM) LastFrame(ctx context.Context) (Frame, error) {
	if len(vm.StackFrame) == 0 {
		return Frame{}, fmt.Errorf("empty StackFrame")
	}
	return vm.StackFrame[len(vm.StackFrame)-1], nil
}
func (vm *CVM) PushFrame(ctx context.Context, fr Frame) error {
	vm.StackFrame = append(vm.StackFrame, fr)
	return nil
}
func (vm *CVM) PopFrame(ctx context.Context) (Frame, error) {
	if len(vm.StackFrame) == 0 {
		return Frame{}, fmt.Errorf("cant pop, stackFrame is empty")
	}
	fr := vm.StackFrame[len(vm.StackFrame)-1]
	vm.StackFrame = vm.StackFrame[:len(vm.StackFrame)-1]
	return fr, nil
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
func (vm *CVM) Trace() string {
	var buf bytes.Buffer
	fmt.Fprint(&buf, "=== Heap:\n")
	for i, obj := range vm.Heap {
		fmt.Fprintf(&buf, "\t$%03d -> %s\n", i, obj.String())
	}
	fmt.Fprint(&buf, "=== StackFrace:\n")
	for i, fr := range vm.StackFrame {
		fmt.Fprintf(&buf, "\t$%03d -> %s\n", i, fr.String())
	}
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
		case OP_BLOCK_START:
			ip++
			addr, err := CreateI32Object(instr.Operands)
			if err != nil {
				return err
			}
			vm.PushFrame(ctx, Frame{
				StackOffset: len(vm.Stack),
				HeapOffset:  len(vm.Heap),
				ReturnIP:    uint32(addr.ToI32()),
				FrameOffset: -1,
			})
		case OP_BLOCK_BR:
			fr, err := vm.LastFrame(ctx)
			if err != nil {
				return err
			}
			ip = fr.ReturnIP
		case OP_BLOCK_END:
			ip++
			fr, err := vm.PopFrame(ctx)
			if err != nil {
				return err
			}
			vm.Heap = vm.Heap[:fr.HeapOffset]
			vm.Stack = vm.Stack[:fr.StackOffset]
		case OP_BLOCK_LOAD:
			ip++
			ind, err := CreateI32Object(instr.Operands)
			if err != nil {
				return err
			}
			fr, err := vm.LastFrame(ctx)
			if err != nil {
				return err
			}
			obj, err := vm.Load(ctx, uint32(int32(fr.HeapOffset)+ind.ToI32()))
			if err != nil {
				return err
			}
			vm.Push(ctx, obj)
		case OP_BLOCK_SAVE:
			ip++
			ind, err := CreateI32Object(instr.Operands)
			if err != nil {
				return err
			}
			obj, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			fr, err := vm.LastFrame(ctx)
			if err != nil {
				return err
			}
			err = vm.Save(ctx, uint32(int32(fr.HeapOffset)+ind.ToI32()), obj)
			if err != nil {
				return err
			}
		case OP_LOAD:
			ip++
			ind, err := CreateI32Object(instr.Operands)
			if err != nil {
				return err
			}
			obj, err := vm.Load(ctx, uint32(ind.ToI32()))
			vm.Push(ctx, obj)
		case OP_SAVE:
			ip++
			ind, err := CreateI32Object(instr.Operands)
			if err != nil {
				return err
			}
			obj, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			err = vm.Save(ctx, uint32(ind.ToI32()), obj)
			if err != nil {
				return err
			}
		case OP_NEW:
			ip++
			obj, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			_, err = vm.New(ctx, obj)
			if err != nil {
				return err
			}
		case OP_FUNC_CALL:
			ip++
			addr, err := CreateI32Object(instr.Operands[:5])
			if err != nil {
				return err
			}
			argsLen, err := CreateI32Object(instr.Operands[5:10])
			if err != nil {
				return err
			}
			vm.PushFrame(ctx, Frame{
				StackOffset: len(vm.Stack) - int(argsLen.ToI32()),
				HeapOffset:  len(vm.Heap),
				ReturnIP:    uint32(ip),
				FrameOffset: len(vm.StackFrame),
			})
			ip = uint32(addr.ToI32())
		case OP_FUNC_RET:
			fr, err := vm.LastFuncFrame(ctx)
			if err != nil {
				return err
			}
			retLen, err := CreateI32Object(instr.Operands[:5])
			if err != nil {
				return err
			}
			ip = fr.ReturnIP
			objs := vm.Stack[len(vm.Stack)-int(retLen.ToI32()):]
			vm.Heap = vm.Heap[:fr.HeapOffset]
			vm.Stack = vm.Stack[:fr.StackOffset]
			vm.StackFrame = vm.StackFrame[:fr.FrameOffset]
			for _, obj := range objs {
				vm.Push(ctx, obj)
			}
		case OP_LOCAL_LOAD:
			ip++
			ind, err := CreateI32Object(instr.Operands)
			if err != nil {
				return err
			}
			fr, err := vm.LastFuncFrame(ctx)
			if err != nil {
				return err
			}
			obj, err := vm.Load(ctx, uint32(int32(fr.HeapOffset)+ind.ToI32()))
			if err != nil {
				return err
			}
			vm.Push(ctx, obj)
		case OP_LOCAL_SAVE:
			ip++
			ind, err := CreateI32Object(instr.Operands)
			if err != nil {
				return err
			}
			obj, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			fr, err := vm.LastFuncFrame(ctx)
			if err != nil {
				return err
			}
			err = vm.Save(ctx, uint32(int32(fr.HeapOffset)+ind.ToI32()), obj)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown instruction of kind 0x%02x", instr.Kind)
		}
	}
	return nil
}
