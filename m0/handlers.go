package main

import "github.com/go-playground/validator/v10"

var dfvalidator = validator.New()

const (
	JOINED_ROOM_KEY = "__JOINED_ROOM"
)

type MsgPubToRoom struct {
	RoomID   string `json:"room_id" validate:"required"`
	PubTopic string `validate:"required"`
	Data     map[string]interface{}
}

// ctx 中应该包含 client 对象
func HandlePubToRoom(ctx *Context, m IMessage) {
	params := &MsgPubToRoom{}
	err := m.DecodeBody(NewJsonDecoder(params))
	if err != nil {
		// 看是否回消息？
		ctx.GetLogger().Printf(LevelError, "decode body fail : %s", err)
		return
	}

	if err := dfvalidator.Struct(params); err != nil {
		ctx.GetLogger().Printf(LevelError, "validate params fail : %s", err)
		return
	}

	b, _ := NewJsonEncoder(params.Data).Encode()

	// 消息
	sendMsg := &JsonMessage{
		MsgCode:   MessageCode(params.PubTopic),
		MsgHeader: MessageHeader{},
		MsgBody:   b,
	}

	err = DefaultRoomHub.BoardCast(sendMsg, nil, params.RoomID)
	if err != nil {
		ctx.GetLogger().Printf(LevelError, "board cast fail : %s", err)
		return
	}
}

type MsgJoinRoom struct {
	RoomID string `json:"room_id" validate:"required"`
}

func HandleJoinRoom(ctx *Context, m IMessage) {
	params := &MsgJoinRoom{}
	if err := m.DecodeBody(NewJsonDecoder(params)); err != nil {
		ctx.GetLogger().Printf(LevelError, "decode body fail : %s", err)
		return
	}

	if err := dfvalidator.Struct(params); err != nil {
		ctx.GetLogger().Printf(LevelError, "validate params fail : %s", err)
		return
	}

	r, err := DefaultRoomHub.GetOrCreateRoom(params.RoomID)
	if err != nil {
		ctx.GetLogger().Printf(LevelError, "get or create room fail : %s", err)
		return
	}

	if err = r.JoinRoom(ctx.GetClient()); err != nil {
		ctx.GetLogger().Printf(LevelError, "join room fail : %s", err)
		return
	}

	ctx.GetClient().Set(JOINED_ROOM_KEY, r)

}

type MsgPubMsg struct {
	Topic string      `json:"topic" validate:"required"`
	Data  interface{} `json:"data" validate:"required"`
}

func HandlePubMsg(ctx *Context, m IMessage) {
	params := &MsgPubMsg{}
	if err := m.DecodeBody(NewJsonDecoder(params)); err != nil {
		ctx.GetLogger().Printf(LevelError, "decode body fail : %s", err)
		return
	}

	if err := dfvalidator.Struct(params); err != nil {
		ctx.GetLogger().Printf(LevelError, "validate params fail : %s", err)
		return
	}

	r := ctx.Value(JOINED_ROOM_KEY)
	if r == nil {
		ctx.GetLogger().Printf(LevelError, "no joined room")
		return
	}

	b, _ := NewJsonEncoder(params.Data).Encode()

	// 消息
	sendMsg := &JsonMessage{
		MsgCode:   MessageCode(params.Topic),
		MsgHeader: MessageHeader{"fromCid": ctx.GetClient().ID()},
		MsgBody:   b,
	}

	if err := r.(*Room).BoardCast(sendMsg, &BoardCastOpt{ExceptIDs: []string{ctx.GetClient().ID()}}); err != nil {
		ctx.GetLogger().Printf(LevelError, "board cast fail : %s", err)
		return
	}
}
