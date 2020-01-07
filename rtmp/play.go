package rtmp

import (
	"go-rtmp/amf"
	"go-rtmp/buffertool"
	"go-rtmp/h2642flv"
	"log"
)

type PlayReq struct {
	TransactionID float64
	CommandObject interface{}
	StreamName    string
	Start         float64
}

func (this *Session) handPlay(b *buffertool.Buffer, version int) {
	req := &PlayReq{}
	amf.UnMarshal(b.RestBytes(), req)
	log.Println("handPlay", req)

	this.StreamBegin()
	this.OnStatus()
	this.RtmpSampleAccess()
	this.OnMetaData()

	h2642flv.NewProcess().Start("resource/"+req.StreamName,
		this.Conn)
}
