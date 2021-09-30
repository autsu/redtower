package heartbeat

//import (
//	"context"
//	"log"
//	"testing"
//	"time"
//	"zinx/client"
//	"zinx/server"
//)
//
//type EchoHandler struct {
//	server.BasicHandler
//}
//
//func NewEchoHandler() *EchoHandler {
//	return &EchoHandler{}
//}
//
//func (e *EchoHandler) Handle(r *server.Request) {
//	data := r.Data()
//
//	msg := server.NewMessage(data, r.MsgType())
//	_, err := r.SocketConn().Send(msg)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//}
//
//func TestServer(t *testing.T) {
//	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Ltime)
//	s := server.NewTCPServer("localhost", "8080", "server1")
//	//s.AddHandler(server.EchoMsg, NewEchoHandler())
//	s.Start(context.Background())
//}
//
//func TestClient(t *testing.T) {
//	cli := client.NewClientWithTCP("localhost", "8080")
//	// 发送心跳包，持续 10 秒
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
//	defer cancel()
//
//	cli.StartHeartbeat(ctx)
//
//	select {}
//}
//
//// 不发心跳包，只发数据
//func TestClient1(t *testing.T) {
//	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Ltime)
//	cli := client.NewClientWithTCP("localhost", "8080")
//
//	msgs := []string{"666", "456", "789", "999", "zxzxz", "12435", "dsfsd"}
//
//	for i := 0; i < len(msgs); i++ {
//		msg := server.NewMessage([]byte(msgs[i]), server.EchoMsg)
//		n, err := cli.Send(msg)
//		log.Printf("send %d bytes", n)
//		if err != nil {
//			log.Println("send error: ", err)
//			continue
//		}
//
//		receive, err := cli.Receive()
//		if err != nil {
//			log.Println("recv error: ", err)
//			continue
//		}
//
//		log.Println(string(receive.Data()))
//		time.Sleep(time.Second * 3)
//
//		if i == len(msgs)-1 {
//			i = 0
//		}
//	}
//}
