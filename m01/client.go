package main

import (
	"context"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
}

func (c *Client) On(eventName string, cb func(ctx context.Context, msg []byte)) {

}
