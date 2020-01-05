package rtmp

import (
	"go-rtmp/amf"
	"go-rtmp/model"
)

type RtmpSampleAccessResponse struct {
	Name  string
	Audio bool
	video bool
}

func (this *Session) RtmpSampleAccess() {
	buffer := model.Chunk{
		BasicHeader: &model.BasicHeader{
			Fmt:           0,
			ChunkStreamId: 5,
		},
		MessageHeader: &model.MessageHeader{
			Timestamp:       0,
			TypeId:          0x12,
			MessageStreamId: 1,
		},
	}

	rsar := RtmpSampleAccessResponse{
		Name:  "|RtmpSampleAccess",
		Audio: true,
		video: false,
	}
	payload, _ := amf.Marshal(rsar)
	buffer.Payload = payload
	this.write(buffer.Bytes())
}
