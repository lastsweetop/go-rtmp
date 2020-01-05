package test

import (
	"go-rtmp/amf"
	"go-rtmp/rtmp"
	"log"
	"testing"
)

func TestAMF(t *testing.T) {
	scc := rtmp.ConnectResult{
		Name:          "connect",
		TransactionID: 1,
		Properties: rtmp.ConnectResultProperties{
			FmsVer:       "FMS/3,0,1,123",
			Capabilities: 31,
		}, Information: rtmp.ConnectResultInformation{
			Level:          "status",
			Code:           "NetConnection.Connect.Success",
			Description:    "Connection succeeded",
			ObjectEncoding: 0,
		},
	}
	log.Println(amf.Marshal(scc))
}
