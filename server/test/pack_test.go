package test

import (
	"fmt"
	"io"
	"log"
	"net"
	_ "net/http/pprof"
	"testing"
	"zinx/server"
)

func init() {
	log.SetFlags(log.Llongfile | log.Ltime)
}

func TestPack(t *testing.T) {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln(err)
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		pack := server.NewDataPack()

		buf := make([]byte, pack.HeadSize())

		_, err = io.ReadFull(conn, buf)
		if err != nil {
			log.Println("read head error: ", err)
			continue
		}

		//fmt.Println("read full...")
		//time.Sleep(time.Second * 10)

		unPack, err := pack.UnPack(buf)
		if err != nil {
			log.Println("unpack error: ", err)
			continue
		}

		//fmt.Println("unpack...")
		//time.Sleep(time.Second * 10)

		b := make([]byte, unPack.DataLen())
		log.Println("data len: ", unPack.DataLen())
		if len(b) > 1024 {
			log.Println("data too large")
			break
		}
		_, err = conn.Read(b)
		if err != nil {
			log.Println("read data error: ", err)
			continue
		}

		fmt.Println((string)(b))


	}
}

func TestPackClient(t *testing.T) {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Fatalln(err)
	}

	var data []byte
	data = ([]byte)("123")

	msg := server.NewMessage(data, server.TextMsg)
	log.Println(msg.DataLen())

	pack := server.NewDataPack()

	packet, err := pack.Packet(msg)
	if err != nil {
		log.Fatalln("packet error: ", err)
	}

	_, err = conn.Write(packet)
	if err != nil {
		log.Fatalln("write error: ", err)
	}


}
