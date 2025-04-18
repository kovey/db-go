package table

import (
	"fmt"
	"strconv"
	"strings"
)

type ScaleType byte

const (
	Scale_Type_Zero ScaleType = 0x1
	Scale_Type_One  ScaleType = 0x2
	Scale_Type_Two  ScaleType = 0x3
	Scale_Type_More ScaleType = 0x4
)

type ColumnType struct {
	Name   string
	Type   ScaleType
	Length int
	Scale  int
	sets   []string
}

func (c *ColumnType) IsNumeric() bool {
	return isNumeric(c.Name)
}

func (c *ColumnType) IsInteger() bool {
	return isInteger(c.Name)
}

func (c *ColumnType) Set(sets ...string) *ColumnType {
	for _, set := range sets {
		c.sets = append(c.sets, fmt.Sprintf("'%s'", set))
	}

	return c
}

func (c *ColumnType) Build(builder *strings.Builder) {
	builder.WriteString(" ")
	builder.WriteString(c.Name)
	switch c.Type {
	case Scale_Type_One:
		builder.WriteString("(")
		builder.WriteString(strconv.Itoa(c.Length))
		builder.WriteString(")")
	case Scale_Type_Two:
		builder.WriteString("(")
		builder.WriteString(strconv.Itoa(c.Length))
		builder.WriteString(",")
		builder.WriteString(strconv.Itoa(c.Scale))
		builder.WriteString(")")
	case Scale_Type_More:
		builder.WriteString("(")
		builder.WriteString(strings.Join(c.sets, ","))
		builder.WriteString(")")
	}
}
