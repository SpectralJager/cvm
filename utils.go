package cvm

import (
	"context"
	"cvm/object"
)

type binFunc func(obj1, obj2 object.CVMObject) (object.CVMObject, error)
type unaryFunc func(obj object.CVMObject) (object.CVMObject, error)
type ternaryFunc func(obj1, obj2, obj3 object.CVMObject) (object.CVMObject, error)
type nFunc func(obj1 object.CVMObject, objs []object.CVMObject) (object.CVMObject, error)

func TernaryOperation(ctx context.Context, vm *CVM, terOperation ternaryFunc) (object.CVMObject, error) {
	obj3, err := vm.Pop(ctx)
	if err != nil {
		return object.CVMObject{}, nil
	}
	obj2, err := vm.Pop(ctx)
	if err != nil {
		return object.CVMObject{}, nil
	}
	obj1, err := vm.Pop(ctx)
	if err != nil {
		return object.CVMObject{}, nil
	}
	return terOperation(obj1, obj2, obj3)
}

func BinaryOperation(ctx context.Context, vm *CVM, binOperation binFunc) (object.CVMObject, error) {
	obj2, err := vm.Pop(ctx)
	if err != nil {
		return object.CVMObject{}, nil
	}
	obj1, err := vm.Pop(ctx)
	if err != nil {
		return object.CVMObject{}, nil
	}
	return binOperation(obj1, obj2)
}

func UnaryOperation(ctx context.Context, vm *CVM, unaryOperation unaryFunc) (object.CVMObject, error) {
	obj, err := vm.Pop(ctx)
	if err != nil {
		return object.CVMObject{}, nil
	}
	return unaryOperation(obj)
}

func NOperation(ctx context.Context, vm *CVM, nOperation nFunc) (object.CVMObject, error) {
	nO, err := vm.Pop(ctx)
	if err != nil {
		return object.CVMObject{}, err
	}
	nV, err := object.ValueI32(nO)
	if err != nil {
		return object.CVMObject{}, err
	}
	objs := make([]object.CVMObject, 0, nV)
	for i := 0; i < int(nV); i++ {
		obj, err := vm.Pop(ctx)
		if err != nil {
			return object.CVMObject{}, err
		}
		objs = append(objs, obj)
	}
	obj, err := vm.Pop(ctx)
	if err != nil {
		return object.CVMObject{}, err
	}
	return nOperation(obj, objs)
}
