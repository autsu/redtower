package server

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	m := make(map[MessageType]struct{})
	g1 := GenMsgTyp("123")
	g2 := GenMsgTyp("456")
	g3 := GenMsgTyp("123")
	m[g1] = struct{}{}
	m[g2] = struct{}{}
	m[g3] = struct{}{}

	fmt.Println(len(m))

	if v, ok := m[g1]; ok {
		fmt.Println(v)
	} else {
		t.Fatal("error")
	}
}
