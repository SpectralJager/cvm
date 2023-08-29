package object

import (
	"bytes"
	"fmt"
)

// constructor
func CreateList(val []byte) (CVMObject, error) {
	list := CVMObject{
		Tag:  TAG_LIST,
		Data: make([]byte, len(val)),
	}
	copy(list.Data, val)
	return list, nil
}

// manipulation

func StringList(obj CVMObject) (string, error) {
	if obj.Tag != TAG_LIST {
		return "", fmt.Errorf("expected list, got %s", TagsName(obj.Tag))
	}
	var buf bytes.Buffer
	l, err := Len(obj)
	if err != nil {
		return buf.String(), err
	}
	fmt.Fprintf(&buf, "(%s.%s)[%d]{ ", TagsName(obj.Tag), TagsName(obj.Data[0]), l)
	size, err := Size(obj)
	if err != nil {
		return buf.String(), err
	}
	for i := 6; i < size-1; {
		s := 0
		str := ""
		switch obj.Data[i] {
		case TAG_I32, TAG_F32, TAG_BOOL:
			s, err = Size(CVMObject{Tag: obj.Data[i]})
			if err != nil {
				return buf.String(), err
			}
			obj, err := CreateObject(obj.Data[i : i+s])
			if err != nil {
				return buf.String(), err
			}
			str, err = String(obj)
			if err != nil {
				return buf.String(), err
			}
		case TAG_LIST:
			s, err = Size(CVMObject{Tag: obj.Data[i], Data: obj.Data[i+1 : i+7]})
			if err != nil {
				return buf.String(), err
			}
			str, err = StringList(CVMObject{Tag: obj.Data[i]})
		default:
			return buf.String(), fmt.Errorf("unexpected list item tag %s", TagsName(obj.Data[i]))
		}
		i += s
		fmt.Fprintf(&buf, "%s ", str)
	}
	fmt.Fprint(&buf, "}")
	return buf.String(), nil
}

func ChangeLengthList(list CVMObject, NewLength int32) (CVMObject, error) {
	len, err := CreateI32(NewLength)
	if err != nil {
		return list, err
	}
	for i := 1; i < 6; i++ {
		list.Data[i] = Bytes(len)[i-1]
	}
	return list, nil
}

// actions

func LenList(list CVMObject) (CVMObject, error) {
	var obj CVMObject
	ln, err := Len(list)
	if err != nil {
		return obj, err
	}
	obj, err = CreateI32(int32(ln))
	return obj, err
}

func GetList(list, ind CVMObject) (CVMObject, error) {
	var obj CVMObject
	if list.Tag != TAG_LIST {
		return obj, fmt.Errorf("expected list, got %s", TagsName(list.Tag))
	}
	temp, err := ValueI32(ind)
	if err != nil {
		return obj, err
	}
	indVal := int(temp)
	ln, err := Len(list)
	if err != nil {
		return obj, err
	}
	if ln <= 0 {
		return obj, fmt.Errorf("list is empty")
	}
	if ln <= indVal {
		return obj, fmt.Errorf("index %d out of range", indVal)
	}
	size, err := Size(CVMObject{Tag: list.Data[0]})
	if err != nil {
		return obj, err
	}
	switch list.Data[0] {
	case TAG_I32:
		offStart, OffEnd := (6 + indVal*size), (indVal*size + 11)
		obj, err = CreateObject(list.Data[offStart:OffEnd])
	case TAG_F32:
		offStart, OffEnd := (6 + indVal*size), (indVal*size + 11)
		obj, err = CreateObject(list.Data[offStart:OffEnd])
	case TAG_BOOL:
		offStart, OffEnd := (6 + indVal*size), (indVal*size + 8)
		obj, err = CreateObject(list.Data[offStart:OffEnd])
	default:
		return obj, fmt.Errorf("unexpected object tag %v", TagsName(list.Data[0]))
	}
	return obj, nil
}

func RemoveList(oldList, ind CVMObject) (CVMObject, error) {
	list, _ := CreateList(oldList.Data)
	if list.Tag != TAG_LIST {
		return list, fmt.Errorf("expected list, got %s", TagsName(list.Tag))
	}
	if len(list.Data) <= 6 {
		return list, fmt.Errorf("trying to pop element from empty list")
	}
	temp, err := ValueI32(ind)
	if err != nil {
		return list, err
	}
	indVal := int(temp)
	ln, err := CreateObject(list.Data[1:6])
	if err != nil {
		return list, err
	}
	lenVal, err := ValueI32(ln)
	if err != nil {
		return list, err
	}
	if lenVal <= 0 {
		return list, fmt.Errorf("list is empty")
	}
	if int(lenVal) <= indVal {
		return list, fmt.Errorf("index %d out of range", indVal)
	}
	size, err := Size(CVMObject{Tag: list.Data[0]})
	if err != nil {
		return list, err
	}
	prefSize := 6
	if err != nil {
		return list, err
	}
	offStart, offEnd := (prefSize + indVal*size), (prefSize + indVal*size + size)
	list.Data = append(list.Data[:offStart], list.Data[offEnd:]...)
	list, err = ChangeLengthList(list, lenVal-1)
	if err != nil {
		return list, err
	}
	return list, nil
}

func InsertList(oldList, ind, obj CVMObject) (CVMObject, error) {
	list, _ := CreateList(oldList.Data)
	if list.Tag != TAG_LIST {
		return list, fmt.Errorf("expected list, got %s", TagsName(list.Tag))
	}
	temp, err := ValueI32(ind)
	if err != nil {
		return list, err
	}
	indVal := int(temp)
	ln, err := CreateObject(list.Data[1:6])
	if err != nil {
		return list, err
	}
	lenVal, err := ValueI32(ln)
	if err != nil {
		return list, err
	}
	if int(lenVal) < indVal {
		return list, fmt.Errorf("index %d out of range", indVal)
	}
	size, err := Size(CVMObject{Tag: list.Data[0]})
	if err != nil {
		return list, err
	}
	prefSize := 6
	offStart := (prefSize + indVal*size)
	list.Data = append(list.Data[:offStart], append(Bytes(obj), list.Data[offStart:]...)...)
	list, err = ChangeLengthList(list, lenVal+1)
	if err != nil {
		return list, err
	}
	return list, nil
}

func ReplaceList(oldList, ind, obj CVMObject) (CVMObject, error) {
	list, err := CreateList(oldList.Data)
	list, err = RemoveList(list, ind)
	if err != nil {
		return list, err
	}
	list, err = InsertList(list, ind, obj)
	return list, err
}
