package main

import (
	"context"
	"log"

	rmap "github.com/iamlongalong/bigmess/pkg/responsivemap"
)

type Handlers struct {
	handlers map[string]HandlerFunc
}

type IHandler interface {
	Handle(ctx context.Context, c *Peer, event *Event)
}

type HandlerFunc func(ctx context.Context, c *Peer, msg []byte)

func (h *Handlers) Handle(ctx context.Context, c *Peer, event *Event) {
	f, ok := h.handlers[event.Name]
	if ok {
		f(ctx, c, event.Message)
	}
}

func (h *Handlers) Register(eventName string, cb func(ctx context.Context, c *Peer, msg []byte)) {
	h.handlers[eventName] = cb
}

// HandleSubChanges 收到其他连接发来的状态更新
func (clu *Cluster) HandleStateChanges(ctx context.Context, c *Peer, msg []byte) {
	m := &ChangeStateMessage{}
	err := ParseMessage(msg, m)
	if err != nil {
		log.Printf("parse subchanges msg fail : %s", err)
		return
	}

	switch m.Option {
	case rmap.MapKeyGet:
		return
	case rmap.MapKeySet:
		err = c.states.Set(m.Key, m.Val, rmap.Option{DisableNotice: true})
		if err != nil {
			log.Printf("set value fail in sub changes : %s", err)
			return
		}
	case rmap.MapKeyDel:
		err = c.states.Del(m.Key, rmap.Option{DisableNotice: true})
		if err != nil {
			log.Printf("del value fail in sub changes : %s", err)
			return
		}
	}
}

// HandleClusterStateChanges 收到其他连接发来的状态更新
func (clu *Cluster) HandleClusterStateChanges(ctx context.Context, c *Peer, msg []byte) {
	m := &ChangeStateMessage{}
	err := ParseMessage(msg, m)
	if err != nil {
		log.Printf("parse subchanges msg fail : %s", err)
		return
	}

	switch m.Option {
	case rmap.MapKeyGet:
		return
	case rmap.MapKeySet:
		err = clu.cinfo.Set(m.Key, m.Val, rmap.Option{DisableNotice: true})
		if err != nil {
			log.Printf("set value fail in sub changes : %s", err)
			return
		}
	case rmap.MapKeyDel:
		err = clu.cinfo.Del(m.Key, rmap.Option{DisableNotice: true})
		if err != nil {
			log.Printf("del value fail in sub changes : %s", err)
			return
		}
	}
}
