package main

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

func NewClient(conn *websocket.Conn, handler IHandler) *Client {
	return &Client{
		id: GetID(), // id 姑且用人名

		msgHandler: handler,
		readChan:   make(chan IMessage, 20),
		writeChan:  make(chan IMessage, 20),
		conn:       conn,

		closedHooks: make([]func(c *Client) error, 0),
	}
}

type Client struct {
	id   string
	conn *websocket.Conn

	readChan  chan IMessage
	writeChan chan IMessage

	vals sync.Map

	msgHandler IHandler

	closedHooks []func(c *Client) error
}

func (c *Client) AddClosedHook(hook func(c *Client) error) {
	c.closedHooks = append(c.closedHooks, hook)
}

func (c *Client) Set(k, v interface{}) {
	c.vals.Store(k, v)
}

func (c *Client) Get(k interface{}) (interface{}, bool) {
	return c.vals.Load(k)
}

func (c *Client) ID() string {
	return c.id
}

func (c *Client) Start() {
	go c.startRead()
	go c.startWrite()
	go c.startHandle()
}

func (c *Client) startHandle() {
	for msg := range c.readChan {
		logger := NewLogger(KV{Key: "clientid", Value: c.id}, KV{Key: "msgcode", Value: string(msg.MessageCode())})
		ctx := NewContext(c, ContextOpt{logger: logger})

		c.msgHandler.Handle(ctx, msg)
	}
}

func (c *Client) startRead() {
	defer func() {
		c.Close()
	}()

	for {
		// err := c.conn.SetReadDeadline(time.Now().Add(time.Second * 30))
		// if err != nil {
		// 	log.Printf("set read deadline fail : %s", err)
		// 	return
		// }

		msg, err := c.read()
		if err != nil {
			log.Printf("read message fail : %s", err)
			return
		}

		c.readChan <- msg
	}
}

func (c *Client) startWrite() {
	defer func() {
		c.Close()
	}()

	for msg := range c.writeChan { // send msg
		b, err := msg.Encode()
		if err != nil {
			log.Printf("encode message fail : %s", err)
		}

		err = c.conn.SetWriteDeadline(time.Now().Add(time.Second * 30))
		if err != nil {
			log.Printf("set write deadline fail : %s", err)
			return
		}

		err = c.conn.WriteMessage(websocket.BinaryMessage, b)
		if err != nil {
			// TODO 判断错误情况
			log.Printf("write message fail : %s", err)
			return
		}
	}
}

func (c *Client) Send(msg IMessage) error {
	// TODO 先做一些检查之类的

	c.writeChan <- msg
	return nil
}

func (c *Client) read() (IMessage, error) {
	_, msgbytes, err := c.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	// 姑且用 json 的格式
	return DecodeJsonMessage(msgbytes)
}

func (c *Client) Close() error {
	errs := &BundleErr{}
	for _, hook := range c.closedHooks {
		errs.AddErr(hook(c))
	}

	errs.AddErr(c.conn.Close())

	return errs.Bundle()
}
