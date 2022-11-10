package main

import "context"

var DfHandlers = Handlers{}

func init() {
	DfHandlers.Register("subchange", HandleSubChanges)
}

type Handlers struct {
	handlers map[string]HandlerFunc
}

type HandlerFunc func(ctx context.Context, c *Client, msg []byte)

func (h *Handlers) Handle(ctx context.Context, eventName string, c *Client, msg []byte) {
	f, ok := h.handlers[eventName]
	if ok {
		f(ctx, c, msg)
	}
}

func (h *Handlers) Register(eventName string, cb func(ctx context.Context, c *Client, msg []byte)) {
	h.handlers[eventName] = cb
}

func HandleSubChanges(ctx context.Context, c *Client, msg []byte) {

}
