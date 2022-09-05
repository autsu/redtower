package server

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/autsu/redtower/conf"
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
	return (uint32)(unsafe.Sizeof(d.m.dataLen) + unsafe.Sizeof(d.m.typ))
}

func (d *DataPack) Packet(msg *Message) ([]byte, error) {
	data := msg.data
	dataLen := msg.dataLen
	dataType := msg.typ

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

	if err := binary.Read(r, binary.BigEndian, &m.typ); err != nil {
		log.Println("packet read dataLen error: ", err)
		return nil, err
	}

	if m.dataLen > uint64(conf.GlobalConf.MaxPackageSize) {
		log.Printf("unpack error: too large message, len: %v, max: %v \n",
			m.dataLen, conf.GlobalConf.MaxPackageSize)
		return nil, errors.New("unpack error: too large message")
	}

	return &m, nil
}
