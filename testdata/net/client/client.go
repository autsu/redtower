package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testdata/net"
	"testdata/pbfile"

	"github.com/autsu/redtower/client"
	"github.com/autsu/redtower/server"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func main() {
	home := os.Getenv("HOME")
	fp := filepath.Join(home, "Downloads/jDKUARa.jpg")

	f, err := os.Open(fp)
	assert(err, "open file error: %v")
	defer f.Close()

	file, err := os.ReadFile(fp)
	assert(err, "read file error: %v")

	msg := pbfile.Message{
		Uid:         wrapperspb.UInt64(10086),
		Data:        file,
		MessageType: wrapperspb.String("image"),
	}

	c, err := client.NewClientWithTCP("localhost", "7788").Init(context.Background())
	assert(err, "dial error: %v")
	defer c.Close()

	b, err := proto.Marshal(&msg)
	assert(err, "proto marshal error: %v")
	fmt.Printf("proto marshal size: %v\n", len(b))

	n, err := c.Send(server.NewMessage(b, net.ProtoMsg))
	assert(err, "write to conn error: %v")
	fmt.Printf("write %v bytes\n", n)
}

func assert(err error, format string) {
	if err != nil {
		panic(fmt.Sprintf(format, err.Error()))
	}
}
