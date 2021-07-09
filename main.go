package main

import "zinx/server"

func main() {
	s := server.NewTCPServer("localhost", "8080", "server1")
	s.AddHandler(server.EchoMsg, server.NewEchoHandler())
	s.AddHandler(server.HeartBeat, server.NewHeartBeatHandler())

	s.Start()
}
