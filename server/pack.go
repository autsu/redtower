package server

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/youseebiggirl/redtower/conf"
	"io"
	"log"
	"unsafe"
)

// 封包格式：
// =======  head ========== ==== body ==
// +-----------+--------------+-----------+
// |  datalen  | dataType  	  |   data	  |
// |-----------|--------------|-----------|
//   uint32		  interface{}     []byte

type Pack interface {
	HeadSize() uint32
	Packet(msg *Message) ([]byte, error)
	UnPack([]byte) (*Message, error)
}

type DataPack struct {
	m Message
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (d *DataPack) HeadSize() uint32 {
	// head = DataLen() + DataType()
	// DataLen() uint32 = 4 byte
	// DataType() interface{} = 16 byte	 ps: go 里面的接口类型大小为 16 bit
	return (uint32)(unsafe.Sizeof(d.m.dataLen) + unsafe.Sizeof(d.m.typ))
}

func (d *DataPack) Packet(msg *Message) ([]byte, error) {
	data := msg.Data()
	dataLen := msg.DataLen()
	dataType := msg.Type()

	var b bytes.Buffer

	// 按照封包格式依次写入，注意顺序不能出错，否则拆包时会出现不可预计的错误
	if err := binary.Write(&b, binary.BigEndian, &dataLen); err != nil {
		log.Println("packet write dataLen error: ", err)
		return nil, err
	}

	if err := binary.Write(&b, binary.BigEndian, dataType); err != nil {
		log.Println("packet write dataType error: ", err)
		return nil, err
	}

	if err := binary.Write(&b, binary.BigEndian, &data); err != nil {
		log.Println("packet write data error: ", err)
		return nil, err
	}

	return b.Bytes(), nil
}

func (d *DataPack) UnPack(pkg []byte) (*Message, error) {
	r := bytes.NewReader(pkg)

	var m Message

	// 拆包同样也按照格式顺序
	if err := binary.Read(r, binary.BigEndian, &m.dataLen); err != nil {
		log.Println("packet read dataLen error: ", err)
		return nil, err
	}

	// 需要先定义一个实现了 MessageType 的 struct
	// 如果 binary.Read 第三个参数直接传 m.typ 会产生 nil pointer 异常
	// 传 &m.typ 也会发生错误：invalid type
	var x __xxx__
	if err := binary.Read(r, binary.BigEndian, &x); err != nil {
		log.Println("packet read dataLen error: ", err)
		return nil, err
	}
	m.typ = &x

	if m.dataLen > conf.MaxPackSize {
		log.Printf("unpack error: too large message, len: %v, max: %v \n",
			m.dataLen, conf.MaxPackSize)
		return nil, errors.New("unpack error: too large message")
	}

	return &m, nil
}

func (d *DataPack) UnPackFromConn(conn Conn) (*Message, error) {
	head := make([]byte, d.HeadSize())

	_, err := io.ReadFull(conn.SocketConn(), head)
	if err != nil {
		log.Println("read head error: ", err)
		return nil, err
	}

	body, err := d.UnPack(head)
	if err != nil {
		return nil, err
	}

	return body, nil
}
