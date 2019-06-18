package common_test

import (
	"errors"
	"github.com/newkedison/go-utils/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewByteArray(t *testing.T) {
	assert := assert.New(t)
	ba := common.NewByteArray()
	assert.Equal(len(*ba), 0)
	ba = common.NewByteArray(10)
	assert.Equal(len(*ba), 10)
	ba = common.NewByteArray([]byte{0x01, 0x02})
	assert.Equal(len(*ba), 2)
	tmp := []byte{0x00}
	ba = common.NewByteArray(tmp)
	assert.Equal(len(*ba), 1)
	(*ba)[0] = 0xFF
	assert.EqualValues(tmp[0], 0xFF)

	assert.Panics(func() { common.NewByteArray("a") })
}

func TestToStringEx(t *testing.T) {
	assert := assert.New(t)
	var ba common.ByteArray
	assert.Equal(ba.ToStringEx(false, "", "", ""), "")
	assert.Equal(ba.ToStringEx(false, ",", "0x", ""), "")
	assert.Equal(ba.ToStringEx(true, ",", "0x", "XXX"), "[0]")
	ba = append(ba, 0x00)
	assert.Equal(ba.ToStringEx(false, "", "", ""), "00")
	assert.Equal(ba.ToStringEx(false, ",", "0x", ""), "0x00")
	assert.Equal(ba.ToStringEx(true, ",", "0x", ""), "[1]0x00")
	assert.Equal(ba.ToStringEx(true, ",", "", "H"), "[1]00H")
	assert.Equal(ba.ToStringEx(true, ",", "<<", ">>"), "[1]<<00>>")
	ba = append(ba, []byte{0x01, 0x02}...)
	assert.Equal(ba.ToStringEx(false, "", "", ""), "000102")
	assert.Equal(ba.ToStringEx(false, ",", "0x", ""), "0x00,0x01,0x02")
	assert.Equal(ba.ToStringEx(true, ",", "0x", ""), "[3]0x00,0x01,0x02")
	assert.Equal(ba.ToStringEx(true, ",", "", "H"), "[3]00H,01H,02H")
	assert.Equal(ba.ToStringEx(true, "|", "<<", ">>"), "[3]<<00>>|<<01>>|<<02>>")
}

func TestToString(t *testing.T) {
	assert := assert.New(t)
	var ba common.ByteArray
	assert.Equal(ba.ToString(), "[0]")
	ba = append(ba, 0x00)
	assert.Equal(ba.ToString(), "[1]00")
	ba = common.ByteArray([]byte{0x00, 0x11, 0x22, 0xFF})
	assert.Equal(ba.ToString(), "[4]00 11 22 FF")
}

func TestLen(t *testing.T) {
	assert := assert.New(t)
	var ba common.ByteArray
	assert.Equal(ba.Len(), 0)
	ba = append(ba, 0x00)
	assert.Equal(ba.Len(), 1)
	ba = common.ByteArray([]byte{0x00, 0x11, 0x22, 0xFF})
	assert.Equal(ba.Len(), 4)
}

func TestAppendByte(t *testing.T) {
	assert := assert.New(t)
	var ba common.ByteArray
	assert.Equal(ba.ToString(), "[0]")
	ba.AppendByte(0x00)
	assert.Equal(ba.ToString(), "[1]00")
}

func TestAppendBytes(t *testing.T) {
	assert := assert.New(t)
	var ba common.ByteArray
	assert.Equal(ba.ToString(), "[0]")
	ba.AppendBytes([]byte{0x00, 0x01, 0x02})
	assert.Equal(ba.ToString(), "[3]00 01 02")
}

func TestAppendIntAsByte(t *testing.T) {
	assert := assert.New(t)
	var ba common.ByteArray
	assert.Equal(ba.ToString(), "[0]")
	ba.AppendIntAsByte(0)
	ba.AppendIntAsByte(255)
	ba.AppendIntAsByte(256)
	ba.AppendIntAsByte(65538)
	assert.Equal(ba.ToString(), "[4]00 FF 00 02")
}

func TestAppendString(t *testing.T) {
	assert := assert.New(t)
	var ba common.ByteArray
	assert.Equal(ba.ToString(), "[0]")
	ba.AppendString("AAA")
	assert.Equal(ba.ToString(), "[3]41 41 41")
}

func TestAddCrc16(t *testing.T) {
	assert := assert.New(t)
	var ba common.ByteArray
	assert.Equal(ba.ToString(), "[0]")
	ba.AddCrc16()
	assert.Equal(ba.ToString(), "[2]FF FF")
	ba = []byte{0x10}
	ba.AddCrc16()
	assert.Equal(ba.ToString(), "[3]10 BE 8C")
	ba.AddCrc16()
	assert.Equal(ba.ToString(), "[5]10 BE 8C 00 00")
}

func TestCrc8(t *testing.T) {
	assert := assert.New(t)
	assert.EqualValues(common.ByteArray{0x10}.Crc8(), 0x70)
}

func TestCrc16(t *testing.T) {
	assert := assert.New(t)
	assert.EqualValues(common.ByteArray{0x10}.Crc16(), 0x8CBE)
}

func TestCrc32(t *testing.T) {
	assert := assert.New(t)
	assert.EqualValues(common.ByteArray{0x10}.Crc32(), 0xCFB5FFE9)
}

func TestClone(t *testing.T) {
	assert := assert.New(t)
	ba := common.ByteArray([]byte{0x01, 0x02})
	ba2 := ba.Clone()
	assert.Equal(len(ba2), len(ba))
	assert.Equal(ba2[0], ba[0])
	ba2[0] += 1
	assert.NotEqual(ba2[0], ba[0])
}

func TestAssign(t *testing.T) {
	assert := assert.New(t)
	ba := common.ByteArray([]byte{0x01, 0x02})
	data := []byte{0xFF}
	ba.Assign(data)
	assert.Equal(len(ba), 1)
	assert.EqualValues(ba[0], 0xFF)
	ba[0] = 0xAA
	assert.EqualValues(data[0], 0xAA)
}

func TestAssignByCopy(t *testing.T) {
	assert := assert.New(t)
	ba := common.ByteArray([]byte{0x01, 0x02})
	data := []byte{0xFF}
	ba.AssignByCopy(data)
	assert.Equal(len(ba), 1)
	assert.EqualValues(ba[0], 0xFF)
	ba[0] = 0xAA
	assert.EqualValues(data[0], 0xFF)
}

func TestMarshalBinary(t *testing.T) {
	assert := assert.New(t)
	ba := common.ByteArray([]byte{0x01, 0x02})
	buf, err := ba.MarshalBinary()
	assert.NoError(err)
	assert.Equal(buf, []byte{0x04, 0x0A, 0x02, 0x01, 0x02, 0xF7, 0x89})
}

func TestUnmarshalBinary(t *testing.T) {
	assert := assert.New(t)
	unmarshalError := common.NewUnmarshalObjectError(errors.New(""))
	ba := common.ByteArray([]byte{0x01, 0x02})
	buf, err := ba.MarshalBinary()
	assert.NoError(err)
	assert.Equal(len(buf), 7)
	var ba2 common.ByteArray
	err = ba2.UnmarshalBinary(nil)
	assert.NotNil(err)
	assert.IsType(unmarshalError, err)
	err = ba2.UnmarshalBinary([]byte{0x00})
	assert.NotNil(err)
	assert.IsType(unmarshalError, err)
	err = ba2.UnmarshalBinary([]byte{0x01, 0x00, 0x00, 0x00})
	assert.NotNil(err)
	assert.IsType(unmarshalError, err)
	err = ba2.UnmarshalBinary(buf)
	assert.NoError(err)
	assert.Equal(ba, ba2)
}
