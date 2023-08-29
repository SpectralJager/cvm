package object

import "fmt"

// constructor

func CreateBool(val bool) (CVMObject, error) {
	var obj CVMObject
	obj.Data = nil
	obj.Tag = TAG_BOOL
	if val {
		obj.Data = append(obj.Data, 1)
	} else {
		obj.Data = append(obj.Data, 0)
	}
	return obj, nil
}

// manipulation

func ValueBool(obj CVMObject) (bool, error) {
	if obj.Tag != TAG_BOOL {
		return false, fmt.Errorf("can't get Data, object tag is %s, not bool", TagsName(obj.Tag))
	}
	return obj.Data[0] > 0, nil
}

func AsBool(obj CVMObject) (CVMObject, error) {
	switch obj.Tag {
	case TAG_BOOL:
		return obj, nil
	case TAG_I32:
		val, err := ValueI32(obj)
		if err != nil {
			return CVMObject{}, err
		}
		return CreateBool(val != 0)
	case TAG_F32:
		val, err := ValueF32(obj)
		if err != nil {
			return CVMObject{}, err
		}
		return CreateBool(val != 0.0)
	case TAG_LIST:
		l, err := Len(obj)
		if err != nil {
			return CVMObject{}, err
		}
		return CreateBool(l > 0)
	default:
		return CVMObject{}, fmt.Errorf("can't convert %s to f32", TagsName(obj.Tag))
	}
}

func StringBool(obj CVMObject) (string, error) {
	if obj.Tag != TAG_BOOL {
		return "", fmt.Errorf("expected bool, got %s", TagsName(obj.Tag))
	}
	val, err := ValueBool(obj)
	return fmt.Sprintf("(%s)%v", TagsName(obj.Tag), val), err
}

// actions

func NotBool(obj CVMObject) (CVMObject, error) {
	v, err := ValueBool(obj)
	if err != nil {
		return CVMObject{}, err
	}
	return CreateBool(!v)
}

func AndBool(obj1, obj2 CVMObject) (CVMObject, error) {
	v1, err := ValueBool(obj1)
	if err != nil {
		return CVMObject{}, err
	}
	v2, err := ValueBool(obj2)
	if err != nil {
		return CVMObject{}, err
	}
	return CreateBool(v1 && v2)
}

func OrBool(obj1, obj2 CVMObject) (CVMObject, error) {
	v1, err := ValueBool(obj1)
	if err != nil {
		return CVMObject{}, err
	}
	v2, err := ValueBool(obj2)
	if err != nil {
		return CVMObject{}, err
	}
	return CreateBool(v1 || v2)
}

func NorBool(obj1, obj2 CVMObject) (CVMObject, error) {
	v1, err := ValueBool(obj1)
	if err != nil {
		return CVMObject{}, err
	}
	v2, err := ValueBool(obj2)
	if err != nil {
		return CVMObject{}, err
	}
	return CreateBool(!(v1 || v2))
}

func NandBool(obj1, obj2 CVMObject) (CVMObject, error) {
	v1, err := ValueBool(obj1)
	if err != nil {
		return CVMObject{}, err
	}
	v2, err := ValueBool(obj2)
	if err != nil {
		return CVMObject{}, err
	}
	return CreateBool(!(v1 && v2))
}

func XorBool(obj1, obj2 CVMObject) (CVMObject, error) {
	v1, err := ValueBool(obj1)
	if err != nil {
		return CVMObject{}, err
	}
	v2, err := ValueBool(obj2)
	if err != nil {
		return CVMObject{}, err
	}
	return CreateBool(v1 != v2)
}
