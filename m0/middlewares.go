package main

import "log"

var AuthKey = "__auth_key_m0"

type MiddleWareAuth struct{}

func (mwa *MiddleWareAuth) Handle(ctx *Context, m IMessage) {
	v := ctx.Value(AuthKey)
	if v == nil {
		log.Printf("not authed")
	}
}

type MiddleWareLog struct{}

func (mwl *MiddleWareLog) Handle(ctx *Context, m IMessage) {
	ctx.GetLogger().Printf(LevelDebug, "get message")
}
