package common

import (
	"encoding/hex"
	"strings"
)

func ByteToHexString(b byte) string {
	return strings.ToUpper(hex.EncodeToString([]byte{b}))
}

func HexStringToByte(s string) (byte, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return 0, err
	}
	return b[0], err
}
