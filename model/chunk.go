package model

import (
	"bytes"
)

type Chunk struct {
	*BasicHeader
	*MessageHeader
	Payload []byte
}

func (c *Chunk) Bytes() []byte {
	buffer := bytes.NewBuffer([]byte{})
	if c.ChunkStreamId < 64 {
		buffer.WriteByte(c.Fmt<<6 + byte(c.ChunkStreamId))
	} else if c.ChunkStreamId < 319 {
		buffer.WriteByte(c.Fmt << 6)
		buffer.WriteByte(byte(c.ChunkStreamId - 64))
	} else {
		buffer.WriteByte(c.Fmt<<6 + 1)
		buffer.WriteByte(byte((c.ChunkStreamId - 64) >> 8))
		buffer.WriteByte(byte(c.ChunkStreamId - 64))
	}
	if c.Fmt < 3 {
		buffer.WriteByte(byte(c.Timestamp >> 16))
		buffer.WriteByte(byte(c.Timestamp >> 8))
		buffer.WriteByte(byte(c.Timestamp))
	}
	if c.Fmt < 2 {
		buffer.WriteByte(byte(len(c.Payload) >> 16))
		buffer.WriteByte(byte(len(c.Payload) >> 8))
		buffer.WriteByte(byte(len(c.Payload)))
		buffer.WriteByte(c.TypeId)

	}
	if c.Fmt < 1 {
		buffer.WriteByte(byte(c.MessageStreamId))
		buffer.WriteByte(byte(c.MessageStreamId >> 8))
		buffer.WriteByte(byte(c.MessageStreamId >> 16))
		buffer.WriteByte(byte(c.MessageStreamId >> 24))
	}
	buffer.Write(c.Payload)
	return buffer.Bytes()
}

func (c *Chunk) BytesWithPayloadLen(length int) []byte {
	buffer := bytes.NewBuffer([]byte{})
	if c.ChunkStreamId < 64 {
		buffer.WriteByte(c.Fmt<<6 + byte(c.ChunkStreamId))
	} else if c.ChunkStreamId < 319 {
		buffer.WriteByte(c.Fmt << 6)
		buffer.WriteByte(byte(c.ChunkStreamId - 64))
	} else {
		buffer.WriteByte(c.Fmt<<6 + 1)
		buffer.WriteByte(byte((c.ChunkStreamId - 64) >> 8))
		buffer.WriteByte(byte(c.ChunkStreamId - 64))
	}
	if c.Fmt < 3 {
		buffer.WriteByte(byte(c.Timestamp >> 16))
		buffer.WriteByte(byte(c.Timestamp >> 8))
		buffer.WriteByte(byte(c.Timestamp))
	}
	if c.Fmt < 2 {
		buffer.WriteByte(byte(length >> 16))
		buffer.WriteByte(byte(length >> 8))
		buffer.WriteByte(byte(length))
		buffer.WriteByte(c.TypeId)

	}
	if c.Fmt < 1 {
		buffer.WriteByte(byte(c.MessageStreamId))
		buffer.WriteByte(byte(c.MessageStreamId >> 8))
		buffer.WriteByte(byte(c.MessageStreamId >> 16))
		buffer.WriteByte(byte(c.MessageStreamId >> 24))
	}
	buffer.Write(c.Payload)
	return buffer.Bytes()
}
