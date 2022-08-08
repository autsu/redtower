package server

import (
	"errors"
	"fmt"
)

type Router struct {
	m map[[16]byte]Handler
}

func NewRouter() *Router {
	// 添加内置的几个消息类型
	m := map[[16]byte]Handler{
		ErrorMsg.ID():     &errorHandler{},
		HeartBeatMsg.ID(): &heartBeatHandler{},
		OriginalMsg.ID():  &originalHandler{},
	}
	return &Router{m: m}
}

func (r *Router) Do(req *Request) error {
	handler, ok := r.m[req.MsgType().ID()]
	//log.Printf("typ id: %x, typ name: %v\n", req.MsgType().ID(), md5Typ[req.MsgType().ID()])
	if !ok {
		errmsg := fmt.Sprintf("unknown message type: %v", md5Typ[req.MsgType().ID()])
		return errors.New(errmsg)
	}

	handler.BeforeHandle(req)
	handler.Handle(req)
	handler.AfterHandle(req)

	return nil
}

func (r *Router) Add(t MessageType, handler Handler) {
	if t == ErrorMsg || t == HeartBeatMsg || t == OriginalMsg {
		panic("the type you define is the same as the system type (error, heartbeat, original)")
	}
	//log.Printf("add a new typ: %v\n", t)
	// 这里似乎不用加锁，因为 Add 是在程序启动前调用的
	r.m[t.ID()] = handler
}
