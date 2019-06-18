package common

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
	"fmt"
	proto "github.com/golang/protobuf/proto"
	"github.com/newkedison/go-utils/algorithm"
	"math"
	"reflect"
)

var defaultByteOrder binary.ByteOrder = binary.LittleEndian

func SetByteOrder(o binary.ByteOrder) {
	defaultByteOrder = o
}

type BinaryUnmarshalerWithSize interface {
	UnmarshalBinaryWithSize(data []byte) (int, error)
}

type MarshalObjectError error
type UnmarshalObjectError error

func NewNotEnoughDataError(required, offer, offset int) UnmarshalObjectError {
	return UnmarshalObjectError(errors.New(fmt.Sprintf(
		"Not enought data, require %d, offer %d-%d=%d",
		required, offer, offset, offer-offset)))
}

func NewMarshalObjectError(err error) MarshalObjectError {
	return MarshalObjectError(err)
}

func NewUnmarshalObjectError(err error) UnmarshalObjectError {
	return UnmarshalObjectError(err)
}

func SetErrorWhenMarshalObjectErrorPanic(typeName string, err *error) func() {
	return func() {
		if r := recover(); r != nil {
			if s, ok := r.(MarshalObjectError); ok {
				*err = errors.New(
					fmt.Sprintf("Marshal %s fail: %s", typeName, error(s).Error()))
			} else {
				panic(r)
			}
		}
	}
}

func SetErrorWhenUnmarshalObjectErrorPanic(typeName string, err *error) func() {
	return func() {
		if r := recover(); r != nil {
			if s, ok := r.(UnmarshalObjectError); ok {
				*err = errors.New(
					fmt.Sprintf("Unmarshal %s fail: %s", typeName, error(s).Error()))
			} else {
				panic(r)
			}
		}
	}
}

func CheckBufferSize(buf []byte, requiredLength int, args ...interface{}) {
	offset := 0
	offsetType := reflect.TypeOf(offset)
	if len(args) > 0 {
		if reflect.TypeOf(args[0]).ConvertibleTo(offsetType) {
			offset = reflect.ValueOf(args[0]).Convert(offsetType).Interface().(int)
		}
	}
	if len(buf) < offset+requiredLength {
		panic(NewNotEnoughDataError(requiredLength, len(buf), offset))
	}
}

func MarshalSimpleType(d interface{}) []byte {
	tmp := make([]byte, 8)
	switch v := d.(type) {
	case byte:
		return []byte{v}
	case int8:
		return []byte{byte(v)}
	case int16:
		defaultByteOrder.PutUint16(tmp, uint16(v))
		return tmp[:2]
	case uint16:
		defaultByteOrder.PutUint16(tmp, v)
		return tmp[:2]
	case int:
		defaultByteOrder.PutUint32(tmp, uint32(v))
		return tmp[:4]
	case uint:
		defaultByteOrder.PutUint32(tmp, uint32(v))
		return tmp[:4]
	case int32:
		defaultByteOrder.PutUint32(tmp, uint32(v))
		return tmp[:4]
	case uint32:
		defaultByteOrder.PutUint32(tmp, v)
		return tmp[:4]
	case int64:
		defaultByteOrder.PutUint64(tmp, uint64(v))
		return tmp
	case uint64:
		defaultByteOrder.PutUint64(tmp, v)
		return tmp
	case float32:
		return MarshalSimpleType(math.Float32bits(v))
	case float64:
		return MarshalSimpleType(math.Float64bits(v))
	}
	panic(NewMarshalObjectError(errors.New("MarshalSimpleType: Unknown type")))
}

func UnmarshalSimpleType(p interface{}, data []byte) int {
	t := reflect.TypeOf(p)
	if t.Kind() == reflect.Ptr {
		switch v := p.(type) {
		case *byte:
			CheckBufferSize(data, 1)
			*v = data[0]
			return 1
		case *int8:
			CheckBufferSize(data, 1)
			*v = int8(data[0])
			return 1
		case *int16:
			CheckBufferSize(data, 2)
			*v = int16(defaultByteOrder.Uint16(data))
			return 2
		case *uint16:
			CheckBufferSize(data, 2)
			*v = defaultByteOrder.Uint16(data)
			return 2
		case *int:
			CheckBufferSize(data, 4)
			*v = int(int32(defaultByteOrder.Uint32(data)))
			return 4
		case *uint:
			CheckBufferSize(data, 4)
			*v = uint(defaultByteOrder.Uint32(data))
			return 4
		case *int32:
			CheckBufferSize(data, 4)
			*v = int32(defaultByteOrder.Uint32(data))
			return 4
		case *uint32:
			CheckBufferSize(data, 4)
			*v = defaultByteOrder.Uint32(data)
			return 4
		case *int64:
			CheckBufferSize(data, 8)
			*v = int64(defaultByteOrder.Uint64(data))
			return 8
		case *uint64:
			CheckBufferSize(data, 8)
			*v = defaultByteOrder.Uint64(data)
			return 8
		case *float32:
			CheckBufferSize(data, 4)
			*v = math.Float32frombits(defaultByteOrder.Uint32(data))
			return 4
		case *float64:
			CheckBufferSize(data, 8)
			*v = math.Float64frombits(defaultByteOrder.Uint64(data))
			return 8
		}
	}
	panic(NewUnmarshalObjectError(errors.New("UnmarshalSimpleType: Unknown type")))
}

func MarshalString(s string) []byte {
	var buf bytes.Buffer
	buf.Write(MarshalSimpleType(uint16(len(s))))
	buf.WriteString(s)
	return buf.Bytes()
}

func UnmarshalString(dest *string, data []byte) int {
	CheckBufferSize(data, 2)
	var len uint16
	offset := UnmarshalSimpleType(&len, data)
	CheckBufferSize(data, int(len), offset)
	*dest = string(data[offset : offset+int(len)])
	return offset + int(len)
}

func MarshalObject(obj encoding.BinaryMarshaler) []byte {
	binary, err := obj.MarshalBinary()
	if err != nil {
		panic(NewMarshalObjectError(err))
	}
	var buf bytes.Buffer
	buf.Write(binary)
	return buf.Bytes()
}

func UnmarshalObject(dest BinaryUnmarshalerWithSize, data []byte) int {
	n, err := dest.UnmarshalBinaryWithSize(data)
	if err != nil {
		panic(NewUnmarshalObjectError(err))
	}
	return n
}

func MarshalProtoMessage(pb proto.Message) ([]byte, error) {
	data, err := proto.Marshal(pb)
	if err != nil {
		panic(NewMarshalObjectError(err))
	}
	result := PackLength(uint32(len(data)))
	result.AppendBytes(data)
	result.AddCrc16()
	return []byte(result), nil
}

func UnmarshalProtoMessage(data []byte, out proto.Message) int {
	dataLength, err := UnpackLength(ByteArray(data))
	if err != nil {
		panic(NewUnmarshalObjectError(err))
	}
	lengthSize := ByteCountOfPackedLength(dataLength)
	totalLength := lengthSize + int(dataLength) + 2
	CheckBufferSize(data, totalLength)
	if !algorithm.VerifyCrc16(data[:totalLength]) {
		panic(NewUnmarshalObjectError(errors.New("CRC check fail")))
	}
	err = proto.Unmarshal(data[lengthSize:lengthSize+int(dataLength)], out)
	if err != nil {
		panic(NewUnmarshalObjectError(err))
	}
	return totalLength
}
