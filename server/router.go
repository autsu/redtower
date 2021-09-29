package server

import (
	"errors"
)

type Router struct {
	m map[MessageType]Handler
}

func NewRouter() *Router {
	m := map[MessageType]Handler{
		ErrorMsg: &BasicHandler{},
		HeartBeatMsg: NewHeartBeatHandler(),
		OriginalMsg: &BasicHandler{},
	}
	return &Router{m: m}
}

func (r *Router) Do(req *Request) error {
	handler, ok := r.m[req.MsgType()]
	if !ok {
		return errors.New("unknown message type")
	}

	handler.BeforeHandle(req)
	handler.Handle(req)
	handler.AfterHandle(req)

	return nil
}

func (r *Router) AddRouter(t MessageType, handler Handler) {
	if t == ErrorMsg || t == HeartBeatMsg || t == OriginalMsg {
		panic("the type you define is the same as the system type (error, heartbeat, original)")
	}
	r.m[t] = handler
}
