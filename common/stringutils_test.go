package common_test

import (
	"encoding/hex"
	"github.com/newkedison/go-utils/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestByteToHexString(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(common.ByteToHexString(0), "00")
	assert.Equal(common.ByteToHexString(0xFF), "FF")
	assert.Equal(common.ByteToHexString(0x0F), "0F")
	assert.Equal(common.ByteToHexString(0x10), "10")
	assert.Equal(common.ByteToHexString(0xa5), "A5")
}

func testOneValidHexString(assert *assert.Assertions, in string, expect byte) {
	b, err := common.HexStringToByte(in)
	assert.Nil(err)
	assert.Equal(b, expect)
}

func TestHexStringToByte(t *testing.T) {
	assert := assert.New(t)
	testOneValidHexString(assert, "00", 0)
	testOneValidHexString(assert, "FF", 0xFF)
	testOneValidHexString(assert, "55", 0x55)
	testOneValidHexString(assert, "A5", 0xA5)

	_, err := common.HexStringToByte("000")
	assert.Equal(err, hex.ErrLength)
	_, err = common.HexStringToByte("ZZ")
	assert.EqualError(err, "encoding/hex: invalid byte: U+005A 'Z'")
	_, err = common.HexStringToByte("0Z")
	assert.EqualError(err, "encoding/hex: invalid byte: U+005A 'Z'")
}
