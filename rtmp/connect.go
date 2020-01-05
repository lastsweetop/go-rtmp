package rtmp

import (
	"go-rtmp/amf"
	"go-rtmp/buffertool"
	"go-rtmp/model"
	"log"
)

type ConnectCommand struct {
	Name          string
	TransactionID int
}

type ConnectResult struct {
	Name          string
	TransactionID float64
	Properties    ConnectResultProperties
	Information   ConnectResultInformation
}

type ConnectResultProperties struct {
	FmsVer       string  `amf:"fmsVer"`
	Capabilities float64 `amf:"capabilities"`
}

type ConnectResultInformation struct {
	Level          string  `amf:"level"`
	Code           string  `amf:"code"`
	Description    string  `amf:"description"`
	ObjectEncoding float64 `amf:"objectEncoding"`
}

func (this *Session) handConnect(b *buffertool.Buffer, AMFVersion int) {
	_ = b.ScanByte()
	tid := b.ScanAMFNumber()
	log.Println("transaction Id", tid)

	this.WriteWindowAcknowledgementSize()
	this.WriteSetPeerBandwidth()
	this.WriteSetChunkSize()

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

	sc := ConnectResult{
		Name:          "_result",
		TransactionID: tid,
		Properties: ConnectResultProperties{
			FmsVer:       "FMS/3,0,1,123",
			Capabilities: 31,
		}, Information: ConnectResultInformation{
			Level:          "status",
			Code:           "NetConnection.Connect.Success",
			Description:    "Connection succeeded",
			ObjectEncoding: 0,
		},
	}
	payload, _ := amf.Marshal(sc)
	buffer.Payload = payload
	this.write(buffer.Bytes())
}
