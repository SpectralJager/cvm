package object

import (
	"encoding/binary"
	"fmt"
)

// constructor
func CreateI32(val int32) (CVMObject, error) {
	var obj CVMObject
	obj.Data = nil
	obj.Tag = TAG_I32
	obj.Data = binary.LittleEndian.AppendUint32(obj.Data, uint32(val))
	return obj, nil
}

// manipulations

func ValueI32(obj CVMObject) (int32, error) {
	if obj.Tag != TAG_I32 {
		return 0, fmt.Errorf("can't get Data, object tag is %s, not i32", TagsName(obj.Tag))
	}
	val := binary.LittleEndian.Uint32(obj.Data[:4])
	return int32(val), nil
}

func AsI32(obj CVMObject) (CVMObject, error) {
	switch obj.Tag {
	case TAG_I32:
		return obj, nil
	case TAG_F32:
		val, err := ValueF32(obj)
		if err != nil {
			return CVMObject{}, err
		}
		return CreateI32(int32(val))
	default:
		return CVMObject{}, fmt.Errorf("can't convert %s to i32", TagsName(obj.Tag))
	}
}

func StringI32(obj CVMObject) (string, error) {
	if obj.Tag != TAG_I32 {
		return "", fmt.Errorf("expected i32, got %s", TagsName(obj.Tag))
	}
	val, err := ValueI32(obj)
	return fmt.Sprintf("(%s)%d", TagsName(obj.Tag), val), err
}

// actions

func NegI32(obj CVMObject) (CVMObject, error) {
	v, err := ValueI32(obj)
	if err != nil {
		return CVMObject{}, err
	}
	return CreateI32(-v)
}

func AddI32(obj1, obj2 CVMObject) (CVMObject, error) {
	v1, err := ValueI32(obj1)
	if err != nil {
		return CVMObject{}, err
	}
	v2, err := ValueI32(obj2)
	if err != nil {
		return CVMObject{}, err
	}
	return CreateI32(v1 + v2)
}

func SubI32(obj1, obj2 CVMObject) (CVMObject, error) {
	v1, err := ValueI32(obj1)
	if err != nil {
		return CVMObject{}, err
	}
	v2, err := ValueI32(obj2)
	if err != nil {
		return CVMObject{}, err
	}
	return CreateI32(v1 - v2)
}

func MulI32(obj1, obj2 CVMObject) (CVMObject, error) {
	v1, err := ValueI32(obj1)
	if err != nil {
		return CVMObject{}, err
	}
	v2, err := ValueI32(obj2)
	if err != nil {
		return CVMObject{}, err
	}
	return CreateI32(v1 * v2)
}

func DivI32(obj1, obj2 CVMObject) (CVMObject, error) {
	v1, err := ValueI32(obj1)
	if err != nil {
		return CVMObject{}, err
	}
	v2, err := ValueI32(obj2)
	if err != nil {
		return CVMObject{}, err
	}
	return CreateI32(v1 / v2)
}

func LtI32(obj1, obj2 CVMObject) (CVMObject, error) {
	v1, err := ValueI32(obj1)
	if err != nil {
		return CVMObject{}, err
	}
	v2, err := ValueI32(obj2)
	if err != nil {
		return CVMObject{}, err
	}
	return CreateBool(v1 < v2)
}

func GtI32(obj1, obj2 CVMObject) (CVMObject, error) {
	v1, err := ValueI32(obj1)
	if err != nil {
		return CVMObject{}, err
	}
	v2, err := ValueI32(obj2)
	if err != nil {
		return CVMObject{}, err
	}
	return CreateBool(v1 > v2)
}

func LeqI32(obj1, obj2 CVMObject) (CVMObject, error) {
	v1, err := ValueI32(obj1)
	if err != nil {
		return CVMObject{}, err
	}
	v2, err := ValueI32(obj2)
	if err != nil {
		return CVMObject{}, err
	}
	return CreateBool(v1 <= v2)
}

func GeqI32(obj1, obj2 CVMObject) (CVMObject, error) {
	v1, err := ValueI32(obj1)
	if err != nil {
		return CVMObject{}, err
	}
	v2, err := ValueI32(obj2)
	if err != nil {
		return CVMObject{}, err
	}
	return CreateBool(v1 >= v2)
}

func EqI32(obj1, obj2 CVMObject) (CVMObject, error) {
	v1, err := ValueI32(obj1)
	if err != nil {
		return CVMObject{}, err
	}
	v2, err := ValueI32(obj2)
	if err != nil {
		return CVMObject{}, err
	}
	return CreateBool(v1 == v2)
}
