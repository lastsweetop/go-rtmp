package model

type BasicHeader struct {
	Fmt           byte
	ChunkStreamId uint32
}

type MessageHeader struct {
	Timestamp       uint32
	PayloadLength   uint32
	TypeId          byte
	MessageStreamId uint32
	Headlen         uint32
}
