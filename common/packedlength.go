package common

import (
	"errors"
)

var InvalidPackedLengthError error = errors.New("Invalid format for PackedLength")

func PackLength(length uint32) ByteArray {
	var ba ByteArray
	AppendPackedLength(&ba, length)
	return ba
}

func AppendPackedLength(dest *ByteArray, length uint32) {
	if length <= 0x7F {
		dest.AppendByte(byte(length & 0xFF))
	} else if length <= 0xFF {
		dest.AppendByte(0x80)
		dest.AppendByte(byte(length & 0xFF))
	} else if length <= 0xFFFF {
		dest.AppendByte(0xC0)
		dest.AppendByte(byte((length >> 8) & 0xFF))
		dest.AppendByte(byte(length & 0xFF))
	} else if length <= 0xFFFFFF {
		dest.AppendByte(0xE0)
		dest.AppendByte(byte((length >> 16) & 0xFF))
		dest.AppendByte(byte((length >> 8) & 0xFF))
		dest.AppendByte(byte(length & 0xFF))
	} else {
		dest.AppendByte(0xF0)
		dest.AppendByte(byte((length >> 24) & 0xFF))
		dest.AppendByte(byte((length >> 16) & 0xFF))
		dest.AppendByte(byte((length >> 8) & 0xFF))
		dest.AppendByte(byte(length & 0xFF))
	}
}

func UnpackLength(data ByteArray) (uint32, error) {
	if len(data) == 0 {
		return 0, InvalidPackedLengthError
	}
	if data[0] <= 0x7F {
		return uint32(data[0]), nil
	}
	length := len(data)
	if length > 1 && data[0] == 0x80 && data[1] > 0x7F {
		return uint32(data[1]), nil
	}
	if length > 2 && data[0] == 0xC0 && data[1] > 0 {
		return (uint32(data[1]) << 8) + uint32(data[2]), nil
	}
	if length > 3 && data[0] == 0xE0 && data[1] > 0 {
		return ((uint32)(data[1]) << 16) +
			(uint32(data[2]) << 8) +
			uint32(data[3]), nil
	}
	if length > 4 && data[0] == 0xF0 && data[1] > 0 {
		return ((uint32)(data[1]) << 24) +
			(uint32(data[2]) << 16) +
			(uint32(data[3]) << 8) +
			uint32(data[4]), nil
	}
	return 0, InvalidPackedLengthError
}

func ByteCountOfPackedLength(length uint32) int {
	if length <= 0x7F {
		return 1
	}
	if length <= 0xFF {
		return 2
	}
	if length <= 0xFFFF {
		return 3
	}
	if length <= 0xFFFFFF {
		return 4
	}
	return 5
}
