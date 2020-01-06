package rtmp

import (
	"go-rtmp/amf"
	"go-rtmp/buffertool"
	"go-rtmp/model"
	"log"
)

type CreateStreamResult struct {
	Name          string
	TransactionID float64
	CommandObject interface{}
	StreamId      float64
}

func (this *Session) handCreateStream(b *buffertool.Buffer, AMFVersion int) {
	_ = b.ScanByte()
	tid := b.ScanAMFNumber()
	log.Println("transaction Id", tid)

	buffer := model.Chunk{
		BasicHeader: &model.BasicHeader{
			Fmt:           0,
			ChunkStreamId: this.ChunkStreamId,
		},
		MessageHeader: &model.MessageHeader{
			Timestamp:       0,
			TypeId:          0x14,
			MessageStreamId: 0,
		},
	}
	cs := CreateStreamResult{
		Name:          "_result",
		TransactionID: tid,
		CommandObject: nil,
		StreamId:      1,
	}
	payload, _ := amf.Marshal(cs)
	buffer.Payload = payload
	this.write(buffer.Bytes())
}

func (this *Session) deleteStream(b *buffertool.Buffer, AMFVersion int) {

}
