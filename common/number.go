package common

import (
	"github.com/newkedison/go-utils/internal/types"
	"math"
	"reflect"
	"strconv"
)

type Number float64

const (
	MaxNumber Number = Number(math.MaxFloat64)
	MinNumber Number = Number(-math.MaxFloat64)
	// Ref: https://en.wikipedia.org/wiki/Double-precision_floating-point_format#IEEE_754_double-precision_binary_floating-point_format:_binary64
	// The 53-bit significand precision gives from 15 to 17 significant decimal digits precision (2^−53 ≈ 1.11 × 10^−16). If a decimal string with at most 15 significant digits is converted to IEEE 754 double-precision representation, and then converted back to a decimal string with the same number of digits, the final result should match the original string. If an IEEE 754 double-precision number is converted to a decimal string with at least 17 significant digits, and then converted back to double-precision representation, the final result must match the original number.
	MaxIntNumber  int64  = 1e16
	MinIntNumber  int64  = -1e16
	MaxUintNumber uint64 = 1e16
)

var (
	NumberType reflect.Type = reflect.TypeOf(Number(0))
)

func Float32Epsilon() float32 {
	return math.Nextafter32(1.0, 2.0) - 1.0
}

func Float64Epsilon() float64 {
	return math.Nextafter(1.0, 2.0) - 1.0
}

func NewNumber(d interface{}) Number {
	switch v := d.(type) {
	case int64:
		if v > MaxIntNumber {
			panic("Assign an int64 bigger than " +
				strconv.FormatInt(MaxIntNumber, 10) +
				" will lost significant digits, " +
				"used float64 instead")
		}
		if v < MinIntNumber {
			panic("Assign an int64 smaller than " +
				strconv.FormatInt(MinIntNumber, 10) +
				" will lost significant digits, " +
				"used float64 instead")
		}
	case uint64:
		if v > uint64(MaxIntNumber) {
			panic("Assign an uint64 bigger than " +
				strconv.FormatUint(MaxUintNumber, 10) +
				" will lost significant digits, " +
				"used float64 instead")
		}
	}
	if reflect.TypeOf(d).ConvertibleTo(NumberType) {
		return reflect.ValueOf(d).Convert(NumberType).Interface().(Number)
	}
	panic("Invalid number type")
}

func (v Number) ToByte() byte {
	return byte(uint64(v))
}

func (v Number) ToInt8() int8 {
	return int8(int64(v))
}

func (v Number) ToUint8() uint8 {
	return uint8(uint64(v))
}

func (v Number) ToInt16() int16 {
	return int16(int64(v))
}

func (v Number) ToUint16() uint16 {
	return uint16(uint64(v))
}

func (v Number) ToInt32() int32 {
	return int32(int64(v))
}

func (v Number) ToUint32() uint32 {
	return uint32(uint64(v))
}

func (v Number) ToInt64() int64 {
	return int64(v)
}

func (v Number) ToUint64() uint64 {
	return uint64(v)
}

func (v Number) ToFloat32() float32 {
	return float32(v)
}

func (v Number) ToFloat64() float64 {
	return float64(v)
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (v Number) MarshalBinary() (_ []byte, err error) {
	defer SetErrorWhenMarshalObjectErrorPanic("common.Number", &err)()
	p := &types.WSNumber{
		Value: v.ToFloat64(),
	}
	return MarshalProtoMessage(p)
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (v *Number) UnmarshalBinary(data []byte) error {
	_, err := v.UnmarshalBinaryWithSize(data)
	return err
}

// UnmarshalBinaryWithSize implements the common.BinaryUnmarshalerWithSize interface.
func (v *Number) UnmarshalBinaryWithSize(data []byte) (_ int, err error) {
	defer SetErrorWhenUnmarshalObjectErrorPanic("common.Number", &err)()
	var result types.WSNumber
	used := UnmarshalProtoMessage(data, &result)
	*v = Number(result.Value)
	return used, nil
}
