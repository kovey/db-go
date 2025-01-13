package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql"
)

const (
	DateTimeMill = "2006-01-02 15:04:05.000"
)

var Engine ksql.EngineInterface = sql.DefaultEngine()

type LogInfo struct {
	start     int64  // ms
	end       int64  // ms
	Delay     string `json:"delay"`
	TraceId   string `json:"trace_id"`
	BeginTime string `json:"begin_time"`
	EndTime   string `json:"end_time"`
	Sql       string `json:"sql"`
}

func NewLogInfo() *LogInfo {
	return &LogInfo{}
}

func (l *LogInfo) Start(traceId string) {
	now := time.Now()
	l.start = now.UnixMicro()
	l.BeginTime = now.Format(DateTimeMill)
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
	l.end = now.UnixMicro()
	l.EndTime = now.Format(DateTimeMill)
	l.Delay = fmt.Sprintf("%.3fms", float64(l.end-l.start)*0.001)
}

func (l *LogInfo) Encode() []byte {
	buffer := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buffer)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(l); err == nil {
		return buffer.Bytes()
	}

	return nil
}
