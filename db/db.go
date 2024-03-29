package db

import (
	"context"
	ds "database/sql"
	"fmt"

	"github.com/kovey/db-go/v2/itf"
	"github.com/kovey/db-go/v2/sql"
	"github.com/kovey/db-go/v2/sql/meta"
	"github.com/kovey/debug-go/debug"
	"github.com/kovey/pool/object"
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
	FetchPage(table string, where meta.Where, model T, page int, pageSize int, orders ...string) (*meta.Page[T], error)
	FetchPageByWhere(table string, where sql.WhereInterface, model T, page int, pageSize int, orders ...string) (*meta.Page[T], error)
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
	FetchPageCtx(ctx context.Context, table string, where meta.Where, model T, page int, pageSize int, orders ...string) (*meta.Page[T], error)
	FetchPageByWhereCtx(ctx context.Context, table string, where sql.WhereInterface, model T, page int, pageSize int, orders ...string) (*meta.Page[T], error)
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
	stmt, err := m.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	data, err := stmt.QueryContext(ctx, args...)
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
	cc, ok := ctx.(object.CtxInterface)
	if !ok {
		cc = nil
	}

	for data.Next() {
		tmp := model.Clone(cc)
		if err := data.Scan(getFields(columns, has, length, tmp)...); err != nil {
			return nil, err
		}

		if tt, ok := tmp.(itf.ModelInterface); ok {
			tt.SetFetch()
		}

		rows = append(rows, tmp.(T))
	}

	return rows, nil
}

func Query[T itf.RowInterface](ctx context.Context, m ConnInterface, query string, model T, args ...any) ([]T, error) {
	stmt, err := m.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	data, err := stmt.QueryContext(ctx, args...)
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
	cc, ok := ctx.(object.CtxInterface)
	if !ok {
		cc = nil
	}
	for data.Next() {
		tmp := model.Clone(cc)
		if err := data.Scan(getFields(columns, has, length, tmp)...); err != nil {
			return nil, err
		}

		if tt, ok := tmp.(itf.ModelInterface); ok {
			tt.SetFetch()
		}

		rows = append(rows, tmp.(T))
	}

	return rows, nil
}

func Exec(ctx context.Context, m ConnInterface, statment string) error {
	stmt, err := m.PrepareContext(ctx, statment)
	if err != nil {
		debug.Erro("exec prepare error: %s, sql: %s", err, statment)
		return err
	}

	defer stmt.Close()
	result, err := stmt.ExecContext(ctx)
	if err != nil {
		debug.Erro("exec error: %s, sql: %s", err, statment)
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
		debug.Erro("insert error: %s, sql: %s", err, insert)
		return 0, err
	}

	return result.LastInsertId()
}

func Update(ctx context.Context, m ConnInterface, update *sql.Update) (int64, error) {
	result, err := prepare(ctx, m, update)
	if err != nil {
		debug.Erro("update error: %s, sql: %s", err, update)
		return 0, err
	}

	return result.RowsAffected()
}

func Delete(ctx context.Context, m ConnInterface, del *sql.Delete) (int64, error) {
	result, err := prepare(ctx, m, del)
	if err != nil {
		debug.Erro("delete error: %s, sql: %s", err, del)
		return 0, err
	}

	return result.RowsAffected()
}

func BatchInsert(ctx context.Context, m ConnInterface, batch *sql.Batch) (int64, error) {
	result, err := prepare(ctx, m, batch)
	if err != nil {
		debug.Erro("batch insert error: %s, sql: %s", err, batch)
		return 0, err
	}

	return result.RowsAffected()
}

func ShowTables[T itf.RowInterface](ctx context.Context, m ConnInterface, show *sql.ShowTables, model T) ([]T, error) {
	rows, err := queryAll(ctx, m, show.Prepare(), model, show.Args()...)
	if err != nil {
		debug.Erro("show tables error: %s, sql: %s", err, show)
	}

	return rows, err
}

func Desc[T itf.RowInterface](ctx context.Context, m ConnInterface, desc *sql.Desc, model T) ([]T, error) {
	rows, err := queryAll(ctx, m, desc.Prepare(), model, desc.Args()...)
	if err != nil {
		debug.Erro("desc error: %s, sql: %s", err, desc)
	}

	return rows, err
}

func Select[T itf.RowInterface](ctx context.Context, m ConnInterface, sel *sql.Select, model T) ([]T, error) {
	rows, err := Query(ctx, m, sel.Prepare(), model, sel.Args()...)
	if err != nil {
		debug.Erro("select error: %s, sql: %s", err, sel)
	}
	return rows, err
}

func FetchRow[T itf.ModelInterface](ctx context.Context, m ConnInterface, table string, where meta.Where, model T) error {
	var sel *sql.Select
	if cc, ok := ctx.(object.CtxInterface); ok {
		sel = sql.NewSelectBy(cc, table, emptyStr)
	} else {
		sel = sql.NewSelect(table, emptyStr)
	}
	sel.WhereByMap(where).Columns(model.Columns()...).Limit(1)
	stmt, err := m.PrepareContext(ctx, sel.Prepare())
	if err != nil {
		debug.Erro("fetch row error: %s, sql: %s", err, sel)
		return err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx, sel.Args()...)
	err = parseError(result.Scan(model.Fields()...), model)
	if err != nil {
		debug.Erro("fetch row error: %s, sql: %s", err, sel)
	}
	return err
}

func parseError[T itf.ModelInterface](err error, model T) error {
	if err == nil {
		model.SetFetch()
		return nil
	}

	model.SetEmpty()
	if err == ds.ErrNoRows {
		return nil
	}

	return err
}

func LockRow[T itf.ModelInterface](ctx context.Context, m ConnInterface, table string, where meta.Where, model T) error {
	var sel *sql.Select
	if cc, ok := ctx.(object.CtxInterface); ok {
		sel = sql.NewSelectBy(cc, table, emptyStr)
	} else {
		sel = sql.NewSelect(table, emptyStr)
	}
	sel.WhereByMap(where).Columns(model.Columns()...).Limit(1).ForUpdate()
	stmt, err := m.PrepareContext(ctx, sel.Prepare())
	if err != nil {
		debug.Erro("lock row error: %s, sql: %s", err, sel)
		return err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx, sel.Args()...)
	err = parseError(result.Scan(model.Fields()...), model)
	if err != nil {
		debug.Erro("lock row error: %s, sql: %s", err, sel)
	}

	return err
}

func FetchAll[T itf.RowInterface](ctx context.Context, m ConnInterface, table string, where meta.Where, model T) ([]T, error) {
	var sel *sql.Select
	if cc, ok := ctx.(object.CtxInterface); ok {
		sel = sql.NewSelectBy(cc, table, emptyStr)
	} else {
		sel = sql.NewSelect(table, emptyStr)
	}
	sel.WhereByMap(where).Columns(model.Columns()...)

	rows, err := queryAll(ctx, m, sel.Prepare(), model, sel.Args()...)
	if err != nil {
		debug.Erro("fetch all error: %s, sql: %s", err, sel)
	}

	return rows, err
}

func FetchBySelect[T itf.RowInterface](ctx context.Context, m ConnInterface, sel *sql.Select, model T) ([]T, error) {
	sel.SetColumns(model.Columns())

	rows, err := queryAll(ctx, m, sel.Prepare(), model, sel.Args()...)
	if err != nil {
		debug.Erro("fetch by select error: %s, sql: %s", err, sel)
	}

	return rows, err
}

func FetchAllByWhere[T itf.RowInterface](ctx context.Context, m ConnInterface, table string, where sql.WhereInterface, model T) ([]T, error) {
	var sel *sql.Select
	if cc, ok := ctx.(object.CtxInterface); ok {
		sel = sql.NewSelectBy(cc, table, emptyStr)
	} else {
		sel = sql.NewSelect(table, emptyStr)
	}
	sel.Where(where).Columns(model.Columns()...)

	rows, err := queryAll(ctx, m, sel.Prepare(), model, sel.Args()...)
	if err != nil {
		debug.Erro("fetch all by where error: %s, sql: %s", err, sel)
	}

	return rows, err
}

func FetchPage[T itf.RowInterface](ctx context.Context, m ConnInterface, table string, where meta.Where, model T, page int, pageSize int, orders ...string) (*meta.Page[T], error) {
	var sel *sql.Select
	if cc, ok := ctx.(object.CtxInterface); ok {
		sel = sql.NewSelectBy(cc, table, emptyStr)
	} else {
		sel = sql.NewSelect(table, emptyStr)
	}
	sel.WhereByMap(where).Columns(model.Columns()...).Limit(pageSize).Offset((page - 1) * pageSize).Order(orders...)

	rows, err := queryAll(ctx, m, sel.Prepare(), model, sel.Args()...)
	if err != nil {
		debug.Erro("fetch page error: %s, sql: %s", err, sel)
		return nil, err
	}

	w := sql.NewWhere()
	for key, val := range where {
		w.Eq(key, val)
	}

	return pageInfo(ctx, m, table, w, rows, pageSize, page)
}

func FetchPageByWhere[T itf.RowInterface](ctx context.Context, m ConnInterface, table string, where sql.WhereInterface, model T, page int, pageSize int, orders ...string) (*meta.Page[T], error) {
	var sel *sql.Select
	if cc, ok := ctx.(object.CtxInterface); ok {
		sel = sql.NewSelectBy(cc, table, emptyStr)
	} else {
		sel = sql.NewSelect(table, emptyStr)
	}
	sel.Where(where).Columns(model.Columns()...).Limit(pageSize).Offset((page - 1) * pageSize).Order(orders...)

	rows, err := queryAll(ctx, m, sel.Prepare(), model, sel.Args()...)
	if err != nil {
		debug.Erro("fetch page by where error: %s, sql: %s", err, sel)
		return nil, err
	}

	return pageInfo(ctx, m, table, where, rows, pageSize, page)
}

func FetchPageBySelect[T itf.RowInterface](ctx context.Context, m ConnInterface, sel *sql.Select, model T) (*meta.Page[T], error) {
	rows, err := queryAll(ctx, m, sel.Prepare(), model, sel.Args()...)
	if err != nil {
		debug.Erro("fetch page by select error: %s, sql: %s", err, sel)
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
	stmt, err := m.PrepareContext(ctx, sel.Prepare())
	if err != nil {
		debug.Erro("fetch page by select count error: %s, sql: %s", err, sel)
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, sel.Args()...)
	count := int64(0)
	err = row.Scan(&count)
	if err != nil {
		debug.Erro("fetch page by select count error: %s, sql: %s", err, sel)
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
	var sel *sql.Select
	if cc, ok := ctx.(object.CtxInterface); ok {
		sel = sql.NewSelectBy(cc, table, emptyStr)
	} else {
		sel = sql.NewSelect(table, emptyStr)
	}
	sel.Where(where).ColMeta(meta.NewColFuncWithNull(meta.NewField(one, emptyStr, true), countField, zero, meta.Func_COUNT, nil))
	stmt, err := m.PrepareContext(ctx, sel.Prepare())
	if err != nil {
		debug.Erro("count error: %s, sql: %s", err, sel)
		return 0, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, sel.Args()...)
	count := int64(0)
	err = row.Scan(&count)
	if err != nil {
		debug.Erro("count error: %s, sql: %s", err, sel)
	}

	return count, err
}
