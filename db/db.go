package db

import (
	"context"
	ds "database/sql"
	"fmt"

	"github.com/kovey/db-go/v2/itf"
	"github.com/kovey/db-go/v2/sql"
	"github.com/kovey/db-go/v2/sql/meta"
	"github.com/kovey/debug-go/debug"
)

const (
	countField = "count"
	emptyStr   = ""
	one        = "1"
	zero       = "0"
)

type DbInterface[T itf.RowInterface] interface {
	SetTx(*Tx)
	Transaction(func(*Tx) error) error
	InTransaction() bool
	Query(string, T, ...any) ([]T, error)
	Exec(string) error
	Desc(*sql.Desc, T) ([]T, error)
	ShowTables(*sql.ShowTables, T) ([]T, error)
	Insert(*sql.Insert) (int64, error)
	Update(*sql.Update) (int64, error)
	Delete(*sql.Delete) (int64, error)
	BatchInsert(*sql.Batch) (int64, error)
	Select(*sql.Select, T) ([]T, error)
	FetchRow(string, meta.Where, T) error
	LockRow(string, meta.Where, T) error
	FetchAll(string, meta.Where, T) ([]T, error)
	FetchAllByWhere(string, sql.WhereInterface, T) ([]T, error)
	FetchBySelect(*sql.Select, T) ([]T, error)
	FetchPage(string, meta.Where, T, int, int) (*meta.Page[T], error)
	FetchPageByWhere(string, sql.WhereInterface, T, int, int) (*meta.Page[T], error)
	FetchPageBySelect(*sql.Select, T) (*meta.Page[T], error)
	Count(string, sql.WhereInterface) (int64, error)
	TransactionCtx(context.Context, func(*Tx) error, *ds.TxOptions) error
	QueryCtx(context.Context, string, T, ...any) ([]T, error)
	ExecCtx(context.Context, string) error
	DescCtx(context.Context, *sql.Desc, T) ([]T, error)
	ShowTablesCtx(context.Context, *sql.ShowTables, T) ([]T, error)
	InsertCtx(context.Context, *sql.Insert) (int64, error)
	UpdateCtx(context.Context, *sql.Update) (int64, error)
	DeleteCtx(context.Context, *sql.Delete) (int64, error)
	BatchInsertCtx(context.Context, *sql.Batch) (int64, error)
	SelectCtx(context.Context, *sql.Select, T) ([]T, error)
	FetchRowCtx(context.Context, string, meta.Where, T) error
	LockRowCtx(context.Context, string, meta.Where, T) error
	FetchAllCtx(context.Context, string, meta.Where, T) ([]T, error)
	FetchAllByWhereCtx(context.Context, string, sql.WhereInterface, T) ([]T, error)
	FetchBySelectCtx(context.Context, *sql.Select, T) ([]T, error)
	FetchPageCtx(context.Context, string, meta.Where, T, int, int) (*meta.Page[T], error)
	FetchPageByWhereCtx(context.Context, string, sql.WhereInterface, T, int, int) (*meta.Page[T], error)
	FetchPageBySelectCtx(context.Context, *sql.Select, T) (*meta.Page[T], error)
	CountCtx(context.Context, string, sql.WhereInterface) (int64, error)
}

type ConnInterface interface {
	Query(string, ...any) (*ds.Rows, error)
	QueryContext(context.Context, string, ...any) (*ds.Rows, error)
	Exec(string, ...any) (ds.Result, error)
	ExecContext(context.Context, string, ...any) (ds.Result, error)
	Prepare(string) (*ds.Stmt, error)
	PrepareContext(context.Context, string) (*ds.Stmt, error)
	QueryRow(string, ...any) *ds.Row
	QueryRowContext(context.Context, string, ...any) *ds.Row
}

func min(left, right int) int {
	if left < right {
		return left
	}

	return right
}

func getFields[T itf.RowInterface](columns []string, has []string, length int, model T) []any {
	fields := make([]any, length)
	index := 0
	mFields := model.Fields()
	for _, column := range columns {
		for i, col := range has {
			if column == col {
				fields[index] = mFields[i]
				index++
				break
			}
		}
	}
	return fields
}

func queryAll[T itf.RowInterface](ctx context.Context, m ConnInterface, query string, model T, args ...any) ([]T, error) {
	data, err := m.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer data.Close()
	rows := make([]T, 0)
	for data.Next() {
		tmp := model.Clone()
		if err := data.Scan(tmp.Fields()...); err != nil {
			return nil, err
		}

		rows = append(rows, tmp.(T))
	}

	return rows, nil
}

func Query[T itf.RowInterface](ctx context.Context, m ConnInterface, query string, model T, args ...any) ([]T, error) {
	data, err := m.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer data.Close()
	columns, err := data.Columns()
	if err != nil {
		return nil, err
	}
	has := model.Columns()
	length := min(len(columns), len(has))
	rows := make([]T, 0)
	for data.Next() {
		tmp := model.Clone()
		if err := data.Scan(getFields(columns, has, length, tmp)...); err != nil {
			return nil, err
		}

		rows = append(rows, tmp.(T))
	}

	return rows, nil
}

func Exec(ctx context.Context, m ConnInterface, stament string) error {
	debug.Info("sql: %s", stament)
	result, err := m.ExecContext(ctx, stament)

	if err != nil {
		return err
	}

	lastId, _ := result.LastInsertId()
	affected, _ := result.RowsAffected()

	if lastId < 1 && affected < 1 {
		return fmt.Errorf("lastId or affectedId is zero")
	}

	return nil
}

func prepare(ctx context.Context, m ConnInterface, pre sql.SqlInterface) (ds.Result, error) {
	debug.Info("sql: %s", pre)
	smt, err := m.PrepareContext(ctx, pre.Prepare())
	if err != nil {
		return nil, err
	}

	defer smt.Close()

	return smt.Exec(pre.Args()...)
}

func Insert(ctx context.Context, m ConnInterface, insert *sql.Insert) (int64, error) {
	result, err := prepare(ctx, m, insert)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func Update(ctx context.Context, m ConnInterface, update *sql.Update) (int64, error) {
	result, err := prepare(ctx, m, update)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func Delete(ctx context.Context, m ConnInterface, del *sql.Delete) (int64, error) {
	result, err := prepare(ctx, m, del)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func BatchInsert(ctx context.Context, m ConnInterface, batch *sql.Batch) (int64, error) {
	result, err := prepare(ctx, m, batch)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func ShowTables[T itf.RowInterface](ctx context.Context, m ConnInterface, show *sql.ShowTables, model T) ([]T, error) {
	return queryAll(ctx, m, show.Prepare(), model, show.Args()...)
}

func Desc[T itf.RowInterface](ctx context.Context, m ConnInterface, desc *sql.Desc, model T) ([]T, error) {
	return queryAll(ctx, m, desc.Prepare(), model, desc.Args()...)
}

func Select[T itf.RowInterface](ctx context.Context, m ConnInterface, sel *sql.Select, model T) ([]T, error) {
	return Query(ctx, m, sel.Prepare(), model, sel.Args()...)
}

func FetchRow[T itf.ModelInterface](ctx context.Context, m ConnInterface, table string, where meta.Where, model T) error {
	sel := sql.NewSelect(table, "")
	sel.WhereByMap(where).Columns(model.Columns()...).Limit(1)
	result := m.QueryRowContext(ctx, sel.Prepare(), sel.Args()...)

	return parseError(result.Scan(model.Fields()...), model)
}

func parseError[T itf.ModelInterface](err error, model T) error {
	if err == nil {
		return nil
	}

	if err == ds.ErrNoRows {
		model.SetEmpty()
		return nil
	}

	return err
}

func LockRow[T itf.ModelInterface](ctx context.Context, m ConnInterface, table string, where meta.Where, model T) error {
	sel := sql.NewSelect(table, "")
	sel.WhereByMap(where).Columns(model.Columns()...).Limit(1).ForUpdate()
	result := m.QueryRowContext(ctx, sel.Prepare(), sel.Args()...)

	return parseError(result.Scan(model.Fields()...), model)
}

func FetchAll[T itf.RowInterface](ctx context.Context, m ConnInterface, table string, where meta.Where, model T) ([]T, error) {
	sel := sql.NewSelect(table, "")
	sel.WhereByMap(where).Columns(model.Columns()...)

	return queryAll(ctx, m, sel.Prepare(), model, sel.Args()...)
}

func FetchBySelect[T itf.RowInterface](ctx context.Context, m ConnInterface, sel *sql.Select, model T) ([]T, error) {
	sel.SetColumns(model.Columns())

	return queryAll(ctx, m, sel.Prepare(), model, sel.Args()...)
}

func FetchAllByWhere[T itf.RowInterface](ctx context.Context, m ConnInterface, table string, where sql.WhereInterface, model T) ([]T, error) {
	sel := sql.NewSelect(table, "")
	sel.Where(where).Columns(model.Columns()...)

	return queryAll(ctx, m, sel.Prepare(), model, sel.Args()...)
}

func FetchPage[T itf.RowInterface](ctx context.Context, m ConnInterface, table string, where meta.Where, model T, page int, pageSize int) (*meta.Page[T], error) {
	sel := sql.NewSelect(table, "")
	sel.WhereByMap(where).Columns(model.Columns()...).Limit(pageSize).Offset((page - 1) * pageSize)

	rows, err := queryAll(ctx, m, sel.Prepare(), model, sel.Args()...)
	if err != nil {
		return nil, err
	}

	w := sql.NewWhere()
	for key, val := range where {
		w.Eq(key, val)
	}

	return pageInfo(ctx, m, table, w, rows, pageSize, page)
}

func FetchPageByWhere[T itf.RowInterface](ctx context.Context, m ConnInterface, table string, where sql.WhereInterface, model T, page int, pageSize int) (*meta.Page[T], error) {
	sel := sql.NewSelect(table, "")
	sel.Where(where).Columns(model.Columns()...).Limit(pageSize).Offset((page - 1) * pageSize)

	rows, err := queryAll(ctx, m, sel.Prepare(), model, sel.Args()...)
	if err != nil {
		return nil, err
	}

	return pageInfo(ctx, m, table, where, rows, pageSize, page)
}

func FetchPageBySelect[T itf.RowInterface](ctx context.Context, m ConnInterface, sel *sql.Select, model T) (*meta.Page[T], error) {
	rows, err := queryAll(ctx, m, sel.Prepare(), model, sel.Args()...)
	if err != nil {
		return nil, err
	}

	pageSize := sel.GetLimit()
	if sel.GetOffset() == 0 && len(rows) < pageSize {
		p := meta.NewPage(rows)
		p.TotalCount = int64(len(rows))
		p.TotalPage = 1
		return p, nil
	}

	sel.Limit(1)
	sel.Offset(0)
	cols := make([]*meta.Column, 1)
	cols[0] = meta.NewColFuncWithNull(meta.NewField(one, emptyStr, true), countField, zero, meta.Func_COUNT, nil)
	sel.SetColMeta(cols)
	row := m.QueryRow(sel.Prepare(), sel.Args()...)
	count := int64(0)
	err = row.Scan(&count)
	if err != nil {
		return nil, err
	}

	p := meta.NewPage(rows)
	p.TotalCount = count
	p.TotalPage = p.TotalCount / int64(pageSize)
	if p.TotalCount%int64(pageSize) != 0 {
		p.TotalPage++
	}

	return p, nil
}

func pageInfo[T itf.RowInterface](ctx context.Context, m ConnInterface, table string, where sql.WhereInterface, rows []T, pageSize, curPage int) (*meta.Page[T], error) {
	if curPage == 1 && len(rows) < pageSize {
		page := meta.NewPage(rows)
		page.TotalCount = int64(len(rows))
		page.TotalPage = 1
		return page, nil
	}

	count, err := Count(ctx, m, table, where)
	if err != nil {
		return nil, err
	}

	p := meta.NewPage(rows)
	p.TotalCount = count
	p.TotalPage = p.TotalCount / int64(pageSize)
	if p.TotalCount%int64(pageSize) != 0 {
		p.TotalPage++
	}

	return p, nil
}

func Count(ctx context.Context, m ConnInterface, table string, where sql.WhereInterface) (int64, error) {
	sel := sql.NewSelect(table, emptyStr)
	sel.Where(where).ColMeta(meta.NewColFuncWithNull(meta.NewField(one, emptyStr, true), countField, zero, meta.Func_COUNT, nil))
	row := m.QueryRowContext(ctx, sel.Prepare(), sel.Args()...)
	count := int64(0)
	err := row.Scan(&count)

	return count, err
}
