package server

import (
	"context"
	reuseport "github.com/kavu/go_reuseport"
	"log"
	"net"
	"zinx/conf"
)

type Server interface {
	Start(context.Context)
	Stop(context.Context)
	Server(context.Context)
	ConnManage() *ConnManage
	AddHandler(MessageType, Handler)
	Pool() *Pool
}

type TCPServer struct {
	port   string
	host   string
	Name   string
	router *Router
	manage *ConnManage
	pool   *Pool
}

func NewTCPServer(host, port, name string) *TCPServer {
	t := &TCPServer{
		port:   port,
		host:   host,
		Name:   name,
		router: NewRouter(),
		manage: NewConnManage(),
	}
	// 默认在路由中添加心跳处理函数
	t.AddHandler(HeartBeatMsg, NewHeartBeatHandler())

	return t
}

func (t *TCPServer) Start(ctx context.Context) {
	addr := t.host + ":" + t.port
	log.Printf("server[name: %v] start in %v... \n", t.Name, addr)

	// 开启 SO_REUSEPORT 端口复用
	l, err := reuseport.Listen("tcp", addr)
	listen, ok := l.(*net.TCPListener)
	if !ok {
		log.Println("not a tcp listener: ", err)
		return
	}

	if err != nil {
		log.Println("listen socket error: ", err)
		return
	}

	var (
		id uint64 = 1
	)

	// 开启任务池
	t.pool = NewPool(conf.DefaultGoroutineMaxNum, t.router)
	t.pool.StartWorkerPool(ctx)

	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			log.Println("accept conn error: ", err)
			continue
		}

		if t.manage.Len() > conf.MaxConnNum {
			conn.Write([]byte("Your connection request was rejected"))
			log.Printf("new conn[%v] was rejected\n", conn.RemoteAddr())
		}

		tcpConn := NewTCPConn(conn, t)
		log.Printf("a new conn, remote addr: [%v]\n", tcpConn.Conn().RemoteAddr())

		id++
		t.manage.Add(tcpConn) // 将连接添加到 connManage

		// 开启心跳检测
		heartBeat := NewHeartBeat(tcpConn)
		go heartBeat.Start()

		go func() {
			defer func() {
				tcpConn.Stop()
				tcpConn.SetIsClose(true)
				// 从 connManage 中移除
				t.manage.Remove(tcpConn)
				log.Printf("conn [id = %v, addr = %v] close\n",
					tcpConn.id, tcpConn.RemoteAddr())
			}()

			for {
				if err := tcpConn.Handler(); err != nil {
					log.Println(err)
					break // EOF 会 break
				}
				if tcpConn.IsClose() {
					break
				}
			}
		}()
	}

}

func (t *TCPServer) Stop(ctx context.Context) {

}

func (t *TCPServer) Server(ctx context.Context) {
	t.Start(ctx)
}

func (t *TCPServer) ConnManage() *ConnManage {
	return t.manage
}
func (t *TCPServer) Pool() *Pool {
	return t.pool
}

func (t *TCPServer) AddHandler(typ MessageType, handler Handler) {
	t.router.AddRouter(typ, handler)
}
