package common_test

import (
	"errors"
	"github.com/newkedison/go-utils/common"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestMarshalSimple(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(common.MarshalSimpleType(byte(0x55)), []byte{0x55})
	assert.Equal(common.MarshalSimpleType(int8(-1)), []byte{0xFF})
	assert.Equal(common.MarshalSimpleType(uint8(255)), []byte{0xFF})
	assert.Equal(common.MarshalSimpleType(int16(-1)), []byte{0xFF, 0xFF})
	assert.Equal(common.MarshalSimpleType(uint16(65535)), []byte{0xFF, 0xFF})
	assert.Equal(common.MarshalSimpleType(1), []byte{0x01, 0x00, 0x00, 0x00})
	assert.Equal(common.MarshalSimpleType(-1), []byte{0xFF, 0xFF, 0xFF, 0xFF})
	assert.Equal(common.MarshalSimpleType(uint(1)), []byte{0x01, 0x00, 0x00, 0x00})
	assert.Equal(common.MarshalSimpleType(int32(-1)), []byte{0xFF, 0xFF, 0xFF, 0xFF})
	assert.Equal(common.MarshalSimpleType(uint32(4294967295)),
		[]byte{0xFF, 0xFF, 0xFF, 0xFF})
	assert.Equal(common.MarshalSimpleType(int64(-1)),
		[]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})
	assert.Equal(common.MarshalSimpleType(uint64(18446744073709551615)),
		[]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})
	assert.Equal(common.MarshalSimpleType(float32(1.234)),
		[]byte{0xB6, 0xF3, 0x9D, 0x3F})
	assert.Equal(common.MarshalSimpleType(float64(1.234)),
		[]byte{0x58, 0x39, 0xB4, 0xC8, 0x76, 0xBE, 0xF3, 0x3F})
	assert.Panics(func() { common.MarshalSimpleType("I will panic...") })
}

func TestUnmarshalSimpleType(t *testing.T) {
	assert := assert.New(t)
	var b byte
	offset := common.UnmarshalSimpleType(&b, []byte{0xAA})
	assert.EqualValues(offset, 1)
	assert.EqualValues(b, 0xAA)
	var i8 int8
	offset = common.UnmarshalSimpleType(&i8, []byte{0xFF})
	assert.EqualValues(offset, 1)
	assert.EqualValues(i8, -1)
	var u8 uint8
	offset = common.UnmarshalSimpleType(&u8, []byte{0xFF})
	assert.EqualValues(offset, 1)
	assert.EqualValues(u8, 0xFF)
	var i16 int16
	offset = common.UnmarshalSimpleType(&i16, []byte{0xFF, 0xFF})
	assert.EqualValues(offset, 2)
	assert.EqualValues(i16, -1)
	var u16 uint16
	offset = common.UnmarshalSimpleType(&u16, []byte{0xFF, 0xFF})
	assert.EqualValues(offset, 2)
	assert.EqualValues(u16, 0xFFFF)
	var i32 int32
	offset = common.UnmarshalSimpleType(&i32, []byte{0xFF, 0xFF, 0xFF, 0xFF})
	assert.EqualValues(offset, 4)
	assert.EqualValues(i32, -1)
	var u32 uint32
	offset = common.UnmarshalSimpleType(&u32, []byte{0xFF, 0xFF, 0xFF, 0xFF})
	assert.EqualValues(offset, 4)
	assert.EqualValues(u32, 0xFFFFFFFF)
	var i int
	offset = common.UnmarshalSimpleType(&i, []byte{0xFF, 0xFF, 0xFF, 0xFF})
	assert.EqualValues(offset, 4)
	assert.EqualValues(i, -1)
	var u uint
	offset = common.UnmarshalSimpleType(&u, []byte{0xFF, 0xFF, 0xFF, 0xFF})
	assert.EqualValues(offset, 4)
	assert.EqualValues(u, 0xFFFFFFFF)
	var i64 int64
	offset = common.UnmarshalSimpleType(&i64,
		[]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})
	assert.EqualValues(offset, 8)
	assert.EqualValues(i64, -1)
	var u64 uint64
	offset = common.UnmarshalSimpleType(&u64,
		[]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})
	assert.EqualValues(offset, 8)
	assert.EqualValues(u64, uint64(0xFFFFFFFFFFFFFFFF))
	epsilon32 := math.Nextafter32(1, 2) - 1
	epsilon64 := math.Nextafter(1, 2) - 1
	var f32 float32
	offset = common.UnmarshalSimpleType(&f32, []byte{0xB6, 0xF3, 0x9D, 0x3F})
	assert.EqualValues(offset, 4)
	assert.InEpsilon(f32, 1.234, float64(epsilon32))
	var f64 float64
	offset = common.UnmarshalSimpleType(&f64,
		[]byte{0x58, 0x39, 0xB4, 0xC8, 0x76, 0xBE, 0xF3, 0x3F})
	assert.EqualValues(offset, 8)
	assert.InEpsilon(f64, 1.234, epsilon64)
	assert.Panics(func() { common.UnmarshalSimpleType(1, []byte{}) })
	assert.Panics(func() { common.UnmarshalSimpleType(1.1, []byte{}) })
	assert.Panics(func() { common.UnmarshalSimpleType("", []byte{}) })
}

func TestMarshalString(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(common.MarshalString("AAA"), []byte{0x03, 0x00, 0x41, 0x41, 0x41})
	assert.Equal(common.MarshalString("你好"),
		[]byte{0x06, 0x00, 0xe4, 0xbd, 0xa0, 0xe5, 0xa5, 0xbd}) // UTF-8
}

func TestUnmarshalString(t *testing.T) {
	assert := assert.New(t)
	var s string
	assert.Panics(func() { common.UnmarshalString(&s, nil) })
	assert.Panics(func() { common.UnmarshalString(&s, []byte{0x00}) })
	assert.Panics(func() { common.UnmarshalString(&s, []byte{0x02, 0x00, 0x00}) })
	n := common.UnmarshalString(&s, []byte{0x03, 0x00, 0x41, 0x41, 0x41})
	assert.Equal(n, 5)
	assert.Equal(s, "AAA")
	n = common.UnmarshalString(
		&s, []byte{0x06, 0x00, 0xe4, 0xbd, 0xa0, 0xe5, 0xa5, 0xbd})
	assert.Equal(n, 8)
	assert.Equal(s, "你好")
}

type marshalableObject int

func (v marshalableObject) MarshalBinary() ([]byte, error) {
	if int(v) == 0 {
		return nil, errors.New("0 is error")
	}
	return common.MarshalSimpleType(int(v)), nil
}

func (v *marshalableObject) UnmarshalBinaryWithSize(data []byte) (int, error) {
	var i int
	common.UnmarshalSimpleType(&i, data)
	*v = marshalableObject(i)
	if *v == 0 {
		return 0, errors.New("0 is error")
	}
	return 4, nil
}

func TestMarshalObject(t *testing.T) {
	assert := assert.New(t)
	data := common.MarshalObject(marshalableObject(1))
	assert.Equal(data, []byte{0x01, 0x00, 0x00, 0x00})
	assert.Panics(func() { common.MarshalObject(marshalableObject(0)) })
}

func TestUnmarshalObject(t *testing.T) {
	assert := assert.New(t)
	var o marshalableObject
	assert.Equal(common.UnmarshalObject(&o, []byte{0x0A, 0x00, 0x00, 0x00}), 4)
	assert.EqualValues(int(o), 10)
	assert.Panics(func() { common.UnmarshalObject(&o, []byte{0x00, 0x00, 0x00, 0x00}) })
}

type myByteOrder struct{}

func (myByteOrder) Uint16(b []byte) uint16 {
	return 0
}

func (myByteOrder) PutUint16(b []byte, v uint16) {
}

func (myByteOrder) Uint32(b []byte) uint32 {
	return 42
}

func (myByteOrder) PutUint32(b []byte, v uint32) {
	b[0] = 0xAA
	b[1] = 0xAA
	b[2] = 0xAA
}

func (myByteOrder) Uint64(b []byte) uint64 {
	return 0
}

func (myByteOrder) PutUint64(b []byte, v uint64) {
}

func (myByteOrder) String() string { return "myByteOrder" }

func TestSetByteOrder(t *testing.T) {
	assert := assert.New(t)
	common.SetByteOrder(myByteOrder{})
	assert.Equal(common.MarshalSimpleType(uint32(0)), []byte{0xAA, 0xAA, 0xAA, 0x00})
	var v uint32
	assert.EqualValues(common.UnmarshalSimpleType(&v, []byte{0xAA, 0xFF, 0xAA, 0xFF}), 4)
	assert.EqualValues(v, 42)
	//   data := common.MarshalObject(marshalableObject(1))
	//   assert.Equal(data, []byte{0xAA, 0xAA, 0xAA, 0x00, 0xAA, 0xAA, 0xAA, 0x00})
}

func TestSetErrorWhenMarshalObjectErrorPanic(t *testing.T) {
	assert := assert.New(t)
	var err error = errors.New("hello")
	defer func() {
		r := recover()
		assert.Nil(r)
		assert.Equal(err.Error(), "Marshal aaa fail: bbb")
	}()
	defer common.SetErrorWhenMarshalObjectErrorPanic("aaa", &err)()
	panic(common.NewMarshalObjectError(errors.New("bbb")))
}

func TestSetErrorWhenMarshalObjectErrorPanic2(t *testing.T) {
	assert := assert.New(t)
	var err error = errors.New("hello")
	defer func() {
		r := recover()
		assert.NotNil(r)
		assert.IsType(r, 0)
		assert.Equal(err.Error(), "hello")
	}()
	defer common.SetErrorWhenMarshalObjectErrorPanic("aaa", &err)()
	panic(42)
}

func TestSetErrorWhenUnmarshalObjectErrorPanic(t *testing.T) {
	assert := assert.New(t)
	var err error = errors.New("hello")
	defer func() {
		r := recover()
		assert.Nil(r)
		assert.Equal(err.Error(), "Unmarshal aaa fail: bbb")
	}()
	defer common.SetErrorWhenUnmarshalObjectErrorPanic("aaa", &err)()
	panic(common.NewUnmarshalObjectError(errors.New("bbb")))
}

func TestSetErrorWhenUnmarshalObjectErrorPanic2(t *testing.T) {
	assert := assert.New(t)
	var err error = errors.New("hello")
	defer func() {
		r := recover()
		assert.NotNil(r)
		assert.IsType(r, 0)
		assert.Equal(err.Error(), "hello")
	}()
	defer common.SetErrorWhenUnmarshalObjectErrorPanic("aaa", &err)()
	panic(42)
}

func TestMarshalProtoMessage(t *testing.T) {
	assert := assert.New(t)
	assert.Panics(func() { common.MarshalProtoMessage(nil) })
}
