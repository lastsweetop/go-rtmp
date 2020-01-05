package rtmp

import (
	"go-rtmp/rtmp/complexHand"
	"go-rtmp/rtmp/simpleHand"
	"log"
)

func (this *Session) handShake() {
	buffer := make([]byte, 1537)
	n, err := this.Conn.Read(buffer)
	if err != nil {
		log.Println(err.Error())
		this.Conn.Close()
	}
	C0 := buffer[0:1]
	C1 := buffer[1:n]
	log.Println("receive c0+c1")
	var s0s1s2 []byte
	if C1[4] == 0 && C1[5] == 0 && C1[6] == 0 && C1[7] == 0 {
		log.Println("simpleHand")
		s0s1s2 = simpleHand.S0S1S2(C0, C1)
	} else {
		log.Println("complexHand")
		s0s1s2 = complexHand.S0S1S2(C0, C1)
	}

	this.Conn.Write(s0s1s2)
	log.Println("send s0+s1+s2")

	C2 := make([]byte, 1536)
	n, err = this.Conn.Read(C2)
	if err != nil {
		log.Println(err.Error())
		this.Conn.Close()
	}
	log.Println("receive c2")
}
