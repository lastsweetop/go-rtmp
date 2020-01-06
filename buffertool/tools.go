package buffertool

import (
	"encoding/binary"
	"math"
)

type Buffer struct {
	data   []byte
	offset uint32
}

func NewBuffer(data []byte) *Buffer {
	return &Buffer{
		data:   data,
		offset: 0,
	}
}

func (i *Buffer) ScanByteForNBit(N byte) byte {
	nByte := i.offset / 8
	fbit := i.offset % 8
	temp := 1<<(8-byte(fbit)) - 1 - (1<<(8-byte(fbit)-N) - 1)
	i.offset += uint32(N)
	return (i.data[nByte] & byte(temp)) >> (8 - byte(fbit) - N)
}

func (i *Buffer) ScanByte() byte {
	i.offset += 8
	return i.data[i.offset/8-1]
}

func (i *Buffer) ScanInt16() int {
	i.offset += 16
	return int(i.data[i.offset/8-2])<<8 + int(i.data[i.offset/8-1])
}

func (i *Buffer) ScanInt() int {
	i.offset += 32
	return int(i.data[i.offset/8-4])<<24 + int(i.data[i.offset/8-3])<<16 + int(i.data[i.offset/8-2])<<8 + int(i.data[i.offset/8-1])
}

func (i *Buffer) ScanBytes(n uint32) []byte {
	i.offset += 8 * n
	return i.data[i.offset/8-n : i.offset/8]
}

func (i *Buffer) TryBytes(n uint32) []byte {
	return i.data[i.offset/8 : i.offset/8+n]
}

func (i *Buffer) ScanAMFString() string {
	length := uint32(i.ScanByte())<<8 + uint32(i.ScanByte())
	return string(i.ScanBytes(length))
}

func (i *Buffer) ScanAMFNumber() float64 {
	return math.Float64frombits(binary.BigEndian.Uint64(i.ScanBytes(8)))
}

func (i *Buffer) ScanAMFBoolean() bool {
	return i.ScanByte() != 0
}

func Merge(s1 []byte, s2 []byte) []byte {
	slice := make([]byte, len(s1)+len(s2))
	copy(slice, s1)
	copy(slice[len(s1):], s2)
	return slice
}
