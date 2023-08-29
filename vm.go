package cvm

import (
	"bytes"
	"context"
	"cvm/instruction"
	"cvm/object"
	"fmt"
)

const (
	STACK_SIZE       = 2048
	HEAP_SIZE        = 2048
	STACK_FRAME_SIZE = 2048
)

type CVM struct {
	Stack      [STACK_SIZE]object.CVMObject
	Heap       [HEAP_SIZE]object.CVMObject
	StackFrame [STACK_FRAME_SIZE]Frame
	SP, HP, FP uint
}

func (vm *CVM) New(ctx context.Context, obj object.CVMObject) error {
	if vm.HP >= HEAP_SIZE {
		return fmt.Errorf("heap overflow")
	}
	vm.Heap[vm.HP] = obj
	vm.HP++
	return nil
}
func (vm *CVM) Load(ctx context.Context, ind uint32) (object.CVMObject, error) {
	if uint32(vm.HP) < ind {
		return object.CVMObject{}, fmt.Errorf("symbol with index %d not found", ind)
	}
	obj := vm.Heap[ind]
	return obj, nil
}
func (vm *CVM) Free(ctx context.Context, ind uint32) error {
	if uint32(vm.HP) < ind {
		return fmt.Errorf("symbol with index %d not found", ind)
	}
	vm.Heap[ind] = object.CVMObject{}
	return nil
}
func (vm *CVM) Save(ctx context.Context, ind uint32, obj object.CVMObject) error {
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
func (vm *CVM) Push(ctx context.Context, obj object.CVMObject) error {
	if vm.SP >= STACK_SIZE {
		return fmt.Errorf("stack overflow")
	}
	vm.Stack[vm.SP] = obj
	vm.SP++
	return nil
}
func (vm *CVM) Pop(ctx context.Context) (object.CVMObject, error) {
	if vm.SP == 0 {
		return object.CVMObject{}, fmt.Errorf("stack is empty")
	}
	vm.SP--
	obj := vm.Stack[vm.SP]
	return obj, nil
}
func (vm *CVM) Trace() string {
	var buf bytes.Buffer
	fmt.Fprint(&buf, "=== Heap:\n")
	for i := 0; i < int(vm.HP); i++ {
		if vm.Heap[i].Data != nil {
			str, err := object.String(vm.Heap[i])
			if err != nil {
				panic(err)
			}
			fmt.Fprintf(&buf, "\t$%03d -> %s\n", i, str)
		}
	}
	fmt.Fprint(&buf, "=== StackFrace:\n")
	for i := 0; i < int(vm.FP); i++ {
		fmt.Fprintf(&buf, "\t$%03d -> %s\n", i, vm.StackFrame[i].String())
	}
	fmt.Fprint(&buf, "=== Stack:\n")
	for i := 0; i < int(vm.SP); i++ {
		str, err := object.String(vm.Stack[i])
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(&buf, "\t$%03d -> %s\n", i, str)
	}
	return buf.String()
}

func (vm *CVM) Execute(ctx context.Context, instrs []instruction.Instruction) error {
	for ip := uint32(0); ip < uint32(len(instrs)); {
		instr := instrs[ip]
		switch instr.Kind {
		case instruction.OP_NULL:
			ip++
		case instruction.OP_HALT:
			return nil
		case instruction.OP_I32_LOAD:
			ip++
			obj, err := object.CreateObject(instr.Operands)
			if err != nil {
				return err
			}
			vm.Push(ctx, obj)
		case instruction.OP_I32_NEG:
			ip++
			resObj, err := UnaryOperation(ctx, vm, object.NegI32)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case instruction.OP_I32_ADD:
			ip++
			resObj, err := BinaryOperation(ctx, vm, object.AddI32)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case instruction.OP_I32_SUB:
			ip++
			resObj, err := BinaryOperation(ctx, vm, object.SubI32)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case instruction.OP_I32_MUL:
			ip++
			resObj, err := BinaryOperation(ctx, vm, object.MulI32)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case instruction.OP_I32_DIV:
			ip++
			resObj, err := BinaryOperation(ctx, vm, object.DivI32)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case instruction.OP_I32_LT:
			ip++
			resObj, err := BinaryOperation(ctx, vm, object.LtI32)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case instruction.OP_I32_GT:
			ip++
			resObj, err := BinaryOperation(ctx, vm, object.GtI32)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case instruction.OP_I32_LEQ:
			ip++
			resObj, err := BinaryOperation(ctx, vm, object.LeqI32)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case instruction.OP_I32_GEQ:
			ip++
			resObj, err := BinaryOperation(ctx, vm, object.GeqI32)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case instruction.OP_I32_EQ:
			ip++
			resObj, err := BinaryOperation(ctx, vm, object.EqI32)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case instruction.OP_I32_TO_F32:
			ip++
			resObj, err := UnaryOperation(ctx, vm, object.AsF32)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case instruction.OP_I32_TO_BOOL:
			ip++
			resObj, err := UnaryOperation(ctx, vm, object.AsBool)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case instruction.OP_BOOL_LOAD:
			ip++
			obj, err := object.CreateObject(instr.Operands)
			if err != nil {
				return err
			}
			vm.Push(ctx, obj)
		case instruction.OP_BOOL_NOT:
			ip++
			resObj, err := UnaryOperation(ctx, vm, object.NotBool)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case instruction.OP_BOOL_AND:
			ip++
			resObj, err := BinaryOperation(ctx, vm, object.AndBool)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case instruction.OP_BOOL_OR:
			ip++
			resObj, err := BinaryOperation(ctx, vm, object.OrBool)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case instruction.OP_BOOL_NAND:
			ip++
			resObj, err := BinaryOperation(ctx, vm, object.NandBool)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case instruction.OP_BOOL_NOR:
			ip++
			resObj, err := BinaryOperation(ctx, vm, object.NorBool)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case instruction.OP_BOOL_XOR:
			ip++
			resObj, err := BinaryOperation(ctx, vm, object.XorBool)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case instruction.OP_F32_LOAD:
			ip++
			obj, err := object.CreateObject(instr.Operands)
			if err != nil {
				return err
			}
			vm.Push(ctx, obj)
		case instruction.OP_F32_NEG:
			ip++
			resObj, err := UnaryOperation(ctx, vm, object.NegF32)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case instruction.OP_F32_ADD:
			ip++
			resObj, err := BinaryOperation(ctx, vm, object.AddF32)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case instruction.OP_F32_SUB:
			ip++
			resObj, err := BinaryOperation(ctx, vm, object.SubF32)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case instruction.OP_F32_MUL:
			ip++
			resObj, err := BinaryOperation(ctx, vm, object.MulF32)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case instruction.OP_F32_DIV:
			ip++
			resObj, err := BinaryOperation(ctx, vm, object.DivF32)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case instruction.OP_F32_LT:
			ip++
			resObj, err := BinaryOperation(ctx, vm, object.LtF32)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case instruction.OP_F32_GT:
			ip++
			resObj, err := BinaryOperation(ctx, vm, object.GtF32)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case instruction.OP_F32_LEQ:
			ip++
			resObj, err := BinaryOperation(ctx, vm, object.LeqF32)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case instruction.OP_F32_GEQ:
			ip++
			resObj, err := BinaryOperation(ctx, vm, object.GeqF32)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case instruction.OP_F32_EQ:
			ip++
			resObj, err := BinaryOperation(ctx, vm, object.EqF32)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case instruction.OP_F32_TO_I32:
			ip++
			resObj, err := UnaryOperation(ctx, vm, object.AsI32)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case instruction.OP_F32_TO_BOOL:
			ip++
			resObj, err := UnaryOperation(ctx, vm, object.AsBool)
			if err != nil {
				return err
			}
			vm.Push(ctx, resObj)
		case instruction.OP_JUMP:
			addr, err := object.CreateObject(instr.Operands)
			if err != nil {
				return err
			}
			v, err := object.ValueI32(addr)
			if err != nil {
				return err
			}
			ip = uint32(v)
		case instruction.OP_JUMPC:
			obj, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			val, err := object.ValueBool(obj)
			if err != nil {
				return err
			}
			if !val {
				ip++
				continue
			}
			addr, err := object.CreateObject(instr.Operands)
			if err != nil {
				return err
			}
			v, err := object.ValueI32(addr)
			if err != nil {
				return err
			}
			ip = uint32(v)
		case instruction.OP_JUMPNC:
			obj, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			val, err := object.ValueBool(obj)
			if err != nil {
				return err
			}
			if val {
				ip++
				continue
			}
			addr, err := object.CreateObject(instr.Operands)
			if err != nil {
				return err
			}
			v, err := object.ValueI32(addr)
			if err != nil {
				return err
			}
			ip = uint32(v)
		case instruction.OP_BLOCK_START:
			ip++
			addr, err := object.CreateObject(instr.Operands)
			if err != nil {
				return err
			}
			retIp, err := object.ValueI32(addr)
			if err != nil {
				return err
			}
			vm.PushFrame(ctx, Frame{
				StackOffset: int(vm.SP),
				HeapOffset:  int(vm.HP),
				ReturnIP:    uint32(retIp),
				FrameOffset: -1,
			})
		case instruction.OP_BLOCK_BR:
			fr, err := vm.LastFrame(ctx)
			if err != nil {
				return err
			}
			ip = fr.ReturnIP
		case instruction.OP_BLOCK_END:
			ip++
			fr, err := vm.PopFrame(ctx)
			if err != nil {
				return err
			}
			vm.HP = uint(fr.HeapOffset)
			vm.SP = uint(fr.StackOffset)
		case instruction.OP_BLOCK_LOAD:
			ip++
			ind, err := object.CreateObject(instr.Operands)
			if err != nil {
				return err
			}
			fr, err := vm.LastFrame(ctx)
			if err != nil {
				return err
			}
			indVal, err := object.ValueI32(ind)
			if err != nil {
				return err
			}
			obj, err := vm.Load(ctx, uint32(int32(fr.HeapOffset)+indVal))
			if err != nil {
				return err
			}
			vm.Push(ctx, obj)
		case instruction.OP_BLOCK_SAVE:
			ip++
			ind, err := object.CreateObject(instr.Operands)
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
			indVal, err := object.ValueI32(ind)
			if err != nil {
				return err
			}
			err = vm.Save(ctx, uint32(int32(fr.HeapOffset)+indVal), obj)
			if err != nil {
				return err
			}
		case instruction.OP_LOAD:
			ip++
			ind, err := object.CreateObject(instr.Operands)
			if err != nil {
				return err
			}
			indVal, err := object.ValueI32(ind)
			if err != nil {
				return err
			}
			obj, err := vm.Load(ctx, uint32(indVal))
			vm.Push(ctx, obj)
		case instruction.OP_SAVE:
			ip++
			ind, err := object.CreateObject(instr.Operands)
			if err != nil {
				return err
			}
			obj, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			indVal, err := object.ValueI32(ind)
			if err != nil {
				return err
			}
			err = vm.Save(ctx, uint32(indVal), obj)
			if err != nil {
				return err
			}
		case instruction.OP_FREE:
			ip++
			ind, err := object.CreateObject(instr.Operands)
			if err != nil {
				return err
			}
			indVal, err := object.ValueI32(ind)
			if err != nil {
				return err
			}
			err = vm.Free(ctx, uint32(indVal))
			if err != nil {
				return err
			}
		case instruction.OP_NEW:
			ip++
			obj, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
			err = vm.New(ctx, obj)
			if err != nil {
				return err
			}
		case instruction.OP_POP:
			ip++
			_, err := vm.Pop(ctx)
			if err != nil {
				return err
			}
		case instruction.OP_FUNC_CALL:
			ip++
			addr, err := object.CreateObject(instr.Operands[:5])
			if err != nil {
				return err
			}
			argsLen, err := object.CreateObject(instr.Operands[5:10])
			if err != nil {
				return err
			}
			argLenVal, err := object.ValueI32(argsLen)
			if err != nil {
				return err
			}
			vm.PushFrame(ctx, Frame{
				StackOffset: int(vm.SP) - int(argLenVal),
				HeapOffset:  int(vm.HP),
				ReturnIP:    uint32(ip),
				FrameOffset: int(vm.FP),
			})
			addrVal, err := object.ValueI32(addr)
			if err != nil {
				return err
			}
			ip = uint32(addrVal)
		case instruction.OP_FUNC_RET:
			fr, err := vm.LastFuncFrame(ctx)
			if err != nil {
				return err
			}
			retLen, err := object.CreateObject(instr.Operands[:5])
			if err != nil {
				return err
			}
			retLenVal, err := object.ValueI32(retLen)
			if err != nil {
				return err
			}
			ip = fr.ReturnIP
			objs := vm.Stack[int(vm.SP)-int(retLenVal) : vm.SP]
			vm.HP = uint(fr.HeapOffset)
			vm.SP = uint(fr.StackOffset)
			vm.FP = uint(fr.FrameOffset)
			for _, obj := range objs {
				vm.Push(ctx, obj)
			}
		case instruction.OP_LOCAL_LOAD:
			ip++
			ind, err := object.CreateObject(instr.Operands)
			if err != nil {
				return err
			}
			fr, err := vm.LastFuncFrame(ctx)
			if err != nil {
				return err
			}
			indVal, err := object.ValueI32(ind)
			if err != nil {
				return err
			}
			obj, err := vm.Load(ctx, uint32(int32(fr.HeapOffset)+indVal))
			if err != nil {
				return err
			}
			vm.Push(ctx, obj)
		case instruction.OP_LOCAL_SAVE:
			ip++
			ind, err := object.CreateObject(instr.Operands)
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
			indVal, err := object.ValueI32(ind)
			if err != nil {
				return err
			}
			err = vm.Save(ctx, uint32(int32(fr.HeapOffset)+indVal), obj)
			if err != nil {
				return err
			}
		case instruction.OP_LIST_NEW:
			ip++
			obj, err := object.CreateList(instr.Operands)
			if err != nil {
				return err
			}
			vm.Push(ctx, obj)
		case instruction.OP_LIST_LENGTH:
			ip++
			obj, err := UnaryOperation(ctx, vm, object.LenList)
			if err != nil {
				return err
			}
			vm.Push(ctx, obj)
		case instruction.OP_LIST_GET:
			ip++
			obj, err := BinaryOperation(ctx, vm, object.GetList)
			if err != nil {
				return err
			}
			vm.Push(ctx, obj)
		case instruction.OP_LIST_INSERT:
			ip++
			list, err := TernaryOperation(ctx, vm, object.InsertList)
			if err != nil {
				return err
			}
			vm.Push(ctx, list)
		case instruction.OP_LIST_REMOVE:
			ip++
			list, err := BinaryOperation(ctx, vm, object.RemoveList)
			if err != nil {
				return err
			}
			vm.Push(ctx, list)
		case instruction.OP_LIST_REPLACE:
			ip++
			list, err := TernaryOperation(ctx, vm, object.ReplaceList)
			if err != nil {
				return err
			}
			vm.Push(ctx, list)
		default:
			return fmt.Errorf("unknown instruction of kind 0x%02x", instr.Kind)
		}
	}
	return nil
}
