package rtmp

import (
	"go-rtmp/amf"
	"go-rtmp/model"
)

type OnMetaDataResponse struct {
	Name     string
	MetaData MetaData
}

type MetaData struct {
	Server        string  `amf:"Server"`
	Width         float64 `amf:"width"`
	Height        float64 `amf:"height"`
	DisplayWidth  float64 `amf:"displayWidth"`
	DisplayHeight float64 `amf:"displayHeight"`
	Duration      float64 `amf:"duration"`
	Framerate     float64 `amf:"framerate"`
	Fps           float64 `amf:"fps"`
	Videodatarate float64 `amf:"videodatarate"`
	Videocodecid  float64 `amf:"videocodecid"`
	Audiodatarate float64 `amf:"audiodatarate"`
	Audiocodecid  float64 `amf:"audiocodecid"`
	Profile       string  `amf:"profile"`
	level         string  `amf:"level"`
}

func (this *Session) OnMetaData() {
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

	omdr := OnMetaDataResponse{
		Name: "onMetaData",
		MetaData: MetaData{
			Server:        "https://github.com/lastsweetop/go-rtmp",
			Width:         1920,
			Height:        1080,
			DisplayWidth:  1920,
			DisplayHeight: 1080,
			Duration:      0,
			Framerate:     25,
			Fps:           25,
			Videodatarate: 1935,
			Videocodecid:  7,
			Audiodatarate: 136,
			Audiocodecid:  10,
			Profile: string([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}),
			level: string([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}),
		},
	}

	payload, _ := amf.Marshal(omdr)
	buffer.Payload = payload
	this.write(buffer.Bytes())
}
