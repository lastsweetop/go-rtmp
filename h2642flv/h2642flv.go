package h2642flv

import (
	"bufio"
	"bytes"
	"go-rtmp/enums/NalUnitTypes"
	"go-rtmp/model"
	"log"
	"net"
	"os"
	"sync"
)

type Process struct {
	bufChan         chan []byte
	frameChan       chan []byte
	closeSplitChan  bool
	closeDecodeChan bool
	wg              sync.WaitGroup
}

func NewProcess() *Process {
	return &Process{
		bufChan:         make(chan []byte, 10),
		frameChan:       make(chan []byte, 10),
		closeSplitChan:  false,
		closeDecodeChan: false,
		wg:              sync.WaitGroup{},
	}
}

func (this *Process) Start(inPath string, conn net.Conn) {
	fi, err := os.Open(inPath)
	if err != nil {
		log.Println("open error : " + err.Error())
	}
	defer fi.Close()
	r := bufio.NewReader(fi)

	go this.splitFrame()
	go this.decodeFrame(conn)

	for {
		buf := make([]byte, 4096)
		n, err := r.Read(buf)
		if err != nil && err.Error() != "EOF" {
			log.Println("read error : " + err.Error())
		}
		if 0 == n {
			break
		}
		this.bufChan <- buf[:n]
	}
	this.closeSplitChan = true
	this.wg.Wait()
}

func (this *Process) splitFrame() {
	temp := make([]byte, 0)
	frame := make([]byte, 0)
	frameNum := 0

	this.wg.Add(1)
	defer func() {
		this.wg.Done()
	}()

	for {
		select {
		case buf := <-this.bufChan:
			buf = append(temp, buf...) //头部追加上个包剩余不满4个字节的数据
			temp = make([]byte, 0)
			i := 0
			for i < len(buf) {
				if i > len(buf)-4 {
					temp = buf[i:] //剩余不足4个字节的数据
					break
				}

				if buf[i] == 0x00 && buf[i+1] == 0x00 && buf[i+2] == 0x00 && buf[i+3] == 0x01 {
					if i != 0 {
						this.frameChan <- frame
						frame = make([]byte, 0)
					}
					frameNum++ //frame数量+1
					i += 4
				} else if buf[i] == 0x00 && buf[i+1] == 0x00 && buf[i+2] == 0x01 {
					if i != 0 {
						this.frameChan <- frame
						frame = make([]byte, 0)
					}
					frameNum++ //frame数量+1
					i += 3
				} else {
					frame = append(frame, buf[i])
					i++
				}
			}
			break
		default:
			if this.closeSplitChan {
				frame = append(frame, temp...)
				this.frameChan <- frame
				this.closeDecodeChan = true
				log.Println("关闭切割split协程")
				log.Println("frame数量：", frameNum)
				return
			}

		}
	}
}

func (this *Process) decodeFrame(conn net.Conn) {
	first := true
	defer func() {
		this.wg.Done()
	}()
	this.wg.Add(1)

	//conn.Write([]byte{0x46, 0x4c, 0x56, 0x01, 0x01, 0x00, 0x00, 0x00, 0x09})
	//w.Write(getHeader())
	i := 0
	tempsps := []byte{}
	for {
		select {
		case frame := <-this.frameChan:

			flen := len(frame)
			i++
			log.Println("帧", i, "的数据长度：", flen)
			dataStream := model.StreamData{Data: frame, Index: 0}
			_ = dataStream.F(1)
			//log.Printf("forbidden_bit %d\n", ForbiddenBit)
			_ = dataStream.U(2)
			//log.Printf("nal_reference_bit %d\n", NalReferenceBit)
			NalUnitType := dataStream.U(5)
			//log.Printf("nal_unit_type %d\n", NalUnitType)

			switch NalUnitType {
			case NalUnitTypes.SPS:
				tempsps = frame
				log.Println("SPS")
				break
			case NalUnitTypes.PPS:
				log.Println("PPS")
				if first {
					flvspspps(conn, tempsps, frame)
					log.Println("flvspspps")
					first = false
				}
				break
			case NalUnitTypes.SEI:
				log.Println("SEI")
				break
			case NalUnitTypes.IDR:
				log.Println("IDR")
				flvnalu(conn, 0x17, frame)
				break
			case NalUnitTypes.NOIDR:
				log.Println("NOIDR")
				flvnalu(conn, 0x27, frame)
				break
			default:
				break
			}
			break
		default:
			if this.closeDecodeChan {
				log.Println("关闭切割decode协程")
				return
			}
		}
	}
}

func flvnalu(conn net.Conn, fc byte, nalu []byte) {
	nlen := len(nalu)
	data := bytes.NewBuffer([]byte{})
	data.Write([]byte{fc,
		0x01,              //AVCPacketType
		0x00, 0x00, 0x00}) //Composition Time

	data.Write([]byte{byte(nlen >> 24), byte(nlen >> 16), byte(nlen >> 8), byte(nlen)})
	data.Write(nalu)
	writeVideo(conn, data.Bytes())
}

func flvspspps(conn net.Conn, sps []byte, pps []byte) {
	data := bytes.NewBuffer([]byte{})
	data.Write([]byte{0x17, 0x00, 0x00, 0x00, 0x00})
	data.Write([]byte{
		0x01,                                //configurationVersion
		sps[1],                              //AVCProfileIndication
		sps[2],                              //profile_compatibiltity
		sps[3],                              //AVCLevelIndication
		0xFF,                                //lengthSizeMinusOne
		0xE1,                                //numOfSequenceParameterSets
		byte(len(sps) >> 8), byte(len(sps)), //sequenceParameterSetLength
	})
	data.Write(sps)
	data.Write([]byte{
		0x01,                                //numOfPictureParamterSets
		byte(len(pps) >> 8), byte(len(pps)), //pictureParameterSetLength
	})
	data.Write(pps)
	writeVideo(conn, data.Bytes())
	//writeAudio(conn)
}

var ChunkStreamId uint32 = 7

func writeAudio(conn net.Conn) {
	conn.Write([]byte{0x07, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x08, 0x01, 0x00, 0x00, 0x00, 0xAF, 0x00, 0x11, 0x90})
}

func writeVideo(conn net.Conn, payload []byte) {
	buffer := model.Chunk{
		BasicHeader: &model.BasicHeader{
			Fmt:           1,
			ChunkStreamId: ChunkStreamId,
		},
		MessageHeader: &model.MessageHeader{
			Timestamp:       40,
			TypeId:          0x09,
			MessageStreamId: 1,
		},
	}
	buffer.Payload = payload
	conn.Write(buffer.Bytes())
}
