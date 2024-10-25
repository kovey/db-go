package table

import (
	"fmt"
	"strings"
)

type ScaleType byte

const (
	Scale_Type_Zero ScaleType = 1
	Scale_Type_One  ScaleType = 2
	Scale_Type_Two  ScaleType = 3
	Scale_Type_More ScaleType = 4
)

type ColumnType struct {
	Name   string
	Type   ScaleType
	Length int
	Scale  int
	sets   []string
}

func (c *ColumnType) Set(sets ...string) *ColumnType {
	for _, set := range sets {
		c.sets = append(c.sets, fmt.Sprintf("'%s'", set))
	}

	return c
}

func (c *ColumnType) Express() string {
	switch c.Type {
	case Scale_Type_One:
		return fmt.Sprintf("%s(%d)", c.Name, c.Length)
	case Scale_Type_Two:
		return fmt.Sprintf("%s(%d,%d)", c.Name, c.Length, c.Scale)
	case Scale_Type_More:
		return fmt.Sprintf("%s(%s)", c.Name, strings.Join(c.sets, ","))
	default:
		return c.Name
	}
}
