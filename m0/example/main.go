package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	url = "ws://127.0.0.1:8080/ws"
)

var (
	MsgCodeJoinRoom = "/fileroom/join"
	MsgCodePub      = "/fileroom/pub"
)

type Message struct {
	MessageCode string                 `json:"msgcode"`
	MsgHeader   map[string]interface{} `json:"header"`
	MsgBody     interface{}            `json:"body"`
}

func main() {
	roomID := ""

	if len(os.Args) <= 1 {
		fmt.Println("please input roomid : ")
		fmt.Scan(&roomID)
	} else {
		roomID = os.Args[1]
	}

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Printf("dial %s fail : %s \n", url, err)
	}

	c := &Client{
		conn: conn,
	}
	c.ctx, c.cancel = context.WithCancel(context.Background())

	go c.Read()
	go c.HeartBeat()

	err = c.JoinRoom(roomID)
	if err != nil {
		log.Printf("join room : %s fail : %s", roomID, err)
		return
	}

	defer c.Close()

	firstFlag := true

	for {
		select {
		case <-c.ctx.Done():
			fmt.Println("exit")
			return
		default:
			if firstFlag {
				fmt.Println("please write something you want to publish: ")
				firstFlag = false
			}

			inputReader := bufio.NewReader(os.Stdin)
			input, _, err := inputReader.ReadLine()
			if err != nil {
				fmt.Printf("read msg from stdin fail : %s\n", err)
				return
			}

			err = c.PubToRoom(roomID, string(input))
			if err != nil {
				fmt.Printf("write msg fail : %s\n", err)
				return
			}
		}
	}
}

type Client struct {
	mu     sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc

	conn *websocket.Conn
}

func (c *Client) HeartBeat() {
	t := time.NewTicker(time.Second * 5)
	defer func() {
		c.Close()
		t.Stop()
	}()

	for {
		select {
		case <-t.C:
			c.mu.Lock()
			err := c.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				log.Printf("write ping message fail : %s", err)
				c.mu.Unlock()
				return
			}
			c.mu.Unlock()

		case <-c.ctx.Done():
			return
		}
	}
}

func (c *Client) Read() {
	defer func() {
		c.Close()
	}()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			_, b, err := c.conn.ReadMessage()
			if err != nil {
				log.Printf("read message fail : %s", err)
				return
			}

			msg := &Message{}
			err = json.Unmarshal(b, msg)
			if err != nil {
				log.Printf("unmarshal fail : %s", err)
				return
			}

			switch msg.MessageCode {
			case "/communicate": // 聊天消息
				from := msg.MsgHeader["fromCid"]
				msgStr := msg.MsgBody

				fmt.Printf("%s: %s\n", from, msgStr)
			default:
				fmt.Printf("get msg : %s \n", string(b))
			}

		}
	}
}

func (c *Client) Close() {
	c.cancel()
	c.conn.Close()
}

func (c *Client) JoinRoom(roomID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	msg := &Message{
		MessageCode: MsgCodeJoinRoom,
		MsgBody: map[string]interface{}{
			"room_id": roomID,
		},
	}

	return c.conn.WriteJSON(msg)
}

func (c *Client) PubToRoom(roomID string, msg string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	pubmsg := &Message{
		MessageCode: MsgCodePub,
		MsgBody: map[string]interface{}{
			"topic": "/communicate",
			"data":  msg,
		},
	}

	return c.conn.WriteJSON(pubmsg)
}
