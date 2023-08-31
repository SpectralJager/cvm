package object

import (
	"bytes"
	"fmt"
	"strconv"
)

// constructor
func CreateString(str string) (CVMObject, error) {
	var obj CVMObject
	obj.Data = nil
	obj.Tag = TAG_STRING
	ln, err := CreateI32(int32(len(str)))
	if err != nil {
		return obj, err
	}
	obj.Data = Bytes(ln)
	obj.Data = append(obj.Data, []byte(str)...)
	return obj, nil
}

// manipulation

func ValueString(obj CVMObject) (string, error) {
	if obj.Tag != TAG_STRING {
		return "", fmt.Errorf("expected string, got %s", TagsName(obj.Tag))
	}
	val := bytes.NewBuffer(obj.Data[5:])
	return val.String(), nil
}

func AsString(obj CVMObject) (CVMObject, error) {
	switch obj.Tag {
	case TAG_STRING:
		return obj, nil
	case TAG_I32:
		val, err := ValueI32(obj)
		if err != nil {
			return CVMObject{}, err
		}
		return CreateString(strconv.FormatInt(int64(val), 10))
	case TAG_F32:
		val, err := ValueF32(obj)
		if err != nil {
			return CVMObject{}, err
		}
		return CreateString(strconv.FormatFloat(float64(val), 'e', -1, 32))
	default:
		return CVMObject{}, fmt.Errorf("can't convert %s to string", TagsName(obj.Tag))
	}
}

func StringString(obj CVMObject) (string, error) {
	if obj.Tag != TAG_STRING {
		return "", fmt.Errorf("expected string, got %s", TagsName(obj.Tag))
	}
	ln, err := Len(obj)
	if err != nil {
		return "", err
	}
	val, err := ValueString(obj)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("(%s)[%d]\"%s\"", TagsName(obj.Tag), ln, val), nil
}

// actions

func ConcatString(str1, str2 CVMObject) (CVMObject, error) {
	resObj := CVMObject{
		Tag: TAG_STRING,
	}
	var temp []byte
	ln1, err := CreateObject(str1.Data[:5])
	if err != nil {
		return resObj, err
	}
	ln2, err := CreateObject(str2.Data[:5])
	if err != nil {
		return resObj, err
	}
	resLen, err := AddI32(ln1, ln2)
	if err != nil {
		return resObj, err
	}
	temp = Bytes(resLen)
	temp = append(temp, append(str1.Data[5:], str2.Data[5:]...)...)
	resObj.Data = make([]byte, len(temp))
	copy(resObj.Data, temp)
	return resObj, nil
}
