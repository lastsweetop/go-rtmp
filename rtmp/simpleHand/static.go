package simpleHand

import (
	"bytes"
	"time"
)

func S0S1S2(C0, C1 []byte) []byte {
	S0 := C0
	S1 := generaS1(C1)
	S2 := generaS2(C1)
	s0s1s2 := bytes.NewBuffer(S0)
	s0s1s2.Write(S1)
	s0s1s2.Write(S2)
	return s0s1s2.Bytes()
}

func generaS1(C1 []byte) []byte {
	S1 := make([]byte, len(C1))
	copy(S1, C1)
	return S1
}

func generaS2(C1 []byte) []byte {
	S2 := C1
	now := int(time.Now().Unix())
	S2[4] = byte(now >> 24)
	S2[5] = byte(now >> 16)
	S2[6] = byte(now >> 8)
	S2[7] = byte(now)
	return S2
}
