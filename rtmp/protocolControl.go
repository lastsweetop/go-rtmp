package rtmp

import (
	"go-rtmp/buffertool"
	"go-rtmp/enums/LimitTypes"
	"go-rtmp/model"
	"log"
)

var WindowAcknowledgementSize = 5000000
var MaxChunkSize = 450000

func (this *Session) WriteWindowAcknowledgementSize() {
	buffer := model.Chunk{
		BasicHeader: &model.BasicHeader{
			Fmt:           0,
			ChunkStreamId: 2,
		},
		MessageHeader: &model.MessageHeader{
			Timestamp:       0,
			TypeId:          0x05,
			MessageStreamId: 0,
		},
	}
	payload := make([]byte, 4)
	payload[0] = byte(WindowAcknowledgementSize >> 24)
	payload[1] = byte(WindowAcknowledgementSize >> 16)
	payload[2] = byte(WindowAcknowledgementSize >> 8)
	payload[3] = byte(WindowAcknowledgementSize)
	buffer.Payload = payload
	this.write(buffer.Bytes())
}

func (this *Session) WriteSetPeerBandwidth() {
	buffer := model.Chunk{
		BasicHeader: &model.BasicHeader{
			Fmt:           0,
			ChunkStreamId: 2,
		},
		MessageHeader: &model.MessageHeader{
			Timestamp:       0,
			TypeId:          0x06,
			MessageStreamId: 0,
		},
	}
	payload := make([]byte, 5)
	payload[0] = byte(WindowAcknowledgementSize >> 24)
	payload[1] = byte(WindowAcknowledgementSize >> 16)
	payload[2] = byte(WindowAcknowledgementSize >> 8)
	payload[3] = byte(WindowAcknowledgementSize)
	payload[4] = LimitTypes.Dynamic

	buffer.Payload = payload
	this.write(buffer.Bytes())
}

func (this *Session) WriteSetChunkSize() {
	buffer := model.Chunk{
		BasicHeader: &model.BasicHeader{
			Fmt:           0,
			ChunkStreamId: 2,
		},
		MessageHeader: &model.MessageHeader{
			Timestamp:       0,
			TypeId:          0x01,
			MessageStreamId: 0,
		},
	}
	payload := make([]byte, 4)
	payload[0] = byte(MaxChunkSize >> 24)
	payload[1] = byte(MaxChunkSize >> 16)
	payload[2] = byte(MaxChunkSize >> 8)
	payload[3] = byte(MaxChunkSize)
	buffer.Payload = payload
	this.write(buffer.Bytes())
}

func (this *Session) handWindowAcknowledgementSize(b *buffertool.Buffer, AMFVersion int) {
	size := b.ScanInt()

	this.WindowAcknowledgementSize = size

	log.Println("WindowAcknowledgementSize:", size)
}

func (this *Session) handAcknowledgement(b *buffertool.Buffer, AMFVersion int) {
	size := b.ScanInt()

	log.Println("Client Acknowledgement:", size)
}
