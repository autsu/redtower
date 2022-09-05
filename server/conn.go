package server

import (
	"errors"
	"github.com/autsu/redtower/logs"
	"io"
	"log"
	"net"
	"sync/atomic"
)

type Conn interface {
	Addr() net.Addr
	ConnID() uint64

	Send(data *Message) (int, error) // 发送数据到 conn
	Receive() (*Message, error)      // 从 conn 中接收数据
	Handle()                         // 处理 conn 中的数据
	Close()                          // 关闭连接
}

// 每创建一条 conn，便从这里取得一个新的 connId，之后 connId++（原子的）
var connId atomic.Uint64

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
		id:            connId.Load(),
		isClose:       false,
		heartbeatChan: make(chan struct{}),
	}
	// global conn id + 1
	connId.Add(1)
	return t
}

func (t *TCPConn) Close() {
	if t.isClose {
		return
	}
	t.isClose = true
	GlobalConnManage.RemoveAndClose(t.server.Name, t)
}

func (t *TCPConn) Addr() net.Addr {
	return t.socketConn.RemoteAddr()
}

func (t *TCPConn) ConnID() uint64 {
	return t.id
}

func (t *TCPConn) Send(data *Message) (int, error) {
	logs.L.Printf_IfOpenDebug("seng to %v ...\n", t.Addr())
	if t.isClose {
		return 0, errors.New("read error: the connection is close")
	}

	pack := NewDataPack()
	packmsg, err := pack.Packet(data)
	if err != nil {
		return 0, err
	}

	n, err := t.socketConn.Write(packmsg)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (t *TCPConn) Receive() (*Message, error) {
	logs.L.Printf_IfOpenDebug("receive from %v ...\n", t.Addr())
	if t.isClose {
		return nil, errors.New("read error: the connection is close")
	}

	pack := NewDataPack()
	header, err := t.UnpackHeader(pack)
	if err != nil {
		return nil, err
	}

	dataLen := header.dataLen
	msgType := header.typ

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
		if err != nil {
			// 虽然没读满，但是返回了 EOF 错误，说明没有数据可读了，可能是对方已经断开了连接，
			// 此时就不需要再尝试读取了
			if err == io.EOF {
				return nil, err
			}
			// 将当前读的这部分添加到 tmp 中，暂时保存
			tmpBuf = append(tmpBuf, buf[:n]...)
			needN -= uint64(n) // 更新 needN 的值
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

func (t *TCPConn) UnpackHeader(d *DataPack) (*Message, error) {
	head := make([]byte, d.HeadSize())

	_, err := io.ReadFull(t.socketConn, head)
	if err != nil {
		if err != io.EOF {
			log.Println("read head error: ", err)
		}
		return nil, err
	}

	body, err := d.UnPack(head)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Handle 处理连接中的数据，
func (t *TCPConn) Handle() {
	defer func() {
		if !t.isClose {
			t.isClose = true
			// 从该连接从当前 server 的 connManage 中移除，并 close 该连接
			GlobalConnManage.RemoveAndClose(t.server.Name, t)
			logs.L.Printf_IfOpenDebug("conn [id = %v, addr = %v] close\n", t.id, t.Addr())
		}
	}()
	// 开启心跳检测
	heartBeat := NewHeartBeat(t.server.Name, t)
	go heartBeat.Start()
	GlobalConnManage.Add(t.server.Name, t) // 将连接添加到 connManage

	for {
		if err := handleTCP(t, heartBeat); err != nil {
			bmsg := []byte(err.Error())
			t.Send(NewErrorMessage(bmsg))
			break // EOF 会 break
		}

		if t.isClose {
			break
		}
	}
}

func handleTCP(conn *TCPConn, beat *HeartBeat) error {
	recvData, err := conn.Receive()
	if err != nil {
		log.Println(err)
		return err
	}
	// 读取到数据了，发送心跳
	beat.trigger()

	msgType := recvData.Type()
	//log.Printf("recvData: data: %v, type: %v, size: %v bytes \n",
	//	string(recvData.Data()), msgType, len(recvData.Data()))

	msg := NewMessage(recvData.Data(), msgType)
	req := NewRequest(msg, conn)
	// 添加到工作池
	conn.server.pool.AddWork(req)
	return nil
}
