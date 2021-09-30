package server

import (
	"errors"
	"io"
	"log"
	"net"
	"sync/atomic"
)

type Conn interface {
	Server() Server
	SocketConn() net.Conn // 原始的 net.conn
	RemoteAddr() net.Addr
	ConnID() uint64
	IsClose() bool
	HeartBeatChan() chan struct{}

	SetIsClose(bool)

	Send(data *Message) (int, error) // 发送数据到 conn
	Receive() (*Message, error)      // 从 conn 中接收数据
	Handler()                        // 处理 conn 中的数据
	Stop()                           // 关闭连接
}

// 每创建一条 conn，便从这里取得一个新的 connId，之后 connId++（原子的）
var connId uint64 = 1

type TCPConn struct {
	server        *TCPServer // 该连接所属的 server
	socketConn    *net.TCPConn
	id            uint64
	isClose       bool
	heartbeatChan chan struct{}
}

func NewTCPConn(conn *net.TCPConn, server *TCPServer) *TCPConn {
	t := &TCPConn{
		server:        server,
		socketConn:    conn,
		id:            atomic.LoadUint64(&connId),
		isClose:       false,
		heartbeatChan: make(chan struct{}),
	}
	// global conn id + 1
	for !atomic.CompareAndSwapUint64(&connId, connId, connId+1) {

	}

	return t
}

func (t *TCPConn) Stop() {
	if t.isClose {
		return
	}
	t.isClose = true
	t.socketConn.Close()
}

func (t *TCPConn) Server() Server {
	return t.server
}

func (t *TCPConn) SocketConn() net.Conn {
	return t.socketConn
}

func (t *TCPConn) RemoteAddr() net.Addr {
	return t.socketConn.RemoteAddr()
}

func (t *TCPConn) ConnID() uint64 {
	return t.id
}

func (t *TCPConn) Send(data *Message) (int, error) {
	log.Printf("seng to %v ...", t.RemoteAddr().String())
	if t.isClose {
		return 0, errors.New("read error: the connection is close")
	}

	pack := NewDataPack()

	packmsg, err := pack.Packet(data)
	if err != nil {
		return 0, err
	}

	n, err := t.SocketConn().Write(packmsg)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (t *TCPConn) Receive() (*Message, error) {
	log.Printf("receive from %v ...", t.RemoteAddr().String())
	if t.isClose {
		return nil, errors.New("read error: the connection is close")
	}

	pack := NewDataPack()
	msg, err := pack.UnPackFromConn(t)
	if err != nil {
		return nil, err
	}

	dataLen := msg.DataLen()
	msgType := msg.Type()

	// 虽然协议中记录了消息的长度，但是 tcp 传输时可能会发生截断，导致读取不完整，比如
	// 记录的长度是 100，但是因为某些原因（流量控制、拥塞控制）导致 tcp 只发送了长度 80
	// 的数据，还有 20 的数据将在下次传输时发送，如果在这种情况下，还是使用
	// buf := make([]byte, dataLen); conn.Read(buf) 这样的方式，会导致没有读取到完整数据
	//
	// 解决方法是：额外创建一个 tmpBuf 用来存储那些没读完的数据，比如在上面的例子中，就可以
	// 先将第一次读出来的 80 放到 tmpBuf 中，使用 io.ReadFull() 可以保证读取到需要的长度,
	// 同时用一个变量记录还剩多少需要读取，通过一个 for 循环，去读取剩下的数据，最终将之前保存
	// 的数据合并，便是一个完整的数据
	var (
		tmpBuf = make([]byte, 0, 4096)
		needN  = dataLen // 还剩多少需要读取
		readN  int64     // 总共需要读取多少
	)

	for {
		buf := make([]byte, needN)
		// 必须读满 len(buf)，否则返回一个 err
		n, err := io.ReadFull(t.socketConn, buf)
		// 没有读满 buf
		if err != nil && err != io.EOF {
			log.Println("read data error: ", err)
			// 将当前读的这部分添加到 tmp 中，暂时保存
			tmpBuf = append(tmpBuf, buf[:n]...)
			needN -= uint32(n) // 更新 needN 的值
			readN += int64(n)
			continue
		}
		// 读满了
		tmpBuf = append(tmpBuf, buf...)
		readN += int64(n)

		if readN == int64(dataLen) {
			break
		}
	}
	m := NewMessage(tmpBuf, msgType)
	return m, nil
}

func (t *TCPConn) Handler() {
	defer func() {
		t.Stop()
		t.SetIsClose(true)
		t.Server().ConnManage().Remove(t) // 从该连接从当前 server 的 connManage 中移除
		log.Printf("conn [id = %v, addr = %v] close\n", t.id, t.RemoteAddr())
	}()

	// 开启心跳检测
	heartBeat := NewHeartBeat(t)
	go heartBeat.Start()

	for {
		if err := handlerConn(t); err != nil {
			log.Println(err)
			bmsg := []byte(err.Error())
			t.Send(NewErrorMessage(bmsg))
			break // EOF 会 break
		}

		if t.IsClose() {
			break
		}
	}
}

func handlerConn(conn Conn) error {
	recvData, err := conn.Receive()
	if err != nil {
		log.Println(err)
		return err
	}

	msgType := recvData.Type()
	log.Printf("recvData: data: %v, type: %v, size: %v bytes \n",
		string(recvData.Data()), msgType, len(recvData.Data()))

	msg := NewMessage(recvData.Data(), msgType)
	req := NewRequest(msg, conn)
	// 添加到工作池
	conn.Server().Pool().AddWork(req)

	return nil
}

func (t *TCPConn) IsClose() bool {
	return t.isClose
}

func (t *TCPConn) SetIsClose(isClose bool) {
	t.isClose = isClose
}

func (t *TCPConn) HeartBeatChan() chan struct{} {
	return t.heartbeatChan
}
