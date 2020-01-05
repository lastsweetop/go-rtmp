package rtmp

import (
	"log"
	"net"
)

type Server struct {
	Port string
}

func (w *Server) Start() {
	netListen, err := net.Listen("tcp", "0.0.0.0:"+w.Port)
	if err != nil {
		log.Println(err.Error())
	}
	for {
		conn, _ := netListen.Accept()
		rConn := New(conn)
		go rConn.Proc()
	}
}
