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
	lV, err := Len(obj)
	if err != nil {
		return obj, err
	}
	for i := 5; i < lV+5; i++ {
		tg := obj.Data[i]
		o, err := CreateDefault(tg)
		if err != nil {
			return obj, err
		}
		obj.Data = append(obj.Data, Bytes(o)...)
	}
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
		case TAG_STRING:
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
