package db

import (
	"context"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/logger"
)

type Context struct {
	context.Context
	logInfo    []*logger.LogInfo
	isStarting bool
	traceId    string
	endIndex   int
}

func NewContext(ctx context.Context) ksql.ContextInterface {
	if tmp, ok := ctx.(ksql.ContextInterface); ok {
		return tmp
	}

	c := &Context{Context: ctx, endIndex: -1}
	if trace, ok := ctx.(ksql.TraceInterface); ok {
		c.traceId = trace.TraceId()
	}

	return c
}

func (c *Context) WithTraceId(traceId string) ksql.ContextInterface {
	if c.isStarting {
		return c
	}

	c.traceId = traceId
	return c
}

func (c *Context) SqlLogStart(sql ksql.SqlInterface) {
	if !logOpen {
		return
	}

	if !c.isStarting {
		c.isStarting = true
	}

	info := logger.NewLogInfo()
	info.Start(c.traceId)
	info.ExecSql(sql)
	c.logInfo = append(c.logInfo, info)
	c.endIndex++
}

func (c *Context) RawSqlLogStart(sql ksql.ExpressInterface) {
	if !logOpen {
		return
	}

	info := logger.NewLogInfo()
	info.Start(c.traceId)
	info.ExecRawSql(sql)
	c.logInfo = append(c.logInfo, info)
	c.endIndex++
}

func (c *Context) reset() {
	c.isStarting = false
	c.endIndex = -1
	c.logInfo = nil
}

func (c *Context) SqlLogEnd() {
	if c.endIndex < 0 {
		return
	}

	info := c.logInfo[c.endIndex]
	info.End()
	logger.Append(info.Encode())
	c.endIndex--
	if c.endIndex < 0 {
		c.reset()
	}
}
