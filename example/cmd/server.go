package main

import (
	"context"
	"github.com/autsu/redtower/example"
	"github.com/autsu/redtower/server"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Ltime)
}

func main() {
	// pprof
	go func() {
		http.ListenAndServe("0.0.0.0:6060", nil)
	}()

	s := server.NewTCPServer("localhost", "8080", "server1")
	s.AddHandler(example.EchoMsg, &example.EchoHandler{})
	s.AddHandler(example.HTTPMsg, &example.HttpEchoPostFormHandler{})
	server.Debug.PrintTypMap(s.Router())

	//go func() {
	//	// 开启监控
	//	monitor := server.NewMonitor(s)
	//	monitor.Start(os.Stdout, time.Second * 10)
	//}()

	s.Server(context.Background())
}
