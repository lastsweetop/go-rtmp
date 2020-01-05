package main

import "go-rtmp/rtmp"

func main() {
	server := rtmp.Server{Port: "1935"}
	server.Start()
}
