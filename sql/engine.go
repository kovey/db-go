package sql

import (
	"fmt"
	"strings"
	"time"

	ksql "github.com/kovey/db-go/v3"
)

type StringInterface interface {
	String() string
}

type Engine struct {
	quote      string
	timeFormat string
}

func NewEngine(quote, timeFormat string) *Engine {
	return &Engine{quote: quote, timeFormat: timeFormat}
}

func DefaultEngine() *Engine {
	return NewEngine("'", time.DateTime)
}

func (e *Engine) Format(sql ksql.SqlInterface) string {
	tpl := sql.Prepare()
	for _, bind := range sql.Binds() {
		tpl = strings.Replace(tpl, "?", e.value(bind), 1)
	}

	return tpl
}

func (e *Engine) FormatRaw(sql ksql.ExpressInterface) string {
	tpl := sql.Statement()
	for _, bind := range sql.Binds() {
		tpl = strings.Replace(tpl, "?", e.value(bind), 1)
	}

	return tpl
}

func (e *Engine) value(val any) string {
	switch tmp := val.(type) {
	case string:
		return fmt.Sprintf("%s%s%s", e.quote, tmp, e.quote)
	case *string:
		return fmt.Sprintf("%s%s%s", e.quote, *tmp, e.quote)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", tmp)
	case float32, float64:
		return fmt.Sprintf("%d", tmp)
	case *int:
		return fmt.Sprintf("%d", *tmp)
	case *int8:
		return fmt.Sprintf("%d", *tmp)
	case *int16:
		return fmt.Sprintf("%d", *tmp)
	case *int32:
		return fmt.Sprintf("%d", *tmp)
	case *int64:
		return fmt.Sprintf("%d", *tmp)
	case *uint:
		return fmt.Sprintf("%d", *tmp)
	case *uint8:
		return fmt.Sprintf("%d", *tmp)
	case *uint16:
		return fmt.Sprintf("%d", *tmp)
	case *uint32:
		return fmt.Sprintf("%d", *tmp)
	case *uint64:
		return fmt.Sprintf("%d", *tmp)
	case *float32:
		return fmt.Sprintf("%f", *tmp)
	case *float64:
		return fmt.Sprintf("%f", *tmp)
	case time.Time:
		return tmp.Format(e.timeFormat)
	case *time.Time:
		return tmp.Format(e.timeFormat)
	case bool:
		return fmt.Sprintf("%t", tmp)
	case *bool:
		return fmt.Sprintf("%t", *tmp)
	default:
		if s, ok := val.(StringInterface); ok {
			return fmt.Sprintf("%s%s%s", e.quote, s.String(), e.quote)
		}

		return fmt.Sprintf("%v", tmp)
	}
}
