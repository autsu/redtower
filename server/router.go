package server

import (
	"errors"
)

type Router struct {
	m map[MessageType]Handler
}

func NewRouter() *Router {
	return &Router{m: make(map[MessageType]Handler)}
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
	r.m[t] = handler
}
