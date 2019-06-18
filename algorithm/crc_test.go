package algorithm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var testData = makeTestData(256)

func makeTestData(n int) []byte {
	ret := make([]byte, n, n)
	for i := 0; i < n; i++ {
		ret[i] = byte(i)
	}
	return ret
}

func TestCrc8(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(Crc8([]byte{0x10}), byte(0x70))
	assert.Equal(Crc8([]byte{0x10, 0x20, 0x30, 0xFF}), byte(0x2E))
	assert.Equal(Crc8(testData), byte(0x14))
}

func TestCrc8Impl(t *testing.T) {
	assert := assert.New(t)
	table := MakeCrc8Table(0x1D)
	assert.Equal(
		Crc8Impl([]byte{0x10}, table, 0xFD, false, false, 0x00),
		byte(0x33))
	assert.Equal(
		Crc8Impl(testData, table, 0xFD, false, true, 0x00),
		byte(0x03))
}

func TestCrc8Predefined(t *testing.T) {
	assert := assert.New(t)
	if _, err := Crc8Predefined(nil, 100); assert.NotNil(err) {
		assert.EqualError(err, "core/algorithm/crc: invalid index: 100")
	}
	crc, _ := Crc8Predefined(testData, 0)
	assert.EqualValues(crc, 0x14)
	crc, _ = Crc8Predefined(testData, 1)
	assert.EqualValues(crc, 0x41)
	crc, _ = Crc8Predefined(testData, 2)
	assert.EqualValues(crc, 0x3C)
	crc, _ = Crc8Predefined(testData, 3)
	assert.EqualValues(crc, 0xCA)
	crc, _ = Crc8Predefined(testData, 4)
	assert.EqualValues(crc, 0xC5)
	crc, _ = Crc8Predefined(testData, 5)
	assert.EqualValues(crc, 0xC0)
	crc, _ = Crc8Predefined(testData, 6)
	assert.EqualValues(crc, 0x41)
	crc, _ = Crc8Predefined(testData, 7)
	assert.EqualValues(crc, 0x18)
	crc, _ = Crc8Predefined(testData, 8)
	assert.EqualValues(crc, 0x8E)
	crc, _ = Crc8Predefined(testData, 9)
	assert.EqualValues(crc, 0x59)
}

func TestCrc8Continue(t *testing.T) {
	crc := Crc8(testData[:100])
	crc, _ = Crc8Continue(testData[100:], crc)
	assert.Equal(t, crc, byte(0x14))
	crc = Crc8(testData[:150])
	crc, _ = Crc8Continue(testData[150:], crc)
	assert.Equal(t, crc, byte(0x14))
}

func TestAppendCrc8(t *testing.T) {
	assert := assert.New(t)
	data := append([]byte(nil), testData...)
	data = AppendCrc8(data)
	assert.Equal(len(data), 257)
	assert.Equal(data[256], byte(0x14))
}

func TestVerifyCrc8(t *testing.T) {
	assert := assert.New(t)
	data := append([]byte(nil), testData...)
	data = AppendCrc8(data)
	assert.True(VerifyCrc8(data))
	data[256] ^= 0xFF
	assert.False(VerifyCrc8(data))
}

func TestCrc16(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(Crc16([]byte{0x10}), uint16(0x8CBE))
	assert.Equal(Crc16(testData), uint16(0xDE6C))
}

func TestCrc16Impl(t *testing.T) {
	assert := assert.New(t)
	table := MakeCrc16Table(0x1021)
	assert.Equal(
		Crc16Impl([]byte{0x10}, table, 0x0000, false, false, 0x0000),
		uint16(0x1231))
	assert.Equal(
		Crc16Impl(testData, table, 0xFFFF, false, false, 0x00),
		uint16(0x3FBD))
}

func TestCrc16Predefined(t *testing.T) {
	assert := assert.New(t)
	if _, err := Crc16Predefined(nil, 100); assert.NotNil(err) {
		assert.EqualError(err, "core/algorithm/crc: invalid index: 100")
	}
	crc, _ := Crc16Predefined(testData, 0)
	assert.EqualValues(crc, 0x3FBD)
	crc, _ = Crc16Predefined(testData, 1)
	assert.EqualValues(crc, 0xBAD3)
	crc, _ = Crc16Predefined(testData, 2)
	assert.EqualValues(crc, 0x3C8E)
	crc, _ = Crc16Predefined(testData, 3)
	assert.EqualValues(crc, 0x3B7A)
	crc, _ = Crc16Predefined(testData, 4)
	assert.EqualValues(crc, 0xF8E6)
	crc, _ = Crc16Predefined(testData, 5)
	assert.EqualValues(crc, 0xB5A1)
	crc, _ = Crc16Predefined(testData, 6)
	assert.EqualValues(crc, 0x9A1C)
	crc, _ = Crc16Predefined(testData, 7)
	assert.EqualValues(crc, 0x9A1D)
	crc, _ = Crc16Predefined(testData, 8)
	assert.EqualValues(crc, 0x4472)
	crc, _ = Crc16Predefined(testData, 9)
	assert.EqualValues(crc, 0xB50D)
	crc, _ = Crc16Predefined(testData, 10)
	assert.EqualValues(crc, 0xC042)
	crc, _ = Crc16Predefined(testData, 11)
	assert.EqualValues(crc, 0x452C)
	crc, _ = Crc16Predefined(testData, 12)
	assert.EqualValues(crc, 0xCFC3)
	crc, _ = Crc16Predefined(testData, 13)
	assert.EqualValues(crc, 0x563B)
	crc, _ = Crc16Predefined(testData, 14)
	assert.EqualValues(crc, 0xE0B5)
	crc, _ = Crc16Predefined(testData, 15)
	assert.EqualValues(crc, 0x2193)
	crc, _ = Crc16Predefined(testData, 16)
	assert.EqualValues(crc, 0xD841)
	crc, _ = Crc16Predefined(testData, 17)
	assert.EqualValues(crc, 0xDE6C)
	crc, _ = Crc16Predefined(testData, 18)
	assert.EqualValues(crc, 0x303C)
	crc, _ = Crc16Predefined(testData, 19)
	assert.EqualValues(crc, 0x7E55)
}

func TestAppendCrc16(t *testing.T) {
	assert := assert.New(t)
	data := append([]byte(nil), testData...)
	data = AppendCrc16(data)
	assert.Equal(len(data), 258)
	assert.Equal(data[256], byte(0x6C))
	assert.Equal(data[257], byte(0xDE))
}

func TestCrc16Continue(t *testing.T) {
	crc := Crc16(testData[:100])
	crc, _ = Crc16Continue(testData[100:], crc)
	assert.Equal(t, crc, uint16(0xDE6C))
	crc = Crc16(testData[:150])
	crc, _ = Crc16Continue(testData[150:], crc)
	assert.Equal(t, crc, uint16(0xDE6C))
}

func TestVerifyCrc16(t *testing.T) {
	assert := assert.New(t)
	data := append([]byte(nil), testData...)
	data = AppendCrc16(data)
	assert.True(VerifyCrc16(data))
	data[256] ^= 0xFF
	assert.False(VerifyCrc16(data))
}

func TestCrc32(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(Crc32([]byte{0x10}), uint32(0xCFB5FFE9))
	assert.Equal(Crc32(testData), uint32(0x29058C73))
}

func TestCrc32Impl(t *testing.T) {
	assert := assert.New(t)
	table := MakeCrc32Table(0xA833982B)
	assert.Equal(
		Crc32Impl([]byte{0x10}, table, 0xFFFFFFFF, false, false, 0xFFFFFFFF),
		uint32(0x0E44E111))
	assert.Equal(
		Crc32Impl(testData, table, 0xFFFFFFFF, true, true, 0xFFFFFFFF),
		uint32(0x6165003E))
}

func TestCrc32Predefined(t *testing.T) {
	assert := assert.New(t)
	if _, err := Crc32Predefined(nil, 100); assert.NotNil(err) {
		assert.EqualError(err, "core/algorithm/crc: invalid index: 100")
	}
	crc, _ := Crc32Predefined(testData, 0)
	assert.EqualValues(crc, 0x29058C73)
	crc, _ = Crc32Predefined(testData, 1)
	assert.EqualValues(crc, 0xB6B5EE95)
	crc, _ = Crc32Predefined(testData, 2)
	assert.EqualValues(crc, 0x9C44184B)
	crc, _ = Crc32Predefined(testData, 3)
	assert.EqualValues(crc, 0x6165003E)
	crc, _ = Crc32Predefined(testData, 4)
	assert.EqualValues(crc, 0x494A116A)
	crc, _ = Crc32Predefined(testData, 5)
	assert.EqualValues(crc, 0x53EB78DA)
	crc, _ = Crc32Predefined(testData, 6)
	assert.EqualValues(crc, 0xAB16EA85)
	crc, _ = Crc32Predefined(testData, 7)
	assert.EqualValues(crc, 0xD6FA738C)
	crc, _ = Crc32Predefined(testData, 8)
	assert.EqualValues(crc, 0x3FF2756B)
}

func TestAppendCrc32(t *testing.T) {
	assert := assert.New(t)
	data := append([]byte(nil), testData...)
	data = AppendCrc32(data)
	assert.Equal(len(data), 260)
	assert.Equal(data[256], byte(0x73))
	assert.Equal(data[257], byte(0x8C))
	assert.Equal(data[258], byte(0x05))
	assert.Equal(data[259], byte(0x29))
}

func TestCrc32Continue(t *testing.T) {
	crc := Crc32(testData[:100])
	crc, _ = Crc32Continue(testData[100:], crc)
	assert.Equal(t, crc, uint32(0x29058C73))
	crc = Crc32(testData[:150])
	crc, _ = Crc32Continue(testData[150:], crc)
	assert.Equal(t, crc, uint32(0x29058C73))
}

func TestVerifyCrc32(t *testing.T) {
	assert := assert.New(t)
	data := append([]byte(nil), testData...)
	data = AppendCrc32(data)
	assert.True(VerifyCrc32(data))
	data[256] ^= 0xFF
	assert.False(VerifyCrc32(data))
}

func TestChangeCrc8ConfigIndex(t *testing.T) {
	assert := assert.New(t)
	DefaultCrc8ConfigIndex = 100
	assert.Panics(func() { Crc8([]byte{0x00}) })
	DefaultCrc16ConfigIndex = 100
	assert.Panics(func() { Crc16([]byte{0x00}) })
	DefaultCrc32ConfigIndex = 100
	assert.Panics(func() { Crc32([]byte{0x00}) })

	_, err := Crc8Continue([]byte{0x10}, 0)
	assert.Error(err)
	_, err = Crc16Continue([]byte{0x10}, 0)
	assert.Error(err)
	_, err = Crc32Continue([]byte{0x10}, 0)
	assert.Error(err)

	Crc8Configs = append(Crc8Configs, Config{0x9B, 0x00, 0x00, true, true})
	DefaultCrc8ConfigIndex = len(Crc8Configs) - 1
	if v, err := Crc8Continue([]byte{0x10},
		byte(Crc8Configs[DefaultCrc8ConfigIndex].InitValue)); assert.Nil(err) {
		assert.EqualValues(v, 0x98)
	}

	Crc16Configs = append(
		Crc16Configs, Config{0x1021, 0x0000, 0x0000, false, false})
	DefaultCrc16ConfigIndex = len(Crc16Configs) - 1
	if v, err := Crc16Continue([]byte{0x10}, uint16(
		Crc16Configs[DefaultCrc16ConfigIndex].InitValue)); assert.Nil(err) {
		assert.EqualValues(v, 0x1231)
	}

	Crc32Configs = append(
		Crc32Configs, Config{0xAF, 0, 0, false, false})
	DefaultCrc32ConfigIndex = len(Crc32Configs) - 1
	if v, err := Crc32Continue([]byte{0x10}, uint32(
		Crc32Configs[DefaultCrc32ConfigIndex].InitValue)); assert.Nil(err) {
		assert.EqualValues(v, 0xAF0)
	}
}
