package main

import (
	"encoding/binary"
	"fmt"
	"math"
)

const (
	TAG_I32 = iota
	TAG_BOOL
	TAG_F32

	TAG_LIST // tag.elemTag.data...
)

type CVMObject struct {
	Tag   byte
	Value []byte
}

func (o *CVMObject) String() string {
	switch o.Tag {
	case TAG_I32:
		return fmt.Sprintf("(i32)%d", o.ToI32())
	case TAG_BOOL:
		return fmt.Sprintf("(bool)%v", o.ToBool())
	case TAG_F32:
		return fmt.Sprintf("(f32)%v", o.ToF32())
	case TAG_LIST:
		str := fmt.Sprintf("(list.")
		switch o.Value[0] {
		case TAG_I32:
			str += fmt.Sprintf("i32)[ ")
		case TAG_BOOL:
			str += fmt.Sprintf("bool)[ ")
		case TAG_F32:
			str += fmt.Sprintf("f32)[ ")
		default:
			str += fmt.Sprintf("unknown)%v", o.Value)
		}
		var obj CVMObject
		var err error
		for i := 1; i < len(o.Value[1:]); {
			switch o.Value[i] {
			case TAG_I32:
				obj, err = CreateObject(o.Value[i : i+5])
				i += 5
			case TAG_BOOL:
				obj, err = CreateObject(o.Value[i : i+2])
				i += 2
			case TAG_F32:
				obj, err = CreateObject(o.Value[i : i+5])
				i += 5
			default:
				err = fmt.Errorf("unknown tag %d", o.Value[i])
			}
			if err != nil {
				return str + "'" + err.Error() + "' ]"
			}
			str += obj.String() + " "
		}
		str += "]"
		return str
	default:
		return fmt.Sprintf("(unknown)%v", o.Value)
	}
}
func (o *CVMObject) Bytes() []byte {
	return append([]byte{o.Tag}, o.Value...)
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
func (o *CVMObject) ToF32() float32 {
	val := math.Float32frombits(binary.LittleEndian.Uint32(o.Value[:4]))
	return val
}
func (o *CVMObject) I32ToF32() float32 {
	val := o.ToI32()
	return float32(val)
}
func (o *CVMObject) I32ToBool() bool {
	val := o.ToI32()
	if val == 0 {
		return false
	}
	return true
}
func (o *CVMObject) F32ToI32() int32 {
	val := o.ToF32()
	return int32(val)
}
func (o *CVMObject) F32ToBool() bool {
	val := o.ToF32()
	if val == 0.0 {
		return false
	}
	return true
}
func CreateI32(val int32) (CVMObject, error) {
	var obj CVMObject
	obj.Value = nil
	obj.Tag = TAG_I32
	obj.Value = binary.LittleEndian.AppendUint32(obj.Value, uint32(val))
	return obj, nil
}
func CreateBool(val bool) (CVMObject, error) {
	var obj CVMObject
	obj.Value = nil
	obj.Tag = TAG_BOOL
	if val {
		obj.Value = append(obj.Value, 1)
	} else {
		obj.Value = append(obj.Value, 0)
	}
	return obj, nil
}
func CreateF32(val float32) (CVMObject, error) {
	var obj CVMObject
	obj.Value = nil
	obj.Tag = TAG_F32
	obj.Value = binary.LittleEndian.AppendUint32(obj.Value, math.Float32bits(val))
	return obj, nil
}
func CreateList(val []byte) (CVMObject, error) {
	return CVMObject{
		Tag:   TAG_LIST,
		Value: val,
	}, nil
}
func CreateObject(val []byte) (CVMObject, error) {
	var obj CVMObject
	obj.Value = nil
	switch val[0] {
	case TAG_I32:
		obj.Tag = val[0]
		obj.Value = val[1:]
	case TAG_F32:
		obj.Tag = val[0]
		obj.Value = val[1:]
	case TAG_BOOL:
		obj.Tag = val[0]
		obj.Value = val[1:]
	default:
		return obj, fmt.Errorf("unknown tag %v", val[0])
	}
	return obj, nil
}
