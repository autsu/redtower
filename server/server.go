package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/autsu/redtower/conf"
	"github.com/autsu/redtower/logs"
	reuseport "github.com/kavu/go_reuseport"
	"net"
	"os"
	"os/signal"
)

type Server interface {
	Shutdown(context.Context)
	Run(context.Context)
}

type TCPServer struct {
	addr          string
	Name          string
	pool          *Pool
	onRequestFunc func(*Request) error
	cf            *conf.Conf
}

func NewTCPServer(host, port, name string, onRequest func(*Request) error) *TCPServer {
	t := &TCPServer{
		addr:          fmt.Sprintf("%s:%s", host, port),
		Name:          name,
		onRequestFunc: onRequest,
	}
	return t
}

func (t *TCPServer) run(ctx context.Context) error {
	logs.L.Printf_IfOpenDebug("server[name: %v] start in %v... \n", t.Name, t.addr)

	// 开启 SO_REUSEPORT 端口复用
	l, err := reuseport.Listen("tcp", t.addr)
	listen, ok := l.(*net.TCPListener)
	if !ok {
		logs.L.Printf("not a tcp listener: %v\n", err)
		return errors.New(fmt.Sprintf("not a tcp listener: %v", err))
	}

	if err != nil {
		logs.L.Printf("listen socket error: %v\n", err)
		return err
	}

	// 进行初始化工作
	t.init(ctx)

	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			logs.L.Printf("accept conn error: %v\n", err)
			continue
		}

		tcpConn := NewTCPConn(conn, t)
		logs.L.Printf_IfOpenDebug("a new conn, remote addr: [%v]\n", tcpConn.socketConn.RemoteAddr())

		if GlobalConnManage.NumsFromServer(t.Name) > uint64(conf.GlobalConf.MaxConnNum) {
			tcpConn.Send(NewMessage([]byte("Your connection request was rejected"), ErrorMsg))
			logs.L.Printf_IfOpenDebug("new conn[%v] was rejected\n", conn.RemoteAddr())
		}
		go tcpConn.Handle()
	}
}

func (t *TCPServer) init(ctx context.Context) {
	// 开启任务池
	t.pool = newPool(uint64(conf.GlobalConf.GoroutineMaxNum), t.onRequestFunc)
	t.pool.StartWorkerPool(ctx)
	go t.signExit() // 开启信号监听
}

func (t *TCPServer) Shutdown() {
	GlobalConnManage.ClearFromServer(t.Name)
}

func (t *TCPServer) Run(ctx context.Context) error {
	return t.run(ctx)
}

func (t *TCPServer) signExit() {
	signch := make(chan os.Signal, 1)
	signal.Notify(signch, os.Interrupt, os.Kill)
	<-signch
	// 如果收到了信号，则断开所有连接
	t.Shutdown()
	os.Exit(0)
}
