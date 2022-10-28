package main

import "encoding/json"

type MessageCode string

type MessageHeader map[string]interface{}

type Message struct {
	// 消息码姑且用 string 表示
	MessageCode MessageCode `json:"msgcode"`

	// header 姑且用 map 表示
	Header MessageHeader `json:"header"`

	// body 姑且保持最原始状态
	Body []byte `json:"body"`
}

func DecodeMessage(d []byte) (*Message, error) {
	msg := &Message{}

	return msg, json.Unmarshal(d, &msg)
}

func EncodeMessage(msg *Message) ([]byte, error) {
	return json.Marshal(msg)
}
