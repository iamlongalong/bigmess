package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	rmap "github.com/iamlongalong/bigmess/pkg/responsivemap"
)

func NewPeer(id string, ipaddr string, conn *websocket.Conn, handler IHandler) *Peer {
	ctx, cancel := context.WithCancel(context.Background())

	m := rmap.NewResponsiveMap()
	c := &Peer{
		ctx:     ctx,
		cancel:  cancel,
		conn:    conn,
		handler: handler,
		id:      id,
		ipaddr:  ipaddr,
		states:  m,
	}

	// TODO
	// 1. 先 sync states
	// 2. 再 监听 states

	// create a key listener
	listener := rmap.NewMapListener("", func(me rmap.MapEvent) {
		c.ChangeState(context.Background(), &ChangeStateMessage{
			Option: string(me.Option),
			Key:    me.Key,
			Val:    me.NewVal,
		})
	})
	m.Watch(listener)

	go c.Start()

	return c
}

type Peer struct {
	ctx    context.Context
	cancel context.CancelFunc

	once sync.Once

	conn *websocket.Conn

	handler IHandler

	states *rmap.ResponsiveMap

	id string

	ipaddr string
}

func (c *Peer) Start() {
	go c.read()
	go c.heartbeat()
}

func (c *Peer) heartbeat() {
	t := time.NewTicker(time.Second * 30)
	defer t.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-t.C:
			w, err := c.conn.NextWriter(websocket.PingMessage)
			if err != nil {
				log.Printf("get next writer fail in heartbeat : %s", err)
				return
			}

			_, err = w.Write(nil)
			if err != nil {
				log.Printf("write heartbeat fail : %s", err)
				return
			}
		}

	}
}
func (c *Peer) read() {
	defer func() {
		c.Close()
	}()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			_, msg, err := c.conn.ReadMessage()
			if err != nil {
				log.Printf("read message fail : %s", err)
				return
			}

			e, err := ParseEvent(msg)
			if err != nil {
				log.Printf("parse event fail : %s", err)
				continue
			}

			c.handler.Handle(context.Background(), c, e)
		}
	}
}

func (c *Peer) write(ctx context.Context, msg []byte) error {
	w, err := c.conn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	return err
}

func (c *Peer) Close() {
	c.once.Do(func() {
		c.conn.Close()
		c.cancel()
	})
}

// ChangeState 改变自己的 state 后的触发
func (c *Peer) ChangeState(ctx context.Context, msg *ChangeStateMessage) {
	b, err := EncodeJsonMessage(EventStateChange, msg)
	if err != nil {
		log.Printf("encode msg fail in changeState : %s", err)
		return
	}

	err = c.write(ctx, b)
	if err != nil {
		log.Printf("write changeState fail : %s", err)
		return
	}
}
