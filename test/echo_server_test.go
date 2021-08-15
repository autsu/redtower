package test

import (
	"context"
	"testing"
	"zinx/server"
)

func TestEchoServer(t *testing.T) {
	ts := server.NewTCPServer("localhost", "8080", "test1")
	ts.Start(context.Background())
}
