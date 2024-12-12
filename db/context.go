package db

import (
	"context"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/logger"
)

type Context struct {
	context.Context
	logInfo *logger.LogInfo
	traceId string
}

func NewContext(ctx context.Context) ksql.ContextInterface {
	if tmp, ok := ctx.(ksql.ContextInterface); ok {
		return tmp
	}

	c := &Context{Context: ctx}
	if trace, ok := ctx.(ksql.TraceInterface); ok {
		c.traceId = trace.TraceId()
	}

	return c
}

func (c *Context) SqlLogStart(sql ksql.SqlInterface) {
	if !logOpen {
		return
	}

	c.logInfo = logger.NewLogInfo()
	c.logInfo.Start(c.traceId)
	c.logInfo.ExecSql(sql)
}

func (c *Context) RawSqlLogStart(sql ksql.ExpressInterface) {
	if !logOpen {
		return
	}

	c.logInfo = logger.NewLogInfo()
	c.logInfo.Start(c.traceId)
	c.logInfo.ExecRawSql(sql)
}

func (c *Context) SqlLogEnd() {
	if c.logInfo == nil {
		return
	}

	c.logInfo.End()
	logger.Append(c.logInfo.Encode())
	c.logInfo = nil
}
