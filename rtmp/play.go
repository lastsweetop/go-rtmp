package rtmp

import (
	"go-rtmp/buffertool"
	"go-rtmp/h2642flv"
)

func (this *Session) handPlay(b *buffertool.Buffer, version int) {
	this.StreamBegin()

	this.OnStatus()
	this.RtmpSampleAccess()
	this.OnMetaData()

	h2642flv.NewProcess().Start("/Users/sweetop/Code/mystudy/golang/goh264/test/source/7.h264",
		this.Conn)
}
