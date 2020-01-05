package rtmp

import (
	"go-rtmp/buffertool"
	"go-rtmp/model"
	"log"
)

func (this *Session) handUserControl(b *buffertool.Buffer, AMFVersion int) {
	eventType := b.ScanInt16()
	log.Println("user control message event type: ", eventType)
	switch eventType {
	case 0:
		break
	case 1:
		break
	case 2:
		break
	case 3:
		handSetBufferLength(b, AMFVersion)
		break
	case 4:
		break
	case 6:
		break
	case 7:
		break
	default:
		log.Println("error event type")
		break
	}
}

func handSetBufferLength(b *buffertool.Buffer, AMFVersion int) {
	streamId := b.ScanInt()
	bufferLength := b.ScanInt()
	log.Println("streamId:", streamId, "bufferLength:", bufferLength)
}

func (this *Session) StreamBegin() {
	buffer := model.Chunk{
		BasicHeader: &model.BasicHeader{
			Fmt:           0,
			ChunkStreamId: 2,
		},
		MessageHeader: &model.MessageHeader{
			Timestamp:       0,
			TypeId:          0x04,
			MessageStreamId: 0,
		},
	}
	buffer.Payload = []byte{0x00, 0x00,
		0x00, 0x00, 0x00, 0x01}
	this.write(buffer.Bytes())
}
