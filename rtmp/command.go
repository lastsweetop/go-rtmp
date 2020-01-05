package rtmp

import (
	"go-rtmp/buffertool"
	"log"
)

func (this *Session) handCommand(b *buffertool.Buffer, AMFVersion int) {
	_ = b.ScanByte()
	name := b.ScanAMFString()
	log.Println("Command Name:", name)
	switch name {
	case "connect":
		this.handConnect(b, AMFVersion)
		break
	case "createStream":
		this.handCreateStream(b, AMFVersion)
		break
	case "play":
		this.handPlay(b, AMFVersion)
		break
	}
}
