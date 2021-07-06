package test

import (
	"testing"
	"zinx/server"
)

func TestEchoServer(t *testing.T) {
	ts := server.NewTCPServer("localhost", "8080", "test1")
	ts.Start()
}
