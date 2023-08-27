package main

import (
	"bytes"
	"context"
	"fmt"
)

const (
	STACK_SIZE       = 2048
	HEAP_SIZE        = 2048
	STACK_FRAME_SIZE = 2048
)

type CVM struct {
	Stack      [STACK_SIZE]CVMObject
	Heap       [HEAP_SIZE]CVMObject
	StackFrame [STACK_FRAME_SIZE]Frame
	SP, HP, FP uint
}

func (vm *CVM) New(ctx context.Context, obj CVMObject) error {
	if vm.HP >= HEAP_SIZE {
		return fmt.Errorf("heap overflow")
	}
	vm.Heap[vm.HP] = obj
	vm.HP++
	return nil
}
func (vm *CVM) Load(ctx context.Context, ind uint32) (CVMObject, error) {
	if uint32(vm.HP) < ind {
		return CVMObject{}, fmt.Errorf("symbol with index %d not found", ind)
	}
	return vm.Heap[ind], nil
}
func (vm *CVM) Free(ctx context.Context, ind uint32) error {
	if uint32(vm.HP) < ind {
		return fmt.Errorf("symbol with index %d not found", ind)
	}
	vm.Heap[ind] = CVMObject{}
	return nil
}
func (vm *CVM) Save(ctx context.Context, ind uint32, obj CVMObject) error {
	if uint32(vm.HP) < ind {
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
	for i := vm.FP - 1; i >= 0; i-- {
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
	if vm.FP == 0 {
		return Frame{}, fmt.Errorf("empty StackFrame")
	}
	return vm.StackFrame[len(vm.StackFrame)-1], nil
}
func (vm *CVM) PushFrame(ctx context.Context, fr Frame) error {
	if vm.FP >= STACK_FRAME_SIZE {
		return fmt.Errorf("stack frame overflow")
	}
	vm.StackFrame[vm.FP] = fr
	vm.FP++
	return nil
}
func (vm *CVM) PopFrame(ctx context.Context) (Frame, error) {
	if vm.FP == 0 {
		return Frame{}, fmt.Errorf("cant pop, stackFrame is empty")
	}
	vm.FP--
	fr := vm.StackFrame[vm.FP]
	return fr, nil
}
func (vm *CVM) Push(ctx context.Context, obj CVMObject) error {
	if vm.SP >= STACK_SIZE {
		return fmt.Errorf("stack overflow")
	}
	vm.Stack[vm.SP] = obj
	vm.SP++
	return nil
}
func (vm *CVM) Pop(ctx context.Context) (CVMObject, error) {
	if vm.SP == 0 {
		return CVMObject{}, fmt.Errorf("stack is empty")
	}
	vm.SP--
	obj := vm.Stack[vm.SP]
	return obj, nil
}
func (vm *CVM) Trace() string {
	var buf bytes.Buffer
	fmt.Fprint(&buf, "=== Heap:\n")
	for i := 0; i < int(vm.HP); i++ {
		if vm.Heap[i].Value != nil {
			fmt.Fprintf(&buf, "\t$%03d -> %s\n", i, vm.Heap[i].String())
		}
	}
	fmt.Fprint(&buf, "=== StackFrace:\n")
	for i := 0; i < int(vm.FP); i++ {
		fmt.Fprintf(&buf, "\t$%03d -> %s\n", i, vm.StackFrame[i].String())
	}
	fmt.Fprint(&buf, "=== Stack:\n")
	for i := 0; i < int(vm.SP); i++ {
		fmt.Fprintf(&buf, "\t$%03d -> %s\n", i, vm.Stack[i].String())
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
			obj, err := CreateObject(instr.Operands)
			if err != nil {
				return err
			}
			vm.Push(ctx, obj)
		case OP_I32_NEG:
			ip++
			obj, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			resVal := -obj.ToI32()
			resObj, err := CreateI32(resVal)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
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
			resObj, err := CreateI32(resVal)
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
			resObj, err := CreateI32(resVal)
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
			resObj, err := CreateI32(resVal)
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
			resObj, err := CreateI32(resVal)
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
			resObj, err := CreateBool(res)
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
			resObj, err := CreateBool(res)
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
			resObj, err := CreateBool(res)
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
			resObj, err := CreateBool(res)
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
			resObj, err := CreateBool(res)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case OP_I32_TO_F32:
			ip++
			obj, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			resVal := obj.I32ToF32()
			resObj, err := CreateF32(resVal)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case OP_I32_TO_BOOL:
			ip++
			obj, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			resVal := obj.I32ToBool()
			resObj, err := CreateBool(resVal)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case OP_BOOL_LOAD:
			ip++
			obj, err := CreateObject(instr.Operands)
			if err != nil {
				return err
			}
			vm.Push(ctx, obj)
		case OP_BOOL_NOT:
			ip++
			obj, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			resVal := !obj.ToBool()
			resObj, err := CreateBool(resVal)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case OP_BOOL_AND:
			ip++
			obj2, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			obj1, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			resVal := obj1.ToBool() && obj2.ToBool()
			resObj, err := CreateBool(resVal)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case OP_BOOL_OR:
			ip++
			obj2, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			obj1, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			resVal := obj1.ToBool() || obj2.ToBool()
			resObj, err := CreateBool(resVal)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case OP_BOOL_XOR:
			ip++
			obj2, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			obj1, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			resVal := false
			b1 := obj1.ToBool()
			b2 := obj2.ToBool()
			if b1 != b2 && (b1 == true || b2 == true) {
				resVal = true
			}
			resObj, err := CreateBool(resVal)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case OP_F32_LOAD:
			ip++
			obj, err := CreateObject(instr.Operands)
			if err != nil {
				return err
			}
			vm.Push(ctx, obj)
		case OP_F32_NEG:
			ip++
			obj, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			resVal := -obj.ToF32()
			resObj, err := CreateF32(resVal)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case OP_F32_ADD:
			ip++
			obj2, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			obj1, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			resVal := obj1.ToF32() + obj2.ToF32()
			resObj, err := CreateF32(resVal)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case OP_F32_SUB:
			ip++
			obj2, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			obj1, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			resVal := obj1.ToF32() - obj2.ToF32()
			resObj, err := CreateF32(resVal)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case OP_F32_MUL:
			ip++
			obj2, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			obj1, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			resVal := obj1.ToF32() * obj2.ToF32()
			resObj, err := CreateF32(resVal)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case OP_F32_DIV:
			ip++
			obj2, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			obj1, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			if obj2.ToF32() == 0.0 {
				return fmt.Errorf("division by zero")
			}
			resVal := obj1.ToF32() / obj2.ToF32()
			resObj, err := CreateF32(resVal)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case OP_F32_LT:
			ip++
			obj2, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			obj1, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			res := obj1.ToF32() < obj2.ToF32()
			resObj, err := CreateBool(res)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case OP_F32_GT:
			ip++
			obj2, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			obj1, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			res := obj1.ToF32() > obj2.ToF32()
			resObj, err := CreateBool(res)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case OP_F32_LEQ:
			ip++
			obj2, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			obj1, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			res := obj1.ToF32() <= obj2.ToF32()
			resObj, err := CreateBool(res)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case OP_F32_GEQ:
			ip++
			obj2, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			obj1, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			res := obj1.ToF32() >= obj2.ToF32()
			resObj, err := CreateBool(res)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case OP_F32_EQ:
			ip++
			obj2, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			obj1, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			res := obj1.ToF32() == obj2.ToF32()
			resObj, err := CreateBool(res)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case OP_F32_TO_I32:
			ip++
			obj, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			resVal := obj.F32ToI32()
			resObj, err := CreateI32(resVal)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case OP_F32_TO_BOOL:
			ip++
			obj, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			resVal := obj.F32ToBool()
			resObj, err := CreateBool(resVal)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case OP_JUMP:
			obj, err := CreateObject(instr.Operands)
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
			addr, err := CreateObject(instr.Operands)
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
			addr, err := CreateObject(instr.Operands)
			if err != nil {
				return err
			}
			ip = uint32(addr.ToI32())
		case OP_BLOCK_START:
			ip++
			addr, err := CreateObject(instr.Operands)
			if err != nil {
				return err
			}
			vm.PushFrame(ctx, Frame{
				StackOffset: int(vm.SP),
				HeapOffset:  int(vm.HP),
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
			vm.HP = uint(fr.HeapOffset)
			vm.SP = uint(fr.StackOffset)
		case OP_BLOCK_LOAD:
			ip++
			ind, err := CreateObject(instr.Operands)
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
			ind, err := CreateObject(instr.Operands)
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
			ind, err := CreateObject(instr.Operands)
			if err != nil {
				return err
			}
			obj, err := vm.Load(ctx, uint32(ind.ToI32()))
			vm.Push(ctx, obj)
		case OP_SAVE:
			ip++
			ind, err := CreateObject(instr.Operands)
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
		case OP_FREE:
			ip++
			ind, err := CreateObject(instr.Operands)
			if err != nil {
				return err
			}
			err = vm.Free(ctx, uint32(ind.ToI32()))
			if err != nil {
				return err
			}
		case OP_NEW:
			ip++
			obj, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			err = vm.New(ctx, obj)
			if err != nil {
				return err
			}
		case OP_POP:
			ip++
			_, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
		case OP_FUNC_CALL:
			ip++
			addr, err := CreateObject(instr.Operands[:5])
			if err != nil {
				return err
			}
			argsLen, err := CreateObject(instr.Operands[5:10])
			if err != nil {
				return err
			}
			vm.PushFrame(ctx, Frame{
				StackOffset: int(vm.SP) - int(argsLen.ToI32()),
				HeapOffset:  int(vm.HP),
				ReturnIP:    uint32(ip),
				FrameOffset: int(vm.FP),
			})
			ip = uint32(addr.ToI32())
		case OP_FUNC_RET:
			fr, err := vm.LastFuncFrame(ctx)
			if err != nil {
				return err
			}
			retLen, err := CreateObject(instr.Operands[:5])
			if err != nil {
				return err
			}
			ip = fr.ReturnIP
			objs := vm.Stack[int(vm.SP)-int(retLen.ToI32()) : vm.SP]
			vm.HP = uint(fr.HeapOffset)
			vm.SP = uint(fr.StackOffset)
			vm.FP = uint(fr.FrameOffset)
			for _, obj := range objs {
				vm.Push(ctx, obj)
			}
		case OP_LOCAL_LOAD:
			ip++
			ind, err := CreateObject(instr.Operands)
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
			ind, err := CreateObject(instr.Operands)
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
		case OP_LIST_NEW:
			ip++
			obj, err := CreateList(instr.Operands)
			if err != nil {
				return err
			}
			vm.Push(ctx, obj)
		case OP_LIST_APPEND:
			ip++
			obj, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			list, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			if list.Tag != TAG_LIST {
				return fmt.Errorf("expected list, got %v", list.Tag)
			}
			if list.Value[0] != obj.Tag {
				return fmt.Errorf("cant add %d to list of %d", obj.Tag, list.Value[0])
			}
			len, err := CreateObject(list.Value[1:6])
			if err != nil {
				return err
			}
			len, err = CreateI32(len.ToI32() + 1)
			if err != nil {
				return err
			}
			for i := 1; i < 6; i++ {
				list.Value[i] = len.Bytes()[i-1]
			}
			list.Value = append(list.Value, obj.Bytes()...)
			vm.Push(ctx, list)
		case OP_LIST_GET:
			ip++
			ind, err := CreateObject(instr.Operands)
			if err != nil {
				return err
			}
			list, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			if list.Tag != TAG_LIST {
				return fmt.Errorf("expected list, got %v", list.Tag)
			}
			var obj CVMObject
			switch list.Value[0] {
			case TAG_I32:
				if len(list.Value[6:])/5-1 < int(ind.ToI32()) {
					return fmt.Errorf("index %d out of range", ind.ToI32())
				}
				offStart, OffEnd := (6 + ind.ToI32()*5), (ind.ToI32()*5 + 11)
				obj, err = CreateObject(list.Value[offStart:OffEnd])
			case TAG_F32:
				if len(list.Value[6:])/5-1 < int(ind.ToI32()) {
					return fmt.Errorf("index %d out of range", ind.ToI32())
				}
				offStart, OffEnd := (6 + ind.ToI32()*5), (ind.ToI32()*5 + 11)
				obj, err = CreateObject(list.Value[offStart:OffEnd])
			case TAG_BOOL:
				if len(list.Value[6:])/2-1 < int(ind.ToI32()) {
					return fmt.Errorf("index %d out of range", ind.ToI32())
				}
				offStart, OffEnd := (6 + ind.ToI32()*2), (ind.ToI32()*2 + 8)
				obj, err = CreateObject(list.Value[offStart:OffEnd])
			default:
				return fmt.Errorf("unexpected object tag %v", list.Value[0])
			}
			if err != nil {
				return err
			}
			vm.Push(ctx, obj)
		case OP_LIST_POP:
			ip++
			list, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			if list.Tag != TAG_LIST {
				return fmt.Errorf("expected list, got %v", list.Tag)
			}
			if len(list.Value) <= 6 {
				return fmt.Errorf("trying to pop element from empty list")
			}
			l, err := CreateObject(list.Value[1:6])
			if err != nil {
				return err
			}
			if l.ToI32() <= 0 {
				return fmt.Errorf("list is empty")
			}
			var obj CVMObject
			switch list.Value[0] {
			case TAG_I32:
				offStart, OffEnd := len(list.Value)-5, len(list.Value)
				obj, err = CreateObject(list.Value[offStart:OffEnd])
				list.Value = list.Value[:offStart]
			case TAG_F32:
				offStart, OffEnd := len(list.Value)-5, len(list.Value)
				obj, err = CreateObject(list.Value[offStart:OffEnd])
				list.Value = list.Value[:offStart]
			case TAG_BOOL:
				offStart, OffEnd := len(list.Value)-2, len(list.Value)
				obj, err = CreateObject(list.Value[offStart:OffEnd])
				list.Value = list.Value[:offStart]
			default:
				return fmt.Errorf("unexpected object tag %v", list.Value[0])
			}
			l, err = CreateObject(list.Value[1:6])
			if err != nil {
				return err
			}
			if l.ToI32() <= 0 {
				return fmt.Errorf("list is empty")
			}
			l, err = CreateI32(l.ToI32() - 1)
			if err != nil {
				return err
			}
			for i := 1; i < 6; i++ {
				list.Value[i] = l.Bytes()[i-1]
			}
			vm.Push(ctx, list)
			vm.Push(ctx, obj)
		case OP_LIST_INSERT:
			ip++
			ind, err := CreateObject(instr.Operands)
			if err != nil {
				return err
			}
			obj, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			list, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			if list.Tag != TAG_LIST {
				return fmt.Errorf("expected list, got %v", list.Tag)
			}
			l, err := CreateObject(list.Value[1:6])
			if err != nil {
				return err
			}
			if l.ToI32() <= 0 {
				return fmt.Errorf("list is empty")
			} else if l.ToI32() < ind.ToI32() {
				return fmt.Errorf("index %d out of range", ind.ToI32())
			}
			switch list.Value[0] {
			case TAG_I32:
				offStart := (6 + ind.ToI32()*5)
				list.Value = append(list.Value[:offStart], append(obj.Bytes(), list.Value[offStart:]...)...)
			case TAG_F32:
			case TAG_BOOL:
			default:
				return fmt.Errorf("unexpected object tag %v", list.Value[0])
			}
			l, err = CreateI32(l.ToI32() + 1)
			if err != nil {
				return err
			}
			for i := 1; i < 6; i++ {
				list.Value[i] = l.Bytes()[i-1]
			}
			vm.Push(ctx, list)
		case OP_LIST_REMOVE:
			ip++
			ind, err := CreateObject(instr.Operands)
			if err != nil {
				return err
			}
			list, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			if list.Tag != TAG_LIST {
				return fmt.Errorf("expected list, got %v", list.Tag)
			}
			l, err := CreateObject(list.Value[1:6])
			if err != nil {
				return err
			}
			if l.ToI32() <= 0 {
				return fmt.Errorf("list is empty")
			} else if l.ToI32() <= ind.ToI32() {
				return fmt.Errorf("index %d out of range", ind.ToI32())
			}
			switch list.Value[0] {
			case TAG_I32:
				offStart, OffEnd := (6 + ind.ToI32()*5), (ind.ToI32()*5 + 11)
				list.Value = append(list.Value[:offStart], list.Value[OffEnd:]...)
			case TAG_F32:
			case TAG_BOOL:
			default:
				return fmt.Errorf("unexpected object tag %v", list.Value[0])
			}
			l, err = CreateI32(l.ToI32() - 1)
			if err != nil {
				return err
			}
			for i := 1; i < 6; i++ {
				list.Value[i] = l.Bytes()[i-1]
			}
			vm.Push(ctx, list)
		case OP_LIST_REPLACE:
			ip++
			ind, err := CreateObject(instr.Operands)
			if err != nil {
				return err
			}
			obj, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			list, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			if list.Tag != TAG_LIST {
				return fmt.Errorf("expected list, got %v", list.Tag)
			}
			l, err := CreateObject(list.Value[1:6])
			if err != nil {
				return err
			}
			if l.ToI32() <= 0 {
				return fmt.Errorf("list is empty")
			} else if l.ToI32() <= ind.ToI32() {
				return fmt.Errorf("index %d out of range", ind.ToI32())
			}
			switch list.Value[0] {
			case TAG_I32:
				offStart, OffEnd := (6 + ind.ToI32()*5), (ind.ToI32()*5 + 11)
				list.Value = append(list.Value[:offStart], append(obj.Bytes(), list.Value[OffEnd:]...)...)
			case TAG_F32:
			case TAG_BOOL:
			default:
				return fmt.Errorf("unexpected object tag %v", list.Value[0])
			}
			vm.Push(ctx, list)
		default:
			return fmt.Errorf("unknown instruction of kind 0x%02x", instr.Kind)
		}
	}
	return nil
}
