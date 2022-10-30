package main

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type MessageCode string

type MessageHeader map[string]interface{}

type Decoder interface {
	Decode([]byte) error
}

type Encoder interface {
	Encode() ([]byte, error)
}

type IMessage interface {
	Encoder

	MessageCode() MessageCode
	Header() MessageHeader
	DecodeBody(d Decoder) error
}

func NewJsonDecoder(tar interface{}) Decoder {
	return &JsonEncoderDecoder{tar: tar}
}
func NewJsonEncoder(tar interface{}) Encoder {
	return &JsonEncoderDecoder{tar: tar}
}

// JsonEncoderDecoder 通用的 json decoder
type JsonEncoderDecoder struct {
	tar interface{}
}

func (j *JsonEncoderDecoder) Decode(d []byte) error {
	return json.Unmarshal(d, j.tar)
}

func (j *JsonEncoderDecoder) Encode() ([]byte, error) {
	return json.Marshal(j.tar)
}

type JsonMessage struct {
	// 消息码姑且用 string 表示
	MsgCode MessageCode `json:"msgcode"`

	// header 姑且用 map 表示
	MsgHeader MessageHeader `json:"header"`

	// body 姑且保持最原始状态
	MsgBody json.RawMessage `json:"body"`
}

func (j *JsonMessage) Encode() ([]byte, error) {
	return json.Marshal(j)
}

func (j *JsonMessage) MessageCode() MessageCode {
	return j.MsgCode
}
func (j *JsonMessage) Header() MessageHeader {
	return j.MsgHeader
}
func (j *JsonMessage) DecodeBody(d Decoder) error {
	return d.Decode(j.MsgBody)
}

func DecodeJsonMessage(d []byte) (IMessage, error) {
	m := &JsonMessage{}
	err := json.Unmarshal(d, m)
	if err != nil {
		return nil, errors.Wrap(err, "decode json message fail")
	}

	return m, nil
}
