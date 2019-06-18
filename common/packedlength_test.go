package common_test

import (
	"github.com/newkedison/go-utils/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPackLength(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(common.PackLength(0), *common.NewByteArray([]byte{0x00}))
	assert.Equal(common.PackLength(1), *common.NewByteArray([]byte{0x01}))
	assert.Equal(common.PackLength(0x7F), *common.NewByteArray([]byte{0x7F}))
	assert.Equal(common.PackLength(0x80), *common.NewByteArray([]byte{0x80, 0x80}))
	assert.Equal(common.PackLength(0xFF), *common.NewByteArray([]byte{0x80, 0xFF}))
	assert.Equal(common.PackLength(0x100), *common.NewByteArray([]byte{0xC0, 0x01, 0x00}))
	assert.Equal(common.PackLength(0xFFF), *common.NewByteArray([]byte{0xC0, 0x0F, 0xFF}))
	assert.Equal(common.PackLength(0xFFFF), *common.NewByteArray([]byte{0xC0, 0xFF, 0xFF}))
	assert.Equal(common.PackLength(0x10000), *common.NewByteArray([]byte{0xE0, 0x01, 0x00, 0x00}))
	assert.Equal(common.PackLength(0xFFFFF), *common.NewByteArray([]byte{0xE0, 0x0F, 0xFF, 0xFF}))
	assert.Equal(common.PackLength(0x1FFFFF), *common.NewByteArray([]byte{0xE0, 0x1F, 0xFF, 0xFF}))
	assert.Equal(common.PackLength(0x200000), *common.NewByteArray([]byte{0xE0, 0x20, 0x00, 0x00}))
	assert.Equal(common.PackLength(0x2FFFFF), *common.NewByteArray([]byte{0xE0, 0x2F, 0xFF, 0xFF}))
	assert.Equal(common.PackLength(0xFFFFFF), *common.NewByteArray([]byte{0xE0, 0xFF, 0xFF, 0xFF}))
	assert.Equal(common.PackLength(0x1000000), *common.NewByteArray([]byte{0xF0, 0x01, 0x00, 0x00, 0x00}))
	assert.Equal(common.PackLength(0xFFFFFFF), *common.NewByteArray([]byte{0xF0, 0x0F, 0xFF, 0xFF, 0xFF}))
	assert.Equal(common.PackLength(0x10000000), *common.NewByteArray([]byte{0xF0, 0x10, 0x00, 0x00, 0x00}))
	assert.Equal(common.PackLength(0x1FFFFFFF), *common.NewByteArray([]byte{0xF0, 0x1F, 0xFF, 0xFF, 0xFF}))
	assert.Equal(common.PackLength(0xFFFFFFFE), *common.NewByteArray([]byte{0xF0, 0xFF, 0xFF, 0xFF, 0xFE}))
	assert.Equal(common.PackLength(0xFFFFFFFF), *common.NewByteArray([]byte{0xF0, 0xFF, 0xFF, 0xFF, 0xFF}))
}

func TestAppendPackedLength(t *testing.T) {
	assert := assert.New(t)
	var ba common.ByteArray
	common.AppendPackedLength(&ba, 0x7F)
	assert.Equal(ba, *common.NewByteArray([]byte{0x7F}))
	common.AppendPackedLength(&ba, 0x10000)
	assert.Equal(ba, *common.NewByteArray([]byte{0x7F, 0xE0, 0x01, 0x00, 0x00}))
	common.AppendPackedLength(&ba, 0x1FFFFFFF)
	assert.Equal(ba, *common.NewByteArray([]byte{0x7F, 0xE0, 0x01, 0x00, 0x00, 0xF0, 0x1F, 0xFF, 0xFF, 0xFF}))
}

func checkUnpackLengthValue(a *assert.Assertions, data []byte, expect uint32) {
	l, err := common.UnpackLength(common.ByteArray(data))
	a.Nil(err)
	a.Equal(l, expect)
}

func checkUnpackLengthError(a *assert.Assertions, data []byte) {
	l, err := common.UnpackLength(common.ByteArray(data))
	a.NotNil(err)
	a.Equal(l, uint32(0))
	a.Equal(err, common.InvalidPackedLengthError)
}

func TestUnpackLength(t *testing.T) {
	assert := assert.New(t)
	checkUnpackLengthValue(assert, []byte{0x00}, 0)
	checkUnpackLengthValue(assert, []byte{0x7F}, 0x7F)
	checkUnpackLengthValue(assert, []byte{0x80, 0x80}, 0x80)
	checkUnpackLengthValue(assert, []byte{0x80, 0xFF}, 0xFF)
	checkUnpackLengthValue(assert, []byte{0xC0, 0x01, 0x00}, 0x100)
	checkUnpackLengthValue(assert, []byte{0xC0, 0xFF, 0xFF}, 0xFFFF)
	checkUnpackLengthValue(assert, []byte{0xE0, 0x01, 0x00, 0x00}, 0x10000)
	checkUnpackLengthValue(assert, []byte{0xE0, 0xFF, 0xFF, 0xFF}, 0xFFFFFF)
	checkUnpackLengthValue(assert, []byte{0xF0, 0x01, 0x00, 0x00, 0x00}, 0x1000000)
	checkUnpackLengthValue(assert, []byte{0xF0, 0xFF, 0xFF, 0xFF, 0xFF}, 0xFFFFFFFF)

	checkUnpackLengthError(assert, []byte{})
	checkUnpackLengthError(assert, []byte{0x80})
	checkUnpackLengthError(assert, []byte{0xC0})
	checkUnpackLengthError(assert, []byte{0xC0, 0xFF})
	checkUnpackLengthError(assert, []byte{0xE0})
	checkUnpackLengthError(assert, []byte{0xE0, 0xFF})
	checkUnpackLengthError(assert, []byte{0xE0, 0xFF, 0xFF})
	checkUnpackLengthError(assert, []byte{0xF0})
	checkUnpackLengthError(assert, []byte{0xF0, 0xFF})
	checkUnpackLengthError(assert, []byte{0xF0, 0xFF, 0xFF})
	checkUnpackLengthError(assert, []byte{0xF0, 0xFF, 0xFF, 0xFF})
	checkUnpackLengthError(assert, []byte{0x80, 0x00})
	checkUnpackLengthError(assert, []byte{0x80, 0x7F})
	checkUnpackLengthValue(assert, []byte{0x80, 0x80}, 0x80)
	checkUnpackLengthError(assert, []byte{0xC0, 0x00, 0x00})
	checkUnpackLengthError(assert, []byte{0xC0, 0x00, 0xFF})
	checkUnpackLengthValue(assert, []byte{0xC0, 0x01, 0x00}, 0x100)
	checkUnpackLengthError(assert, []byte{0xE0, 0x00, 0x00, 0x00})
	checkUnpackLengthError(assert, []byte{0xE0, 0x00, 0xFF, 0xFF})
	checkUnpackLengthValue(assert, []byte{0xE0, 0x01, 0x00, 0x00}, 0x10000)
	checkUnpackLengthError(assert, []byte{0xF0, 0x00, 0x00, 0x00, 0x00})
	checkUnpackLengthError(assert, []byte{0xF0, 0x00, 0xFF, 0xFF, 0xFF})
	checkUnpackLengthValue(assert, []byte{0xF0, 0x01, 0x00, 0x00, 0x00}, 0x1000000)
}

func TestByteCountOfPackedLength(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(common.ByteCountOfPackedLength(0), 1)
	assert.Equal(common.ByteCountOfPackedLength(0x7F), 1)
	assert.Equal(common.ByteCountOfPackedLength(0x80), 2)
	assert.Equal(common.ByteCountOfPackedLength(0xFF), 2)
	assert.Equal(common.ByteCountOfPackedLength(0x100), 3)
	assert.Equal(common.ByteCountOfPackedLength(0xFFFF), 3)
	assert.Equal(common.ByteCountOfPackedLength(0x10000), 4)
	assert.Equal(common.ByteCountOfPackedLength(0xFFFFFF), 4)
	assert.Equal(common.ByteCountOfPackedLength(0x1000000), 5)
	assert.Equal(common.ByteCountOfPackedLength(0xFFFFFFFF), 5)
}
