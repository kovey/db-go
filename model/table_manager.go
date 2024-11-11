package model

import (
	"context"
	"fmt"
	"time"

	"github.com/kovey/db-go/v3/db"
	"github.com/kovey/db-go/v3/sql"
)

type TableType byte

const (
	Table_Type_Day   TableType = 1
	Table_Type_Month TableType = 2
)

const (
	Day_Format   = "20060102"
	Month_Format = "200601"
)

type Template struct {
	Table         string
	Keep          int
	Type          TableType
	TemplateTable string
}

type TableManager struct {
	templates []*Template
	now       time.Time
}

func NewTableManager(templates []*Template) *TableManager {
	return &TableManager{templates: templates, now: time.Now()}
}

func (t *TableManager) Append(temp *Template) {
	t.templates = append(t.templates, temp)
}

func (t *TableManager) Create(ctx context.Context) {
	for _, temp := range t.templates {
		switch temp.Type {
		case Table_Type_Day:
			t.createDay(ctx, temp)
		case Table_Type_Month:
			t.createMonth(ctx, temp)
		}
	}
}

func (t *TableManager) createMonth(ctx context.Context, temp *Template) {
	for i := -3; i <= 3; i++ {
		tableName := fmt.Sprintf("%s_%s", temp.Table, t.now.AddDate(0, i, 0).Format(Month_Format))
		if has, err := db.HasTable(ctx, tableName); err != nil || has {
			continue
		}

		if _, err := db.ExecRaw(ctx, sql.Raw("CREATE TABLE ? LIKE ?", tableName, temp.TemplateTable)); err != nil {
			fmt.Println(err)
		}
	}
}

func (t *TableManager) createDay(ctx context.Context, temp *Template) {
	for i := -3; i <= 3; i++ {
		tableName := fmt.Sprintf("%s_%s", temp.Table, t.now.AddDate(0, 0, i).Format(Day_Format))
		if has, err := db.HasTable(ctx, tableName); err != nil || has {
			continue
		}

		if _, err := db.ExecRaw(ctx, sql.Raw("CREATE TABLE ? LIKE ?", tableName, temp.TemplateTable)); err != nil {
			fmt.Println(err)
		}
	}
}

func (t *TableManager) delDay(ctx context.Context, temp *Template) {
	for i := 0; i < 3; i++ {
		tableName := fmt.Sprintf("%s_%s", temp.Table, t.now.AddDate(0, 0, -(temp.Keep+i)).Format(Day_Format))
		if has, err := db.HasTable(ctx, tableName); err != nil || !has {
			continue
		}

		if err := db.DropTable(ctx, tableName); err != nil {
			fmt.Println(err)
		}
	}
}

func (t *TableManager) delMonth(ctx context.Context, temp *Template) {
	for i := 0; i < 3; i++ {
		tableName := fmt.Sprintf("%s_%s", temp.Table, t.now.AddDate(0, -(temp.Keep+i), 0).Format(Month_Format))
		if has, err := db.HasTable(ctx, tableName); err != nil || !has {
			continue
		}

		if err := db.DropTable(ctx, tableName); err != nil {
			fmt.Println(err)
		}
	}
}

func (t *TableManager) Delete(ctx context.Context, temp *Template) {
	for _, temp := range t.templates {
		switch temp.Type {
		case Table_Type_Day:
			t.delDay(ctx, temp)
		case Table_Type_Month:
			t.delMonth(ctx, temp)
		}
	}
}
