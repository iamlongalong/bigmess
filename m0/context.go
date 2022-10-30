package main

import "context"

type ContextOpt struct {
	ctx    context.Context
	logger ILogger
}

func NewContext(client *Client, opt ContextOpt) *Context {
	if opt.ctx == nil {
		opt.ctx = context.Background()
	}

	ctx, cancle := context.WithCancel(opt.ctx)

	return &Context{
		Context: ctx,
		cancle:  cancle,
		client:  client,
		logger:  opt.logger,
	}
}

type Context struct {
	context.Context
	cancle context.CancelFunc

	client *Client
	logger ILogger
}

// 姑且把 Client 实例返回，之后可以看是否需要做内部保护
func (c *Context) GetClient() *Client {
	return c.client
}

func (c *Context) GetLogger() ILogger {
	return c.logger
}

func (c *Context) Value(key any) any {
	v := c.Context.Value(key)
	if v != nil {
		return v
	}

	v, ok := c.client.Get(key)
	if ok {
		return v
	}

	return nil
}

func (c *Context) Stop() {
	c.cancle()
}
