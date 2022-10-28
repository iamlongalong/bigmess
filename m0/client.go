package main

import (
	"context"
	"log"

	"github.com/gorilla/websocket"
)

func NewClient(conn *websocket.Conn, handler IHandler) *Client {
	return &Client{
		msgHandler: handler,
		readChan:   make(chan *Message, 20),
		writeChan:  make(chan *Message, 20),
		conn:       conn,
	}
}

type Client struct {
	conn *websocket.Conn

	readChan  chan *Message
	writeChan chan *Message

	msgHandler IHandler
}

func (c *Client) Start() {
	go c.startRead()
	go c.startWrite()
	go c.startHandle()
}

func (c *Client) startHandle() {
	for msg := range c.readChan {
		ctx := &Context{
			Context: context.Background(),
			client:  c,
			logger:  &Logger{},
		}
		c.msgHandler.Handle(ctx, msg)
	}
}

func (c *Client) startRead() {
	defer func() {
		c.Close()
	}()

	for {
		msg, err := c.read()
		if err != nil {
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
		b, err := EncodeMessage(msg)
		if err != nil {
			log.Printf("encode message fail : %s", err)
		}

		err = c.conn.WriteMessage(websocket.BinaryMessage, b)
		if err != nil {
			// TODO 判断错误情况
			log.Printf("write message fail : %s", err)
			return
		}
	}
}

func (c *Client) Send(msg *Message) error {
	// TODO 先做一些检查之类的

	c.writeChan <- msg
	return nil
}

func (c *Client) read() (*Message, error) {
	_, msgbytes, err := c.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	return DecodeMessage(msgbytes)
}

func (c *Client) Close() error {
	return c.conn.Close()
}
