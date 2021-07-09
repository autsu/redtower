package server

import (
	"log"
	"net"
	"zinx/conf"
)

type Server interface {
	Start()
	Stop()
	Server()
}

type TCPServer struct {
	Port       string
	Host       string
	Name       string
	Router     *Router
	ConnManage *ConnManage
	Pool       *Pool
}

func NewTCPServer(host, port, name string) *TCPServer {
	return &TCPServer{
		Port:       port,
		Host:       host,
		Name:       name,
		Router:     NewRouter(),
		ConnManage: NewConnManage(),
	}
}

func (t *TCPServer) Start() {
	log.Println("server start...")

	addr := t.Host + ":" + t.Port

	tcpaddr, _ := net.ResolveTCPAddr("tcp", addr)

	listen, err := net.ListenTCP("tcp", tcpaddr)
	if err != nil {
		log.Println("listen socket error: ", err)
		return
	}

	var id uint64 = 1

	go func() {
		// 开启任务池
		t.Pool = NewPool(conf.DefaultGoroutineMaxNum, t.Router)
		if err := t.Pool.StartWorkerPool(); err != nil {
			log.Println("start worker pool error: ", err)
			return
		}
	}()

	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			log.Println("accept conn error: ", err)
			continue
		}

		if t.ConnManage.Len() > conf.MaxConnNum {
			conn.Write([]byte("Your connection request was rejected"))
			log.Printf("new conn[%v] was rejected\n", conn.RemoteAddr())
		}

		tcpConn := NewTCPConn(conn, t, id)
		log.Printf("a new conn, remote addr: [%v]\n", tcpConn.Conn().RemoteAddr())

		id++
		t.ConnManage.Add(tcpConn) // 将连接添加到 connManage

		// 开启心跳检测
		go tcpConn.Heartbeat()

		go func() {
			defer tcpConn.Stop()
			defer func() {
				tcpConn.IsClose = true
				// 从 connManage 中移除
				t.ConnManage.Remove(tcpConn)
			}()
			defer log.Printf("conn [%v] close\n", tcpConn.RemoteAddr())


			for {
				if err := tcpConn.Handler(); err != nil {
					log.Println(err)
					break // EOF 会 break
				}
				if tcpConn.IsClose {
					break
				}
			}
		}()
	}

}

func (t *TCPServer) Stop() {

}

func (t *TCPServer) Server() {

}

func (t *TCPServer) AddHandler(typ MessageType, handler Handler) {
	t.Router.AddRouter(typ, handler)
}
