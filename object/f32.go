package object

import (
	"encoding/binary"
	"fmt"
	"math"
)

func CreateF32(val float32) (CVMObject, error) {
	var obj CVMObject
	obj.Data = nil
	obj.Tag = TAG_F32
	obj.Data = binary.LittleEndian.AppendUint32(obj.Data, math.Float32bits(val))
	return obj, nil
}

func ValueF32(obj CVMObject) (float32, error) {
	if obj.Tag != TAG_F32 {
		return 0, fmt.Errorf("can't get Data, object tag is %s, not f32", TagsName(obj.Tag))
	}
	val := math.Float32frombits(binary.LittleEndian.Uint32(obj.Data[:4]))
	return val, nil
}

func AsF32(obj CVMObject) (CVMObject, error) {
	switch obj.Tag {
	case TAG_F32:
		return obj, nil
	case TAG_I32:
		val, err := ValueI32(obj)
		if err != nil {
			return CVMObject{}, err
		}
		return CreateF32(float32(val))
	default:
		return CVMObject{}, fmt.Errorf("can't convert %s to f32", TagsName(obj.Tag))
	}
}

func StringF32(obj CVMObject) (string, error) {
	if obj.Tag != TAG_F32 {
		return "", fmt.Errorf("expected f32, got %s", TagsName(obj.Tag))
	}
	val, err := ValueF32(obj)
	return fmt.Sprintf("(%s)%f", TagsName(obj.Tag), val), err
}

// actions

func NegF32(obj CVMObject) (CVMObject, error) {
	v, err := ValueF32(obj)
	if err != nil {
		return CVMObject{}, err
	}
	return CreateF32(-v)
}

func AddF32(obj1, obj2 CVMObject) (CVMObject, error) {
	v1, err := ValueF32(obj1)
	if err != nil {
		return CVMObject{}, err
	}
	v2, err := ValueF32(obj2)
	if err != nil {
		return CVMObject{}, err
	}
	return CreateF32(v1 + v2)
}

func SubF32(obj1, obj2 CVMObject) (CVMObject, error) {
	v1, err := ValueF32(obj1)
	if err != nil {
		return CVMObject{}, err
	}
	v2, err := ValueF32(obj2)
	if err != nil {
		return CVMObject{}, err
	}
	return CreateF32(v1 - v2)
}

func MulF32(obj1, obj2 CVMObject) (CVMObject, error) {
	v1, err := ValueF32(obj1)
	if err != nil {
		return CVMObject{}, err
	}
	v2, err := ValueF32(obj2)
	if err != nil {
		return CVMObject{}, err
	}
	return CreateF32(v1 * v2)
}

func DivF32(obj1, obj2 CVMObject) (CVMObject, error) {
	v1, err := ValueF32(obj1)
	if err != nil {
		return CVMObject{}, err
	}
	v2, err := ValueF32(obj2)
	if err != nil {
		return CVMObject{}, err
	}
	return CreateF32(v1 / v2)
}

func LtF32(obj1, obj2 CVMObject) (CVMObject, error) {
	v1, err := ValueF32(obj1)
	if err != nil {
		return CVMObject{}, err
	}
	v2, err := ValueF32(obj2)
	if err != nil {
		return CVMObject{}, err
	}
	return CreateBool(v1 < v2)
}

func GtF32(obj1, obj2 CVMObject) (CVMObject, error) {
	v1, err := ValueF32(obj1)
	if err != nil {
		return CVMObject{}, err
	}
	v2, err := ValueF32(obj2)
	if err != nil {
		return CVMObject{}, err
	}
	return CreateBool(v1 > v2)
}

func LeqF32(obj1, obj2 CVMObject) (CVMObject, error) {
	v1, err := ValueF32(obj1)
	if err != nil {
		return CVMObject{}, err
	}
	v2, err := ValueF32(obj2)
	if err != nil {
		return CVMObject{}, err
	}
	return CreateBool(v1 <= v2)
}

func GeqF32(obj1, obj2 CVMObject) (CVMObject, error) {
	v1, err := ValueF32(obj1)
	if err != nil {
		return CVMObject{}, err
	}
	v2, err := ValueF32(obj2)
	if err != nil {
		return CVMObject{}, err
	}
	return CreateBool(v1 >= v2)
}

func EqF32(obj1, obj2 CVMObject) (CVMObject, error) {
	v1, err := ValueF32(obj1)
	if err != nil {
		return CVMObject{}, err
	}
	v2, err := ValueF32(obj2)
	if err != nil {
		return CVMObject{}, err
	}
	return CreateBool(v1 == v2)
}
