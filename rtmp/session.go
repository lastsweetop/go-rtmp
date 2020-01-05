package rtmp

import (
	"go-rtmp/buffertool"
	"go-rtmp/enums/AMFVersion"
	"go-rtmp/model"
	"log"
	"net"
)

type Session struct {
	Conn     net.Conn
	datachan chan []byte
	buffer   []byte

	ChunksMap map[uint32]*model.MessageHeader

	ChunkStreamId uint32

	WindowAcknowledgementSize int
}

func New(conn net.Conn) *Session {
	return &Session{Conn: conn,
		datachan:      make(chan []byte, 10),
		buffer:        make([]byte, 0),
		ChunksMap:     make(map[uint32]*model.MessageHeader),
		ChunkStreamId: 3,
	}
}

func (this *Session) Proc() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	this.handShake()

	go this.pkgChunk()
	this.readdata()
}

func (this *Session) readdata() {
	for {
		buffer := make([]byte, 1000)
		n, err := this.Conn.Read(buffer)
		if err != nil {
			log.Println(this.Conn.RemoteAddr(), "error read", err.Error())
			this.Conn.Close()
			return
		}
		log.Println(buffer[:n])
		this.datachan <- buffer[:n]
	}
}

func (this *Session) write(data []byte) {
	this.Conn.Write(data)
}

func (this *Session) pkgChunk() {
	for {
		basicHeader := this.scanBasicHeader()
		log.Println("basicHeader", basicHeader)
		messageHeader := this.scanMessageHeader(basicHeader)
		log.Println("messageHeader", messageHeader)
		this.ChunksMap[basicHeader.ChunkStreamId] = messageHeader
		temps := make([]byte, 0)
		partNum := int((messageHeader.PayloadLength + 1) / 128)
		log.Println("partNum", partNum)
		if partNum > 0 {
			for i := 0; i < partNum; i++ {
				temps = buffertool.Merge(temps, this.getData(128))
				this.getData(1)
			}
			temps = buffertool.Merge(temps, this.getData(int(messageHeader.PayloadLength)-128*partNum))
		} else {
			temps = buffertool.Merge(temps, this.getData(int(messageHeader.PayloadLength)))
		}

		switch messageHeader.TypeId {
		case 0x14:
			this.handCommand(buffertool.NewBuffer(temps), AMFVersion.AMF0)
			break
		case 0x05:
			this.handWindowAcknowledgementSize(buffertool.NewBuffer(temps), AMFVersion.AMF0)
			break
		case 0x04:
			this.handUserControl(buffertool.NewBuffer(temps), AMFVersion.AMF0)
			break
		}
	}
}

func (this *Session) getData(n int) []byte {
	for len(this.buffer) < n {
		data := <-this.datachan
		this.buffer = buffertool.Merge(this.buffer, data)
	}
	temp := this.buffer[:n]
	this.buffer = this.buffer[n:]
	return temp
}

func (this *Session) scanBasicHeader() *model.BasicHeader {
	b := buffertool.NewBuffer(this.getData(1))
	basicHeader := &model.BasicHeader{}
	basicHeader.Fmt = b.ScanByteForNBit(2)
	basicHeader.ChunkStreamId = uint32(b.ScanByteForNBit(6))
	switch basicHeader.ChunkStreamId {
	case 0:
		b = buffertool.NewBuffer(this.getData(1))
		basicHeader.ChunkStreamId = uint32(b.ScanByte()) + 64
		break
	case 1:
		b = buffertool.NewBuffer(this.getData(2))
		basicHeader.ChunkStreamId = uint32(b.ScanByte())<<8 + uint32(b.ScanByte()) + 64
		break
	case 2:
		break
	default:
		break
	}
	return basicHeader
}

func (this *Session) scanMessageHeader(bh *model.BasicHeader) *model.MessageHeader {
	messageHeader := &model.MessageHeader{}
	mh := this.ChunksMap[bh.ChunkStreamId]
	switch bh.Fmt {
	case 0:
		b := buffertool.NewBuffer(this.getData(11))
		messageHeader.Timestamp = uint32(b.ScanByte())<<16 + uint32(b.ScanByte())<<8 + uint32(b.ScanByte())
		messageHeader.PayloadLength = uint32(b.ScanByte())<<16 + uint32(b.ScanByte())<<8 + uint32(b.ScanByte())
		messageHeader.TypeId = b.ScanByte()
		messageHeader.MessageStreamId = uint32(b.ScanByte())<<24 + uint32(b.ScanByte())<<16 + uint32(b.ScanByte())<<8 + uint32(b.ScanByte())
		messageHeader.Headlen = 12
		break
	case 1:
		b := buffertool.NewBuffer(this.getData(7))
		messageHeader.Timestamp = uint32(b.ScanByte())<<16 + uint32(b.ScanByte())<<8 + uint32(b.ScanByte())
		messageHeader.PayloadLength = uint32(b.ScanByte())<<16 + uint32(b.ScanByte())<<8 + uint32(b.ScanByte())
		messageHeader.TypeId = b.ScanByte()
		messageHeader.MessageStreamId = mh.MessageStreamId
		messageHeader.Headlen = 8
		break
	case 2:
		b := buffertool.NewBuffer(this.getData(7))
		messageHeader.Timestamp = uint32(b.ScanByte())<<16 + uint32(b.ScanByte())<<8 + uint32(b.ScanByte())
		messageHeader.PayloadLength = mh.PayloadLength
		messageHeader.TypeId = mh.TypeId
		messageHeader.MessageStreamId = mh.MessageStreamId
		messageHeader.Headlen = 4
		break
	case 3:
		messageHeader.Timestamp = mh.Timestamp
		messageHeader.PayloadLength = mh.PayloadLength
		messageHeader.TypeId = mh.TypeId
		messageHeader.MessageStreamId = mh.MessageStreamId
		messageHeader.Headlen = 1
		break
	}
	return messageHeader
}
