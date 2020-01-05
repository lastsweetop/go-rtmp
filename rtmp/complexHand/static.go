package complexHand

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"go-rtmp/tools"
	"log"
	"math/rand"
	"time"
)

var FP_KEY = []byte{
	0x47, 0x65, 0x6E, 0x75, 0x69, 0x6E, 0x65, 0x20,
	0x41, 0x64, 0x6F, 0x62, 0x65, 0x20, 0x46, 0x6C,
	0x61, 0x73, 0x68, 0x20, 0x50, 0x6C, 0x61, 0x79,
	0x65, 0x72, 0x20, 0x30, 0x30, 0x31, // Genuine Adobe Flash Player 001
	0xF0, 0xEE, 0xC2, 0x4A, 0x80, 0x68, 0xBE, 0xE8,
	0x2E, 0x00, 0xD0, 0xD1, 0x02, 0x9E, 0x7E, 0x57,
	0x6E, 0xEC, 0x5D, 0x2D, 0x29, 0x80, 0x6F, 0xAB,
	0x93, 0xB8, 0xE6, 0x36, 0xCF, 0xEB, 0x31, 0xAE}

var FMSKey = []byte{
	0x47, 0x65, 0x6e, 0x75, 0x69, 0x6e, 0x65, 0x20,
	0x41, 0x64, 0x6f, 0x62, 0x65, 0x20, 0x46, 0x6c,
	0x61, 0x73, 0x68, 0x20, 0x4d, 0x65, 0x64, 0x69,
	0x61, 0x20, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x20, 0x30, 0x30, 0x31, // Genuine Adobe Flash Media Server 001
	0xf0, 0xee, 0xc2, 0x4a, 0x80, 0x68, 0xbe, 0xe8,
	0x2e, 0x00, 0xd0, 0xd1, 0x02, 0x9e, 0x7e, 0x57,
	0x6e, 0xec, 0x5d, 0x2d, 0x29, 0x80, 0x6f, 0xab,
	0x93, 0xb8, 0xe6, 0x36, 0xcf, 0xeb, 0x31, 0xae}

func S0S1S2(C0, C1 []byte) []byte {
	S0 := C0
	S1 := generaS1(C1)
	S2 := generaS2(C1)
	s0s1s2 := bytes.NewBuffer(S0)
	s0s1s2.Write(S1)
	s0s1s2.Write(S2)
	return s0s1s2.Bytes()
}

func generaS2(C1 []byte) []byte {
	S2 := make([]byte, 1536)
	tools.RandomBytesArray(S2)
	c1Digest := calcC1Digest(C1, 0)
	h := hmac.New(sha256.New, FMSKey[:68])
	h.Write(c1Digest)
	tmp_key := h.Sum(nil)

	h2 := hmac.New(sha256.New, tmp_key[:32])
	h2.Write(S2[:1504])
	s2digest := h2.Sum(nil)
	for i := 0; i < 32; i++ {
		S2[1504+i] = s2digest[i]
	}
	return S2
}

func generaS1(C1 []byte) []byte {
	S1 := make([]byte, 1536)
	tools.RandomBytesArray(S1)
	now := int(time.Now().Unix())
	S1[0] = byte(now >> 24)
	S1[1] = byte(now >> 16)
	S1[2] = byte(now >> 8)
	S1[3] = byte(now)
	S1[4] = 0x04
	S1[5] = 0x05
	S1[6] = 0x00
	S1[7] = 0x01

	schemal := 0
	if !verifyC1Digest(C1, schemal) {
		schemal = 1
	}
	if !verifyC1Digest(C1, schemal) {
		schemal = 2
	}

	log.Println("schemal", schemal)
	//key_offset := rand.Intn(628)
	digest_offset := rand.Intn(728)
	base := 0
	if schemal == 1 {
		base = 764
	}
	temp := byte(digest_offset / 4)
	S1[base+8] = temp
	S1[base+9] = temp
	S1[base+10] = temp
	S1[base+11] = temp + byte(digest_offset%4)

	offset := digest_offset + base + 12
	p1 := S1[:offset]
	p2 := S1[offset+32:]
	h := hmac.New(sha256.New, FMSKey[:36])
	h.Write(p1)
	h.Write(p2)
	digest := h.Sum(nil)
	for i := 0; i < 32; i++ {
		S1[offset+i] = digest[i]
	}
	return S1
}

func calcC1Digest(C1 []byte, schemal int) []byte {
	base := 0
	if schemal == 1 {
		base = 764
	}
	offset := (int(C1[base+8]) + int(C1[base+9]) + int(C1[base+10]) + int(C1[base+11])) + 12 + base
	p1 := C1[:offset]
	p2 := C1[offset+32:]

	h := hmac.New(sha256.New, FP_KEY[:30])
	h.Write(p1)
	h.Write(p2)
	digest := h.Sum(nil)
	return digest
}

func verifyC1Digest(C1 []byte, schemal int) bool {
	base := 0
	if schemal == 1 {
		base = 764
	}
	offset := (int(C1[base+8]) + int(C1[base+9]) + int(C1[base+10]) + int(C1[base+11])) + 12 + base
	p1 := C1[:offset]
	p2 := C1[offset+32:]

	h := hmac.New(sha256.New, FP_KEY[:30])
	h.Write(p1)
	h.Write(p2)
	digest := h.Sum(nil)
	for i := 0; i < 32; i++ {
		if C1[offset+i] != digest[i] {
			return false
		}
	}
	return true
}
