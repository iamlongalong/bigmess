package main

import "context"

type Context struct {
	context.Context

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
