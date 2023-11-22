package meta

import (
	"fmt"
	"strings"

	"github.com/kovey/pool"
	"github.com/kovey/pool/object"
)

const (
	when_then       = "WHEN %s THEN %s"
	else_condition  = "ELSE %s"
	case_expression = "CASE"
	end             = "END"
	case_when       = "(%s %s %s) AS `%s`"
	case_when_else  = "(%s %s %s %s) AS `%s`"
	space           = " "
	emptyStr        = ""
	as              = "%s AS %s"
	using           = "%s USING %s"
	qua             = "''"
	strFormat       = "`%s`"
	strTowFormat    = "`%s`.`%s`"
	namespace       = "ko.db.sql.meta"
	cw_name         = "CaseWhen"
)

func init() {
	pool.DefaultNoCtx(namespace, cw_name, func() any {
		return &CaseWhen{ObjNoCtx: object.NewObjNoCtx(namespace, cw_name)}
	})
}

type CaseWhen struct {
	*object.ObjNoCtx
	Conditions []*WhenThen
	elseResult string
	Alias      string
}

func NewCaseWhen(alias string) *CaseWhen {
	return &CaseWhen{Alias: alias, Conditions: make([]*WhenThen, 0)}
}

func NewCaseWhenBy(ctx object.CtxInterface, alias string) *CaseWhen {
	obj := ctx.GetNoCtx(namespace, cw_name).(*CaseWhen)
	obj.Alias = alias
	return obj
}

func (c *CaseWhen) Reset() {
	c.Alias = emptyStr
	c.elseResult = emptyStr
	c.Conditions = nil
}

func (c *CaseWhen) AddWhenThen(when, then string) {
	c.Conditions = append(c.Conditions, NewWhenThen(when, then))
}

func (c *CaseWhen) Else(expression string) {
	c.elseResult = expression
}

func (c *CaseWhen) String() string {
	if c.elseResult == emptyStr {
		return fmt.Sprintf(case_when, case_expression, c.whens(), end, c.Alias)
	}

	return fmt.Sprintf(case_when_else, case_expression, c.whens(), fmt.Sprintf(else_condition, c.elseResult), end, c.Alias)
}

func (c *CaseWhen) whens() string {
	res := make([]string, len(c.Conditions))
	for index, co := range c.Conditions {
		res[index] = co.String()
	}

	return strings.Join(res, space)
}

type WhenThen struct {
	When string
	Then string
}

func NewWhenThen(when, then string) *WhenThen {
	return &WhenThen{When: when, Then: then}
}

func (w *WhenThen) String() string {
	return fmt.Sprintf(when_then, w.When, w.Then)
}
