package tools

import "math/rand"

func RandomBytesArray(b []byte) {
	for i, _ := range b {
		b[i] = byte(rand.Intn(255))
	}
}
