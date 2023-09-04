package object

import (
	"bufio"
	"fmt"
	"os"
)

const (
	TAG_I32 byte = iota // tag.data
	TAG_BOOL
	TAG_F32

	TAG_LIST   // tag.elemTag.len.data...
	TAG_STRING // tag.len.data...
	TAG_STRUCT // tag.len.{fieldTags}...{data}...
)

type CVMObject struct {
	Tag  byte
	Data []byte
}

func String(obj CVMObject) (string, error) {
	switch obj.Tag {
	case TAG_I32:
		return StringI32(obj)
	case TAG_BOOL:
		return StringBool(obj)
	case TAG_F32:
		return StringF32(obj)
	case TAG_LIST:
		return StringList(obj)
	case TAG_STRING:
		return StringString(obj)
	case TAG_STRUCT:
		return StringStruct(obj)
	default:
		return fmt.Sprintf("(unknown)%v", obj.Data), nil
	}
}

func Value(obj CVMObject) (any, error) {
	switch obj.Tag {
	case TAG_I32:
		return ValueI32(obj)
	case TAG_F32:
		return ValueF32(obj)
	case TAG_BOOL:
		return ValueBool(obj)
	case TAG_STRING:
		return ValueString(obj)
	default:
		return nil, fmt.Errorf("can't get value for tag %v", TagsName(obj.Tag))
	}
}

func TagsName(tag byte) string {
	switch tag {
	case TAG_I32:
		return "i32"
	case TAG_F32:
		return "f32"
	case TAG_BOOL:
		return "bool"
	case TAG_LIST:
		return "list"
	case TAG_STRING:
		return "string"
	default:
		return "unknown"
	}
}

func Bytes(obj CVMObject) []byte {
	return append([]byte{obj.Tag}, obj.Data...)
}

func CreateObject(val []byte) (CVMObject, error) {
	var obj CVMObject
	obj.Data = nil
	switch val[0] {
	case TAG_I32, TAG_F32, TAG_BOOL, TAG_STRING, TAG_LIST, TAG_STRUCT:
		obj.Tag = val[0]
		obj.Data = val[1:]
	default:
		return obj, fmt.Errorf("unknown tag %v", val[0])
	}
	return obj, nil
}

func CreateDefault(target byte) (CVMObject, error) {
	switch target {
	case TAG_I32:
		return CreateI32(0)
	case TAG_F32:
		return CreateF32(0.0)
	case TAG_BOOL:
		return CreateBool(false)
	case TAG_STRING:
		return CreateString("")
	default:
		return CVMObject{}, fmt.Errorf("cant create object with target %s", TagsName(target))
	}
}

func Len(obj CVMObject) (int, error) {
	switch obj.Tag {
	case TAG_LIST:
		l, err := CreateObject(obj.Data[1:6])
		if err != nil {
			return 0, err
		}
		val, err := ValueI32(l)
		return int(val), err
	case TAG_STRING, TAG_STRUCT:
		l, err := CreateObject(obj.Data[:5])
		if err != nil {
			return 0, err
		}
		val, err := ValueI32(l)
		return int(val), err
	default:
		return 0, fmt.Errorf("can't get len of %s", TagsName(obj.Tag))
	}
}

func Size(obj CVMObject) (int, error) {
	switch obj.Tag {
	case TAG_I32:
		return 5, nil
	case TAG_F32:
		return 5, nil
	case TAG_BOOL:
		return 2, nil
	case TAG_STRING:
		l, err := Len(obj)
		if err != nil {
			return 0, err
		}
		return 6 + l, nil
	case TAG_LIST:
		l, err := Len(obj)
		if err != nil {
			return 0, err
		}
		itemSize := 0
		switch obj.Data[0] {
		case TAG_I32, TAG_BOOL, TAG_F32:
			itemSize, err = Size(CVMObject{Tag: obj.Data[0]})
			if err != nil {
				return 0, err
			}
			return l*itemSize + 7, nil
		case TAG_STRING:
			for i := 6; i < len(obj.Data); {
				s, err := Size(CVMObject{Tag: obj.Data[i], Data: obj.Data[i+1 : i+6]})
				if err != nil {
					return 0, err
				}
				i += s
				itemSize += s
			}
			return itemSize + 7, nil
		default:
			return 0, fmt.Errorf("unexpected item tag %s", TagsName(obj.Data[0]))
		}
	default:
		return 0, fmt.Errorf("unknown tag %v", obj.Tag)
	}
}

func Print(obj CVMObject) (CVMObject, error) {
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

func Printf(f CVMObject, objs []CVMObject) (CVMObject, error) {
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

func Println(obj CVMObject) (CVMObject, error) {
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

func Read() (CVMObject, error) {
	res, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return CVMObject{}, err
	}
	resObj, err := CreateString(res)
	return resObj, err
}
