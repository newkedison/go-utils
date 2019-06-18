package algorithm

import (
	"fmt"
	"math/bits"
)

type Config struct {
	Polynomial uint32
	InitValue  uint32
	XorResult  uint32
	ReflectIn  bool
	ReflectOut bool
}

// Reference:
// http://crccalc.com/
// http://www.sunshine2k.de/coding/javascript/crc/crc_js.html
// http://www.sunshine2k.de/articles/coding/crc/understanding_crc.html
var (
	Crc8Configs = []Config{
		{0x07, 0x00, 0x00, false, false}, // CRC-8
		{0x9B, 0xFF, 0x00, false, false}, // CRC-8/CDMA2000
		{0x39, 0x00, 0x00, true, true},   // CRC-8/DARC
		{0xD5, 0x00, 0x00, false, false}, // CRC-8/DVB-S2
		{0x1D, 0xFF, 0x00, true, true},   // CRC-8/EBU
		{0x1D, 0xFD, 0x00, false, false}, // CRC-8/I-CODE
		{0x07, 0x00, 0x55, false, false}, // CRC-8/ITU
		{0x31, 0x00, 0x00, true, true},   // CRC-8/MAXIM
		{0x07, 0xFF, 0x00, true, true},   // CRC-8/ROHC
		{0x9B, 0x00, 0x00, true, true},   // CRC-8/WCDMA
	}
	Crc16Configs = []Config{
		{0x1021, 0xFFFF, 0x0000, false, false}, // CRC-16/CCITT-FALSE
		{0x8005, 0x0000, 0x0000, true, true},   // CRC-16/ARC
		{0x1021, 0x1D0F, 0x0000, false, false}, // CRC-16/AUG-CCITT
		{0x8005, 0x0000, 0x0000, false, false}, // CRC-16/BUYPASS
		{0xC867, 0xFFFF, 0x0000, false, false}, // CRC-16/CDMA2000
		{0x8005, 0x800D, 0x0000, false, false}, // CRC-16/DDS-110
		{0x0589, 0x0000, 0x0001, false, false}, // CRC-16/DECT-R
		{0x0589, 0x0000, 0x0000, false, false}, // CRC-16/DECT-X
		{0x3D65, 0x0000, 0xFFFF, true, true},   // CRC-16/DNP
		{0x3D65, 0x0000, 0xFFFF, false, false}, // CRC-16/EN-13757
		{0x1021, 0xFFFF, 0xFFFF, false, false}, // CRC-16/GENIBUS
		{0x8005, 0x0000, 0xFFFF, true, true},   // CRC-16/MAXIM
		{0x1021, 0xFFFF, 0x0000, true, true},   // CRC-16/MCRF4XX
		{0x8BB7, 0x0000, 0x0000, false, false}, // CRC-16/T10-DIF
		{0xA097, 0x0000, 0x0000, false, false}, // CRC-16/TELEDISK
		{0x8005, 0xFFFF, 0xFFFF, true, true},   // CRC-16/USB
		{0x1021, 0x0000, 0x0000, true, true},   // CRC-16/KERMIT
		{0x8005, 0xFFFF, 0x0000, true, true},   // CRC-16/MODBUS
		{0x1021, 0xFFFF, 0xFFFF, true, true},   // CRC-16/X-25
		{0x1021, 0x0000, 0x0000, false, false}, // CRC-16/XMODEM
	}
	Crc32Configs = []Config{
		{0x04C11DB7, 0xFFFFFFFF, 0xFFFFFFFF, true, true},   // CRC-32
		{0x04C11DB7, 0xFFFFFFFF, 0xFFFFFFFF, false, false}, // CRC-32/BZIP2
		{0x1EDC6F41, 0xFFFFFFFF, 0xFFFFFFFF, true, true},   // CRC-32C
		{0xA833982B, 0xFFFFFFFF, 0xFFFFFFFF, true, true},   // CRC-32D
		{0x04C11DB7, 0xFFFFFFFF, 0x00000000, false, false}, // CRC-32/MPEG-2
		{0x04C11DB7, 0x00000000, 0xFFFFFFFF, false, false}, // CRC-32/POSIX
		{0x814141AB, 0x00000000, 0x00000000, false, false}, // CRC-32Q
		{0x04C11DB7, 0xFFFFFFFF, 0x00000000, true, true},   // CRC-32/JAMCRC
		{0x000000AF, 0x00000000, 0x00000000, false, false}, // CRC-32/XFER
	}
)

var (
	DefaultCrc8ConfigIndex  = 0
	DefaultCrc16ConfigIndex = 17
	DefaultCrc32ConfigIndex = 0
)

var (
	crc8Tables  = make([][]byte, len(Crc8Configs))
	crc16Tables = make([][]uint16, len(Crc16Configs))
	crc32Tables = make([][]uint32, len(Crc32Configs))
)

type InvalidIndexError int

func (e InvalidIndexError) Error() string {
	return fmt.Sprintf("core/algorithm/crc: invalid index: %d", int(e))
}

func MakeCrc8Table(poly byte) []byte {
	tbl := make([]byte, 256, 256)
	for i := 0; i < 256; i++ {
		crc := byte(i)
		for j := 0; j < 8; j++ {
			if crc&0x80 > 0 {
				crc = (crc << 1) ^ poly
			} else {
				crc <<= 1
			}
		}
		tbl[i] = crc
	}
	return tbl
}

func Crc8Impl(data []byte, table []byte, init byte,
	refin bool, refout bool, xorout byte) byte {
	crc := init
	for _, d := range data {
		if refin {
			d = bits.Reverse8(d)
		}
		crc = table[crc^d]
	}
	if refout {
		crc = bits.Reverse8(crc)
	}
	return crc ^ xorout
}

func checkCrc8ConfigIndex(index int) bool {
	if index >= len(Crc8Configs) {
		return false
	}
	for len(crc8Tables) <= index {
		crc8Tables = append(crc8Tables, nil)
	}
	if crc8Tables[index] == nil {
		crc8Tables[index] = MakeCrc8Table(
			byte(Crc8Configs[index].Polynomial))
	}
	return true
}

func Crc8Predefined(data []byte, configIndex int) (byte, error) {
	if !checkCrc8ConfigIndex(configIndex) {
		return 0, InvalidIndexError(configIndex)
	}
	return Crc8Impl(data, crc8Tables[configIndex],
		byte(Crc8Configs[configIndex].InitValue),
		Crc8Configs[configIndex].ReflectIn,
		Crc8Configs[configIndex].ReflectOut,
		byte(Crc8Configs[configIndex].XorResult)), nil
}

func Crc8(data []byte) byte {
	crc, err := Crc8Predefined(data, DefaultCrc8ConfigIndex)
	if err != nil {
		panic("CRC8 fail")
	}
	return crc
}

func Crc8Continue(data []byte, prev byte) (byte, error) {
	configIndex := DefaultCrc8ConfigIndex
	if !checkCrc8ConfigIndex(configIndex) {
		return 0, InvalidIndexError(configIndex)
	}
	if Crc8Configs[configIndex].ReflectOut {
		prev = bits.Reverse8(prev)
	}
	prev ^= byte(Crc8Configs[configIndex].XorResult)
	return Crc8Impl(data, crc8Tables[configIndex],
		prev,
		Crc8Configs[configIndex].ReflectIn,
		Crc8Configs[configIndex].ReflectOut,
		byte(Crc8Configs[configIndex].XorResult)), nil
}

func AppendCrc8(data []byte) []byte {
	crc := Crc8(data)
	return append(data, crc)
}

func VerifyCrc8(data []byte) bool {
	return Crc8(data) == 0
}

func MakeCrc16Table(poly uint16) []uint16 {
	tbl := make([]uint16, 256, 256)
	for i := 0; i < 256; i++ {
		crc := uint16(i << 8)
		for j := 0; j < 8; j++ {
			if crc&0x8000 > 0 {
				crc = (crc << 1) ^ poly
			} else {
				crc <<= 1
			}
		}
		tbl[i] = crc
	}
	return tbl
}

func Crc16Impl(data []byte, table []uint16, init uint16,
	refin bool, refout bool, xorout uint16) uint16 {
	crc := init
	for _, d := range data {
		if refin {
			d = bits.Reverse8(d)
		}
		crc = uint16(crc<<8) ^ table[byte(crc>>8)^d]
	}
	if refout {
		crc = bits.Reverse16(crc)
	}
	return crc ^ xorout
}

func checkCrc16ConfigIndex(index int) bool {
	if index >= len(Crc16Configs) {
		return false
	}
	for len(crc16Tables) <= index {
		crc16Tables = append(crc16Tables, nil)
	}
	if crc16Tables[index] == nil {
		crc16Tables[index] = MakeCrc16Table(
			uint16(Crc16Configs[index].Polynomial))
	}
	return true
}

func Crc16Predefined(data []byte, configIndex int) (uint16, error) {
	if !checkCrc16ConfigIndex(configIndex) {
		return 0, InvalidIndexError(configIndex)
	}
	return Crc16Impl(data, crc16Tables[configIndex],
		uint16(Crc16Configs[configIndex].InitValue),
		Crc16Configs[configIndex].ReflectIn,
		Crc16Configs[configIndex].ReflectOut,
		uint16(Crc16Configs[configIndex].XorResult)), nil
}

func Crc16(data []byte) uint16 {
	crc, err := Crc16Predefined(data, DefaultCrc16ConfigIndex)
	if err != nil {
		panic("CRC16 fail")
	}
	return crc
}

func Crc16Continue(data []byte, prev uint16) (uint16, error) {
	configIndex := DefaultCrc16ConfigIndex
	if !checkCrc16ConfigIndex(configIndex) {
		return 0, InvalidIndexError(configIndex)
	}
	if Crc16Configs[configIndex].ReflectOut {
		prev = bits.Reverse16(prev)
	}
	prev ^= uint16(Crc16Configs[configIndex].XorResult)
	return Crc16Impl(data, crc16Tables[configIndex],
		prev,
		Crc16Configs[configIndex].ReflectIn,
		Crc16Configs[configIndex].ReflectOut,
		uint16(Crc16Configs[configIndex].XorResult)), nil
}

func AppendCrc16(data []byte) []byte {
	crc := Crc16(data)
	return append(data, []byte{byte(crc), byte(crc >> 8)}...)
}

func VerifyCrc16(data []byte) bool {
	return Crc16(data) == 0
}

func MakeCrc32Table(poly uint32) []uint32 {
	tbl := make([]uint32, 256, 256)
	for i := 0; i < 256; i++ {
		crc := uint32(i << 24)
		for j := 0; j < 8; j++ {
			if crc&0x80000000 > 0 {
				crc = (crc << 1) ^ poly
			} else {
				crc <<= 1
			}
		}
		tbl[i] = crc
	}
	return tbl
}

func Crc32Impl(data []byte, table []uint32, init uint32,
	refin bool, refout bool, xorout uint32) uint32 {
	crc := init
	for _, d := range data {
		if refin {
			d = bits.Reverse8(d)
		}
		crc = uint32(crc<<8) ^ table[byte(crc>>24)^d]
	}
	if refout {
		crc = bits.Reverse32(crc)
	}
	return crc ^ xorout
}

func checkCrc32ConfigIndex(index int) bool {
	if index >= len(Crc32Configs) {
		return false
	}
	for len(crc32Tables) <= index {
		crc32Tables = append(crc32Tables, nil)
	}
	if crc32Tables[index] == nil {
		crc32Tables[index] = MakeCrc32Table(
			uint32(Crc32Configs[index].Polynomial))
	}
	return true
}

func Crc32Predefined(data []byte, configIndex int) (uint32, error) {
	if !checkCrc32ConfigIndex(configIndex) {
		return 0, InvalidIndexError(configIndex)
	}
	return Crc32Impl(data, crc32Tables[configIndex],
		uint32(Crc32Configs[configIndex].InitValue),
		Crc32Configs[configIndex].ReflectIn,
		Crc32Configs[configIndex].ReflectOut,
		uint32(Crc32Configs[configIndex].XorResult)), nil
}

func Crc32(data []byte) uint32 {
	crc, err := Crc32Predefined(data, DefaultCrc32ConfigIndex)
	if err != nil {
		panic("CRC32 fail")
	}
	return crc
}

func Crc32Continue(data []byte, prev uint32) (uint32, error) {
	configIndex := DefaultCrc32ConfigIndex
	if !checkCrc32ConfigIndex(configIndex) {
		return 0, InvalidIndexError(configIndex)
	}
	if Crc32Configs[configIndex].ReflectOut {
		prev = bits.Reverse32(prev)
	}
	prev ^= uint32(Crc32Configs[configIndex].XorResult)
	return Crc32Impl(data, crc32Tables[configIndex],
		prev,
		Crc32Configs[configIndex].ReflectIn,
		Crc32Configs[configIndex].ReflectOut,
		uint32(Crc32Configs[configIndex].XorResult)), nil
}

func AppendCrc32(data []byte) []byte {
	crc := Crc32(data)
	return append(data,
		[]byte{byte(crc), byte(crc >> 8), byte(crc >> 16), byte(crc >> 24)}...)
}

func VerifyCrc32(data []byte) bool {
	return Crc32(data) == 0x2144DF1C
}
