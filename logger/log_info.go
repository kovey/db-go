package logger

import (
	"encoding/json"
	"fmt"
	"time"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql"
)

var Engine ksql.EngineInterface = sql.DefaultEngine()

type LogInfo struct {
	start     int64  // ms
	end       int64  // ms
	TraceId   string `json:"trace_id"`
	BeginTime string `json:"begin_time"`
	EndTime   string `json:"end_time"`
	Delay     string `json:"delay"`
	Sql       string `json:"sql"`
}

func NewLogInfo() *LogInfo {
	return &LogInfo{}
}

func (l *LogInfo) Start(traceId string) {
	now := time.Now()
	l.start = now.UnixMilli()
	l.BeginTime = now.Format(time.StampMilli)
	l.TraceId = traceId
}

func (l *LogInfo) ExecSql(s ksql.SqlInterface) {
	l.Sql = Engine.Format(s)
}

func (l *LogInfo) ExecRawSql(s ksql.ExpressInterface) {
	l.Sql = Engine.FormatRaw(s)
}

func (l *LogInfo) End() {
	now := time.Now()
	l.end = now.UnixMilli()
	l.EndTime = now.Format(time.StampMilli)
	l.Delay = fmt.Sprintf("%.3fms", float64(l.end-l.start)*0.001)
}

func (l *LogInfo) Encode() []byte {
	if logBytes, err := json.Marshal(l); err == nil {
		return logBytes
	}

	return nil
}
