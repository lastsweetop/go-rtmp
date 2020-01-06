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

	h2642flv.NewProcess().Start("resource/test2.h264",
		this.Conn)
}
