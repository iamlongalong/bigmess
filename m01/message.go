package main

import "encoding/json"

type Event struct {
	Name    string
	Message json.RawMessage
}

type ChangeStateMessage struct {
	Option string
	Key    string
	Val    interface{}
}

func ParseMessage(msg json.RawMessage, p interface{}) error {
	return json.Unmarshal(msg, p)
}

func ParseEvent(msg []byte) (*Event, error) {
	e := &Event{}
	return e, json.Unmarshal(msg, e)
}

func EncodeJsonMessage(eventName string, msg interface{}) (b []byte, err error) {
	mb, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return json.Marshal(Event{
		Name:    eventName,
		Message: mb,
	})
}
