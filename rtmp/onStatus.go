package rtmp

import (
	"go-rtmp/amf"
	"go-rtmp/model"
)

type OnStatusResponse struct {
	Name          string
	TransactionID float64
	CommandObject interface{}
	InfoObject    InfoObject
}

type InfoObject struct {
	Level       string `amf:"level"`
	Code        string `amf:"code"`
	Description string `amf:"description"`
}

func (this *Session) OnStatus() {
	buffer := model.Chunk{
		BasicHeader: &model.BasicHeader{
			Fmt:           0,
			ChunkStreamId: 5,
		},
		MessageHeader: &model.MessageHeader{
			Timestamp:       0,
			TypeId:          0x14,
			MessageStreamId: 1,
		},
	}

	osr := OnStatusResponse{
		Name:          "onStatus",
		TransactionID: 0,
		CommandObject: nil,
		InfoObject: InfoObject{
			Level:       "status",
			Code:        "NetStream.Play.Start",
			Description: "Start live",
		},
	}
	payload, _ := amf.Marshal(osr)
	buffer.Payload = payload
	this.write(buffer.Bytes())
}
