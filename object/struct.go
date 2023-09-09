package object

import (
	"bytes"
	"fmt"
)

// constructor

func CreateStruct(data []byte) (CVMObject, error) {
	var obj CVMObject
	obj.Data = nil
	obj.Tag = TAG_STRUCT
	if data == nil {
		return obj, fmt.Errorf("empty data")
	}
	obj.Data = make([]byte, len(data[1:]))
	copy(obj.Data, data[1:])
	return obj, nil
}

// manipulation

func StringStruct(obj CVMObject) (string, error) {
	if obj.Tag != TAG_STRUCT {
		return "", fmt.Errorf("invalid struct tag %s", TagsName(obj.Tag))
	}
	var buf bytes.Buffer
	fmt.Fprint(&buf, "struct{ ")
	ln, err := Len(obj)
	if err != nil {
		return "", err
	}
	for i := 5; i < ln+5; i++ {
		fmt.Fprintf(&buf, "%s ", TagsName(obj.Data[i]))
	}
	fmt.Fprint(&buf, "}{ ")
	for i := ln + 5; i < len(obj.Data); {
		switch obj.Data[i] {
		case TAG_F32, TAG_I32, TAG_BOOL:
			s, err := Size(CVMObject{Tag: obj.Data[i]})
			if err != nil {
				return buf.String(), err
			}
			tO, err := CreateObject(obj.Data[i : i+s])
			if err != nil {
				return buf.String(), err
			}
			tS, err := String(tO)
			if err != nil {
				return buf.String(), err
			}
			fmt.Fprintf(&buf, "%s ", tS)
			i += s
		case TAG_STRING, TAG_LIST:
			s, err := Size(CVMObject{Tag: obj.Data[i], Data: obj.Data[i+1 : i+6]})
			if err != nil {
				return buf.String(), err
			}
			tO, err := CreateObject(obj.Data[i : i+s])
			if err != nil {
				return buf.String(), err
			}
			tS, err := String(tO)
			if err != nil {
				return buf.String(), err
			}
			fmt.Fprintf(&buf, "%s ", tS)
			i += s
		default:
			return buf.String(), fmt.Errorf("unexpected struct field tag %s", TagsName(obj.Data[i]))
		}
	}
	fmt.Fprint(&buf, "}")
	return buf.String(), nil
}

// actions

func GetStruct(strct, ind CVMObject) (CVMObject, error) {
	var obj CVMObject
	if strct.Tag != TAG_STRUCT {
		return obj, fmt.Errorf("expected struct, got %s", TagsName(strct.Tag))
	}
	temp, err := ValueI32(ind)
	if err != nil {
		return obj, err
	}
	indVal := int(temp)
	ln, err := Len(strct)
	if err != nil {
		return obj, err
	}
	if ln <= 0 {
		return obj, fmt.Errorf("struct is empty")
	}
	if ln <= indVal {
		return obj, fmt.Errorf("index %d out of range", indVal)
	}
	offStart := 5 + ln
	offEnd := 0
	switch strct.Data[offStart] {
	case TAG_BOOL, TAG_F32, TAG_I32:
		size, err := Size(CVMObject{Tag: strct.Data[0]})
		if err != nil {
			return strct, err
		}
		offStart += indVal * size
		offEnd = offStart + size
	case TAG_STRING:
		size := 0
		for i := 0; i < indVal+1; i++ {
			s, err := Size(CVMObject{Tag: TAG_STRING, Data: strct.Data[offStart+1 : offStart+5]})
			if err != nil {
				return strct, err
			}
			offStart += s
			size = s
		}
		offEnd = offStart
		offStart -= size
	case TAG_LIST:
		size := 0
		for i := 0; i < indVal+1; i++ {
			s, err := Size(CVMObject{Tag: TAG_LIST, Data: strct.Data[offStart+1 : offStart+7]})
			if err != nil {
				return strct, err
			}
			offStart += s
			size = s
		}
		offEnd = offStart
		offStart -= size
	}
	obj, err = CreateObject(strct.Data[offStart:offEnd])
	return obj, err
}

func SetStruct(oldStruct, ind, obj CVMObject) (CVMObject, error) {
	return oldStruct, nil
}
