package server

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
)

// 封包格式：
// =======  head ========== ==== body ==
// +-----------+-----------+-----------+
// |  datalen  | dataType  |   data	   |
// |-----------|-----------|-----------|
//   uint32		  int32       []byte

const MaxPackSize = 25535

type Pack interface {
	HeadSize() uint32
	Packet(msg *Message) ([]byte, error)
	UnPack([]byte) (*Message, error)
}

type DataPack struct {}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (d *DataPack) HeadSize() uint32 {
	// head = DataLen() + DataType()
	// DataLen() uint32 = 4 byte
	// DataType() int32 = 4 byte
	return 8
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

	if err := binary.Write(&b, binary.BigEndian, &dataType); err != nil {
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

	if err := binary.Read(r, binary.BigEndian, &m.type_); err != nil {
		log.Println("packet read dataLen error: ", err)
		return nil, err
	}

	if m.dataLen > MaxPackSize {
		return nil, errors.New("unpack error: too large message")
	}

	return &m, nil
}



