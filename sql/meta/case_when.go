package meta

import (
	"fmt"
	"strings"
)

const (
	when_then       = "WHEN %s THEN %s"
	else_condition  = "ELSE %s"
	case_expression = "CASE"
	end             = "END"
	case_when       = "(%s %s %s) AS `%s`"
	case_when_else  = "(%s %s %s %s) AS `%s`"
)

type CaseWhen struct {
	Conditions []*WhenThen
	elseResult string
	Alias      string
}

func NewCaseWhen(alias string) *CaseWhen {
	return &CaseWhen{Alias: alias, Conditions: make([]*WhenThen, 0)}
}

func (c *CaseWhen) AddWhenThen(when, then string) {
	c.Conditions = append(c.Conditions, NewWhenThen(when, then))
}

func (c *CaseWhen) Else(expression string) {
	c.elseResult = expression
}

func (c *CaseWhen) String() string {
	if c.elseResult == "" {
		return fmt.Sprintf(case_when, case_expression, c.whens(), end, c.Alias)
	}

	return fmt.Sprintf(case_when_else, case_expression, c.whens(), fmt.Sprintf(else_condition, c.elseResult), end, c.Alias)
}

func (c *CaseWhen) whens() string {
	res := make([]string, len(c.Conditions))
	for index, co := range c.Conditions {
		res[index] = co.String()
	}

	return strings.Join(res, " ")
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
