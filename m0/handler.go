package main

import (
	"log"

	"github.com/pkg/errors"
)

var (
	ErrRegisterSameMsgCode = errors.New("registry same MsgCode")
	ErrNoHandlerFound      = errors.New("no handler found")
)

func NewEngine() IEngine {
	return &WsEngine{
		IRouters: &WsRouter{},
	}
}

type WsEngine struct {
	IRouters
}

type IEngine interface {
	IRouters
}

func NewHandler(f HandleFunc) Handler {
	return Handler{f: f}
}

type Handler struct {
	f HandleFunc
}

func (h *Handler) Handle(ctx *Context, m *Message) {
	h.f(ctx, m)
}

type IHandler interface {
	Handle(ctx *Context, m *Message)
}

type HandleFunc func(ctx *Context, m *Message)

type IRouters interface {
	Register(MessageCode, ...Handler) // 注册 handlers
	Handle(*Context, *Message)        // 处理消息
}

// router 是用于做消息分发的，消息码
type WsRouter struct {
	// 这里姑且就不加读写锁了，请求访问量大时会有损耗，在启动的时候加增加写锁就行了。
	freezed bool // 启动冻结锁

	NotFoundHandler IHandler

	// 分发方案先用最简单的 hashmap 实现，最简单
	routes map[MessageCode][]Handler
}

func (r *WsRouter) Handle(ctx *Context, m *Message) {
	defer func() { // 没有分组 middleware，recover 姑且写在这里
		if err := recover(); err != nil {
			log.Print(err) // 可以考虑添加堆栈信息等
		}
	}()

	hs := r.getHandlers(m.MessageCode)
	if len(hs) == 0 && r.NotFoundHandler != nil {
		r.NotFoundHandler.Handle(ctx, m)
	}

	for _, h := range hs {
		h.Handle(ctx, m)
	}
}

func (r *WsRouter) Register(mc MessageCode, hs ...Handler) {
	if r.freezed {
		panic("can not registry hanlder func when running")
	}

	_, ok := r.routes[mc]
	if ok {
		panic(errors.Wrapf(ErrRegisterSameMsgCode, string(mc)))
	}

	r.routes[mc] = hs
}

func (r *WsRouter) getHandlers(mc MessageCode) []Handler {
	return r.routes[mc]
}
