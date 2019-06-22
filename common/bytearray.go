package common

import (
	"bytes"
	"github.com/newkedison/go-utils/algorithm"
	"github.com/newkedison/go-utils/internal/types"
	"strconv"
)

type ByteArray []byte

func NewByteArray(args ...interface{}) *ByteArray {
	size := 0
	if len(args) > 0 {
		switch v := args[0].(type) {
		case int:
			size = v
		case []byte:
			ret := ByteArray(v)
			return &ret
		default:
			panic("Invalid type of 1st parameter for ByteArray.New, must be int")
		}
	}
	ret := make(ByteArray, size)
	return &ret
}

func (ba ByteArray) ToStringEx(
	withLen bool, sep string, prefix string, suffix string) string {
	var buffer bytes.Buffer
	if withLen {
		buffer.WriteString("[")
		buffer.WriteString(strconv.Itoa(len(ba)))
		buffer.WriteString("]")
	}
	for _, b := range ba {
		buffer.WriteString(prefix)
		buffer.WriteString(ByteToHexString(b))
		buffer.WriteString(suffix)
		buffer.WriteString(sep)
	}
	result := buffer.String()
	if len(ba) > 0 && len(sep) > 0 {
		return result[:len(result)-len(sep)]
	}
	return result
}

func (ba ByteArray) ToString() string {
	return ba.ToStringEx(true, " ", "", "")
}

func (ba ByteArray) Len() int {
	return len(ba)
}

func (ba *ByteArray) AppendByte(b byte) {
	*ba = append(*ba, b)
}

func (ba *ByteArray) AppendBytes(arr []byte) {
	*ba = append(*ba, arr...)
}

func (ba *ByteArray) AppendIntAsByte(i int) {
	*ba = append(*ba, byte(i&0xFF))
}

func (ba *ByteArray) AppendString(s string) {
	*ba = append(*ba, []byte(s)...)
}

func (ba *ByteArray) AddCrc16() {
	*ba = algorithm.AppendCrc16([]byte(*ba))
}

func (ba ByteArray) Crc8() byte {
	return algorithm.Crc8([]byte(ba))
}

func (ba ByteArray) Crc16() uint16 {
	return algorithm.Crc16([]byte(ba))
}

func (ba ByteArray) Crc32() uint32 {
	return algorithm.Crc32([]byte(ba))
}

func (ba ByteArray) Clone() ByteArray {
	return append(ByteArray{}, ba...)
}

func (ba *ByteArray) Assign(data []byte) {
	*ba = ByteArray(data)
}

func (ba *ByteArray) AssignByCopy(data []byte) {
	*ba = append(ByteArray{}, data...)
}

func (ba *ByteArray) ToProtoMessage() *types.WSByteArray {
	return &types.WSByteArray{
		Data: []byte(*ba),
	}
}

func (ba *ByteArray) FromProtoMessage(p *types.WSByteArray) {
	*ba = ByteArray(p.Data)
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (ba *ByteArray) MarshalBinary() (_ []byte, err error) {
	defer SetErrorWhenMarshalObjectErrorPanic("common.ByteArray", &err)()
	return MarshalProtoMessage(ba.ToProtoMessage())
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (ba *ByteArray) UnmarshalBinary(data []byte) error {
	_, err := ba.UnmarshalBinaryWithSize(data)
	return err
}

// UnmarshalBinaryWithSize implements the common.BinaryUnmarshalerWithSize interface.
func (ba *ByteArray) UnmarshalBinaryWithSize(data []byte) (_ int, err error) {
	defer SetErrorWhenUnmarshalObjectErrorPanic("common.ByteArray", &err)()
	var result types.WSByteArray
	used := UnmarshalProtoMessage(data, &result)
	ba.FromProtoMessage(&result)
	return used, nil
}

// vim: fdm=syntax fdn=1
