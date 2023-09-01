package object

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	case TAG_BOOL:
		val, err := ValueBool(obj)
		if err != nil {
			return CVMObject{}, err
		}
		return CreateString(strconv.FormatBool(val))
	case TAG_LIST:
		ln, err := Len(obj)
		if err != nil {
			return CVMObject{}, err
		}
		var buf strings.Builder
		buf.WriteString("[")
		for i := 0; i < ln; i++ {
			ind, err := CreateI32(int32(i))
			if err != nil {
				return CVMObject{}, err
			}
			item, err := GetList(obj, ind)
			if err != nil {
				return CVMObject{}, err
			}
			res, err := AsString(item)
			if err != nil {
				return CVMObject{}, err
			}
			val, err := ValueString(res)
			if err != nil {
				return CVMObject{}, err
			}
			buf.WriteString(" ")
			buf.WriteString(val)
		}
		buf.WriteString(" ]")
		return CreateString(buf.String())
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

func LenString(str CVMObject) (CVMObject, error) {
	ln, err := Len(str)
	if err != nil {
		return CVMObject{}, err
	}
	lnObj, err := CreateI32(int32(ln))
	return lnObj, err
}

func SplitString(str CVMObject, sep CVMObject) (CVMObject, error) {
	var list CVMObject
	if str.Tag != TAG_STRING {
		return list, fmt.Errorf("invalid string, got %s", TagsName(str.Tag))
	}
	if sep.Tag != TAG_STRING {
		return list, fmt.Errorf("invalid string, got %s", TagsName(sep.Tag))
	}
	buf := make([]byte, 0, 6)
	buf = append(buf, TAG_STRING)
	buf = append(buf, TAG_I32)
	buf = binary.LittleEndian.AppendUint32(buf, 0)
	list, err := CreateList(buf)
	if err != nil {
		return list, err
	}
	valueStr, err := ValueString(str)
	if err != nil {
		return list, err
	}
	sepStr, err := ValueString(sep)
	if err != nil {
		return list, err
	}
	res := strings.Split(valueStr, sepStr)
	for i, s := range res {
		sO, err := CreateString(s)
		if err != nil {
			return list, err
		}
		iO, err := CreateI32(int32(i))
		if err != nil {
			return list, err
		}
		list, err = InsertList(list, iO, sO)
		if err != nil {
			return list, err
		}
	}
	return list, err
}

func FormatString(str CVMObject, data []CVMObject) (CVMObject, error) {
	var resObj CVMObject
	valStr, err := ValueString(str)
	if err != nil {
		return resObj, err
	}
	// temp := strings.Split(valStr, "%.")
	// l := len(temp) - 1
	l := strings.Count(valStr, "%.")
	g := len(data)
	if l != g {
		return resObj, fmt.Errorf("invalid value format points (%d) and processing data (%d)", l, g)
	}
	// var resVal strings.Builder
	for i := l - 1; i >= 0; i-- {
		sO, err := AsString(data[i])
		if err != nil {
			return resObj, err
		}
		sV, err := ValueString(sO)
		if err != nil {
			return resObj, err
		}
		valStr = strings.Replace(valStr, "%.", sV, 1)
	}
	resObj, err = CreateString(valStr)
	return resObj, nil
}

func PrintString(obj CVMObject) (CVMObject, error) {
	frmtO, err := CreateString("%.")
	if err != nil {
		return CVMObject{}, err
	}
	res, err := FormatString(frmtO, []CVMObject{obj})
	if err != nil {
		return CVMObject{}, err
	}
	resV, err := ValueString(res)
	if err != nil {
		return CVMObject{}, err
	}
	fmt.Fprint(os.Stdout, resV)
	return CVMObject{}, nil
}

func PrintfString(f CVMObject, objs []CVMObject) (CVMObject, error) {
	res, err := FormatString(f, objs)
	if err != nil {
		return CVMObject{}, err
	}
	resV, err := ValueString(res)
	if err != nil {
		return CVMObject{}, err
	}
	fmt.Fprint(os.Stdout, resV)
	return CVMObject{}, nil
}

func PrintlnString(obj CVMObject) (CVMObject, error) {
	f, err := CreateString("%.")
	if err != nil {
		return CVMObject{}, err
	}

	res, err := FormatString(f, []CVMObject{obj})
	if err != nil {
		return CVMObject{}, err
	}
	resV, err := ValueString(res)
	if err != nil {
		return CVMObject{}, err
	}
	fmt.Fprintln(os.Stdout, resV)
	return CVMObject{}, nil
}
