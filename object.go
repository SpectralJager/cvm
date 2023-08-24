package main

import (
	"encoding/binary"
	"fmt"
)

const (
	TAG_I32 = iota
	TAG_BOOL
)

type CVMObject struct {
	Tag   byte
	Value []byte
}

func (o *CVMObject) String() string {
	switch o.Tag {
	case TAG_I32:
		return fmt.Sprintf("int.%d", o.ToI32())
	case TAG_BOOL:
		return fmt.Sprintf("bool.%v", o.ToBool())
	default:
		return fmt.Sprintf("unknown.%v", o.Value)
	}
}

func (o *CVMObject) ToI32() int32 {
	if o.Tag != TAG_I32 {
		return 0
	}
	val := binary.LittleEndian.Uint32(o.Value[:4])
	return int32(val)
}
func (o *CVMObject) ToBool() bool {
	if o.Tag != TAG_BOOL {
		return false
	}
	return o.Value[0] > 0
}

func CreateI32Object[T []byte | int32](val T) (CVMObject, error) {
	var obj CVMObject
	obj.Value = nil
	switch val := any(val).(type) {
	case []byte:
		obj.Tag = val[0]
		if obj.Tag != TAG_I32 {
			return obj, fmt.Errorf("invalid tag for i32 object: %02x", obj.Tag)
		}
		obj.Value = val[1:]
	case int32:
		obj.Tag = TAG_I32
		obj.Value = binary.LittleEndian.AppendUint32(obj.Value, uint32(val))
	}
	return obj, nil
}
func CreateBool[T []byte | bool](val T) (CVMObject, error) {
	var obj CVMObject
	obj.Value = nil
	switch val := any(val).(type) {
	case []byte:
		obj.Tag = val[0]
		if obj.Tag != TAG_BOOL {
			return obj, fmt.Errorf("invalid tag for bool object: %02x", obj.Tag)
		}
		obj.Value = val[1:]
	case bool:
		obj.Tag = TAG_BOOL
		if val {
			obj.Value = append(obj.Value, 1)
		} else {
			obj.Value = append(obj.Value, 0)
		}
	}
	return obj, nil
}
