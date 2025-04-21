package sql

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type columnInfo struct {
	fun    string
	column string
	as     string
	prefix string
	isFunc bool
	expr   ksql.ExpressInterface
}

func (c *columnInfo) Binds() []any {
	if c.expr == nil {
		return nil
	}

	return c.expr.Binds()
}

type columnInfos struct {
	columns []*columnInfo
	binds   []any
}

func (c *columnInfos) Empty() bool {
	return len(c.columns) == 0
}

func (c *columnInfos) Build(builder *strings.Builder) {
	for index, column := range c.columns {
		if index > 0 {
			builder.WriteString(", ")
		}
		column.Build(builder)
		c.binds = append(c.binds, column.Binds()...)
	}
}

func (c *columnInfos) Append(column *columnInfo) {
	c.columns = append(c.columns, column)
}

func (c *columnInfo) Build(builder *strings.Builder) {
	if c.expr != nil {
		builder.WriteString("(")
		builder.WriteString(c.expr.Statement())
		builder.WriteString(")")
	} else if c.isFunc {
		builder.WriteString(c.fun)
		builder.WriteString("(")
		if c.prefix != "" {
			builder.WriteString(c.prefix)
			operator.BuildColumnString(c.column, builder)
		} else {
			operator.Column(c.column, builder)
		}
		builder.WriteString(")")
	} else {
		builder.WriteString(c.prefix)
		if c.prefix != "" {
			operator.BuildColumnString(c.column, builder)
		} else {
			operator.Column(c.column, builder)
		}
	}

	if c.as != "" {
		builder.WriteString(" AS")
		operator.BuildColumnString(c.as, builder)
	}
}

type orderMeta struct {
	column *columnInfo
	typ    string
}

func (o *orderMeta) Build(builder *strings.Builder) {
	o.column.Build(builder)
	builder.WriteString(" ")
	builder.WriteString(o.typ)
}

type orderInfo struct {
	columns []*orderMeta
	with    string
}

func (o *orderInfo) Append(column *orderMeta) {
	o.columns = append(o.columns, column)
}

func (o *orderInfo) Empty() bool {
	return len(o.columns) == 0
}

func (o *orderInfo) Build(builder *strings.Builder) {
	builder.WriteString(" ORDER BY ")
	for index, column := range o.columns {
		if index > 0 {
			builder.WriteString(", ")
		}
		column.Build(builder)
	}

	if o.with != "" {
		builder.WriteString(", ")
		builder.WriteString(o.with)
	}
}

type groupInfo struct {
	columns *columnInfos
	with    string
}

func newGroupInfo() *groupInfo {
	return &groupInfo{columns: &columnInfos{}}
}

func (o *groupInfo) Empty() bool {
	return o.columns.Empty()
}

func (o *groupInfo) Build(builder *strings.Builder) {
	builder.WriteString("GROUP BY ")
	o.columns.Build(builder)
	if o.with != "" {
		builder.WriteString(", ")
		builder.WriteString(o.with)
	}
}

type tableInfo struct {
	table string
	sub   ksql.QueryInterface
	as    string
}

func (t *tableInfo) Build(builder *strings.Builder) {
	if t.sub != nil {
		builder.WriteString(" (")
		builder.WriteString(t.sub.Prepare())
		builder.WriteString(")")
	} else {
		operator.BuildColumnString(t.table, builder)
	}

	if t.as != "" {
		builder.WriteString(" AS")
		operator.BuildColumnString(t.as, builder)
	}
}

func (t *tableInfo) Binds() []any {
	if t.sub == nil {
		return nil
	}

	return t.sub.Binds()
}

type window struct {
	name string
	as   string
}

func (w *window) Build(builder *strings.Builder) {
	operator.BuildColumnString(w.name, builder)
	builder.WriteString(" AS")
	operator.BuildColumnString(w.as, builder)
}

type windows struct {
	data []*window
}

func (w *windows) Empty() bool {
	return len(w.data) == 0
}

func (w *windows) Append(wd *window) {
	w.data = append(w.data, wd)
}

func (w *windows) Build(builder *strings.Builder) {
	builder.WriteString("WINDOW")
	for index, window := range w.data {
		if index > 0 {
			builder.WriteString(",")
		}
		window.Build(builder)
	}
}

type limitInfo struct {
	hasLimit  bool
	limit     int
	hasOffset bool
	offset    int
}

func (l *limitInfo) Binds() []any {
	var res []any
	if l.hasLimit {
		res = append(res, l.limit)
	}
	if l.hasOffset {
		res = append(res, l.offset)
	}

	return res
}

func (l *limitInfo) Build(builder *strings.Builder) {
	if l.hasLimit {
		builder.WriteString(" LIMIT ?")
	}

	if l.hasOffset {
		builder.WriteString(" OFFSET ?")
	}
}

type Query struct {
	*base
	modifer          string
	table            *tableInfo
	columns          *columnInfos
	intoVars         []string
	where            ksql.WhereInterface
	join             []ksql.JoinInterface
	limitInfo        *limitInfo
	order            *orderInfo
	group            *groupInfo
	having           ksql.HavingInterface
	sharding         ksql.Sharding
	initBinds        []any
	forSql           *For
	partitions       []string
	highPriority     string
	straightJoin     string
	sqlSmallResult   string
	sqlBigResult     string
	sqlBufferResult  string
	sqlNoCache       string
	sqlCalcFoundRows string
	windows          *windows
}

func NewQuery() *Query {
	q := &Query{
		base: newBase(), where: NewWhere(), having: NewHaving(), sharding: ksql.Sharding_None, columns: &columnInfos{}, group: newGroupInfo(), order: &orderInfo{},
		forSql: &For{}, table: &tableInfo{}, windows: &windows{}, limitInfo: &limitInfo{},
	}
	q.opChain.Append(q._keyword, q._columns, q._into, q._from, q._joinInfo, q._partition, q._where, q._group, q._having, q._window, q._order, q._limit, q._for)
	return q
}

func (o *Query) _keyword(builder *strings.Builder) {
	builder.WriteString("SELECT")
	operator.BuildPureString(o.modifer, builder)
	operator.BuildPureString(o.highPriority, builder)
	operator.BuildPureString(o.straightJoin, builder)
	operator.BuildPureString(o.sqlSmallResult, builder)
	operator.BuildPureString(o.sqlBigResult, builder)
	operator.BuildPureString(o.sqlBufferResult, builder)
	operator.BuildPureString(o.sqlNoCache, builder)
	operator.BuildPureString(o.sqlCalcFoundRows, builder)
}

func (o *Query) _columns(builder *strings.Builder) {
	if o.columns.Empty() {
		builder.WriteString(" *")
		return
	}

	builder.WriteString(" ")
	o.columns.Build(builder)
	o.binds = append(o.binds, o.columns.binds...)
}

func (o *Query) _into(builder *strings.Builder) {
	if len(o.intoVars) == 0 {
		return
	}

	builder.WriteString(" INTO ")
	for index, intoVar := range o.intoVars {
		if index > 0 {
			builder.WriteString(",")
		}

		operator.BuildColumnString(intoVar, builder)
	}
}

func (o *Query) _from(builder *strings.Builder) {
	builder.WriteString(" FROM")
	o.table.Build(builder)
	o.binds = append(o.binds, o.table.Binds()...)
}

func (o *Query) _joinInfo(builder *strings.Builder) {
	if len(o.join) == 0 {
		return
	}

	for _, join := range o.join {
		builder.WriteString(" ")
		join.Build(builder)
		o.binds = append(o.binds, join.Binds()...)
	}
}

func (o *Query) _partition(builder *strings.Builder) {
	if len(o.partitions) == 0 {
		return
	}

	builder.WriteString(" PARTITION")
	for index, partition := range o.partitions {
		if index > 0 {
			builder.WriteString(",")
		}

		operator.BuildColumnString(partition, builder)
	}
}

func (o *Query) _where(builder *strings.Builder) {
	if o.where.Empty() {
		return
	}

	builder.WriteString(" ")
	o.where.Build(builder)
	o.binds = append(o.binds, o.where.Binds()...)
}

func (o *Query) _group(builder *strings.Builder) {
	if o.group.Empty() {
		return
	}

	builder.WriteString(" ")
	o.group.Build(builder)
}

func (o *Query) _having(builder *strings.Builder) {
	if o.having.Empty() {
		return
	}

	builder.WriteString(" ")
	o.having.Build(builder)
	o.binds = append(o.binds, o.having.Binds()...)
}

func (o *Query) _window(builder *strings.Builder) {
	if o.windows.Empty() {
		return
	}

	builder.WriteString(" ")
	o.windows.Build(builder)
}

func (o *Query) _order(builder *strings.Builder) {
	if o.order.Empty() {
		return
	}

	o.order.Build(builder)
}

func (o *Query) _limit(builder *strings.Builder) {
	if o.limitInfo.hasLimit || o.limitInfo.hasOffset {
		o.limitInfo.Build(builder)
		o.binds = append(o.binds, o.limitInfo.Binds()...)
	}
}

func (o *Query) _for(builder *strings.Builder) {
	if o.forSql.Empty() {
		return
	}

	o.forSql.Build(builder)
}

func (o *Query) Clone() ksql.QueryInterface {
	q := &Query{
		base: newBase(), where: o.where.Clone(), having: o.having.Clone(),
		table: o.table, join: o.join, group: o.group, initBinds: o.initBinds, order: o.order, intoVars: o.intoVars, limitInfo: o.limitInfo, modifer: o.modifer,
		forSql: o.forSql, partitions: o.partitions, highPriority: o.highPriority, straightJoin: o.straightJoin, windows: o.windows, columns: &columnInfos{},
	}
	q.opChain.Append(q._keyword, q._columns, q._into, q._from, q._joinInfo, q._partition, q._where, q._group, q._having, q._window, q._order, q._limit, q._for)
	return q
}

func (o *Query) TableBy(operater ksql.QueryInterface, as string) ksql.QueryInterface {
	o.table.sub = operater
	return o
}

func (o *Query) Sharding(sharding ksql.Sharding) {
	o.sharding = sharding
}

func (o *Query) GetSharding() ksql.Sharding {
	return o.sharding
}

func (o *Query) Table(table string) ksql.QueryInterface {
	o.table.table = table
	return o
}

func (o *Query) As(as string) ksql.QueryInterface {
	o.table.as = as
	return o
}

func (o *Query) Func(fun, column, as string) ksql.QueryInterface {
	o.columns.Append(&columnInfo{isFunc: true, column: column, fun: fun, as: as})
	return o
}

func (o *Query) Column(column, as string) ksql.QueryInterface {
	o.columns.Append(&columnInfo{column: column, as: as})
	return o
}

func (o *Query) IntoVar(vars ...string) ksql.QueryInterface {
	o.intoVars = append(o.intoVars, vars...)
	return o
}

func (o *Query) Columns(columns ...string) ksql.QueryInterface {
	for _, column := range columns {
		o.columns.Append(&columnInfo{column: column})
	}

	return o
}

func (o *Query) ColumnsExpress(columns ...ksql.ExpressInterface) ksql.QueryInterface {
	for _, column := range columns {
		o.columns.Append(&columnInfo{expr: column})
	}

	return o
}

func (o *Query) _data(data ...any) {
	o.binds = append(o.binds, data...)
}

func (o *Query) WhereIsNull(column string) ksql.QueryInterface {
	o.where.IsNull(column)
	return o
}

func (o *Query) WhereIsNotNull(column string) ksql.QueryInterface {
	o.where.IsNotNull(column)
	return o
}

func (o *Query) WhereIn(column string, val []any) ksql.QueryInterface {
	o.where.In(column, val)
	return o
}

func (o *Query) WhereNotIn(column string, val []any) ksql.QueryInterface {
	o.where.NotIn(column, val)
	return o
}

func (o *Query) WhereInBy(column string, sub ksql.QueryInterface) ksql.QueryInterface {
	o.where.InBy(column, sub)
	return o
}

func (o *Query) WhereNotInBy(column string, sub ksql.QueryInterface) ksql.QueryInterface {
	o.where.NotInBy(column, sub)
	return o
}

func (o *Query) Where(column string, op ksql.Op, val any) ksql.QueryInterface {
	o.where.Where(column, op, val)
	return o
}

func (o *Query) WhereExpress(expresses ...ksql.ExpressInterface) ksql.QueryInterface {
	for _, express := range expresses {
		o.where.Express(express)
	}

	return o
}

func (o *Query) OrWhere(call func(w ksql.WhereInterface)) ksql.QueryInterface {
	if call == nil {
		return o
	}
	o.where.OrWhere(call)
	return o
}

func (q *Query) AndWhere(call func(w ksql.WhereInterface)) ksql.QueryInterface {
	if call == nil {
		return q
	}

	q.where.AndWhere(call)
	return q
}

func (o *Query) Between(column string, begin, end any) ksql.QueryInterface {
	o.where.Between(column, begin, end)
	return o
}

func (o *Query) NotBetween(column string, begin, end any) ksql.QueryInterface {
	o.where.NotBetween(column, begin, end)
	return o
}

func (o *Query) Limit(limit int) ksql.QueryInterface {
	o.limitInfo.hasLimit = true
	o.limitInfo.limit = limit
	return o
}

func (o *Query) Offset(offset int) ksql.QueryInterface {
	o.limitInfo.hasOffset = true
	o.limitInfo.offset = offset
	return o
}

func (o *Query) Pagination(page, pageSize int) ksql.QueryInterface {
	o.Limit(pageSize)
	o.Offset((page - 1) * pageSize)
	return o
}

func (o *Query) Order(columns ...string) ksql.QueryInterface {
	for _, column := range columns {
		o.order.Append(&orderMeta{column: &columnInfo{column: column}, typ: "ASC"})
	}
	return o
}

func (o *Query) OrderDesc(columns ...string) ksql.QueryInterface {
	for _, column := range columns {
		o.order.Append(&orderMeta{column: &columnInfo{column: column}, typ: "DESC"})
	}
	return o
}

func (o *Query) Group(columns ...string) ksql.QueryInterface {
	for _, column := range columns {
		o.group.columns.Append(&columnInfo{column: column})
	}
	return o
}

func (o *Query) _join(table string) ksql.JoinInterface {
	join := NewJoin()
	join.Table(table)
	o.join = append(o.join, join)
	return join
}

func (o *Query) Join(table string) ksql.JoinInterface {
	return o._join(table).Inner()
}

func (o *Query) _joinExpress(express ksql.ExpressInterface) ksql.JoinInterface {
	join := NewJoin()
	join.Express(express)
	o.join = append(o.join, join)
	return join
}

func (o *Query) JoinExpress(express ksql.ExpressInterface) ksql.JoinInterface {
	return o._joinExpress(express)
}

func (o *Query) LeftJoin(table string) ksql.JoinInterface {
	return o._join(table).Left()
}

func (o *Query) RightJoin(table string) ksql.JoinInterface {
	return o._join(table).Right()
}

func (o *Query) HavingBetween(column string, begin, end any) ksql.QueryInterface {
	o.having.Between(column, begin, end)
	return o
}

func (o *Query) HavingNotBetween(column string, begin, end any) ksql.QueryInterface {
	o.having.NotBetween(column, begin, end)
	return o
}

func (o *Query) HavingIsNull(column string) ksql.QueryInterface {
	o.having.IsNull(column)
	return o
}

func (o *Query) HavingIsNotNull(column string) ksql.QueryInterface {
	o.having.IsNotNull(column)
	return o
}

func (o *Query) HavingIn(column string, val []any) ksql.QueryInterface {
	o.having.In(column, val)
	return o
}

func (o *Query) HavingNotIn(column string, val []any) ksql.QueryInterface {
	o.having.NotIn(column, val)
	return o
}

func (o *Query) HavingInBy(column string, sub ksql.QueryInterface) ksql.QueryInterface {
	o.having.InBy(column, sub)
	return o
}

func (o *Query) HavingNotInBy(column string, sub ksql.QueryInterface) ksql.QueryInterface {
	o.having.NotInBy(column, sub)
	return o
}

func (o *Query) Having(column string, op ksql.Op, val any) ksql.QueryInterface {
	o.having.Having(column, op, val)
	return o
}

func (o *Query) HavingExpress(expresses ...ksql.ExpressInterface) ksql.QueryInterface {
	for _, raw := range expresses {
		o.having.Express(raw)
	}
	return o
}

func (q *Query) AndHaving(call func(h ksql.HavingInterface)) ksql.QueryInterface {
	if call == nil {
		return q
	}
	q.having.AndHaving(call)
	return q
}

func (o *Query) OrHaving(call func(h ksql.HavingInterface)) ksql.QueryInterface {
	if call == nil {
		return o
	}
	o.having.OrHaving(call)
	return o
}

func (o *Query) ForUpdate() ksql.QueryInterface {
	o.forSql.Update()
	return o
}

func (o *Query) Distinct() ksql.QueryInterface {
	o.modifer = "DISTINCT"
	return o
}

func (o *Query) FuncDistinct(fun, column, as string) ksql.QueryInterface {
	o.columns.Append(&columnInfo{prefix: "DISTINCT", column: column, fun: fun, as: as})
	return o
}

func (o *Query) All() ksql.QueryInterface {
	o.modifer = "ALL"
	return o
}

func (o *Query) DistinctRow() ksql.QueryInterface {
	o.modifer = "DISTINCTROW"
	return o
}

func (o *Query) HighPriority() ksql.QueryInterface {
	o.highPriority = "HIGH_PRIORITY"
	return o
}

func (o *Query) StraightJoin() ksql.QueryInterface {
	o.straightJoin = "STRAIGHT_JOIN"
	return o
}

func (o *Query) SqlSmallResult() ksql.QueryInterface {
	o.sqlSmallResult = "SQL_SMALL_RESULT"
	return o
}

func (o *Query) SqlBigResult() ksql.QueryInterface {
	o.sqlBigResult = "SQL_BIG_RESULT"
	return o
}

func (o *Query) SqlBufferResult() ksql.QueryInterface {
	o.sqlBufferResult = "SQL_BUFFER_RESULT"
	return o
}

func (o *Query) SqlNoCache() ksql.QueryInterface {
	o.sqlNoCache = "SQL_NO_CACHE"
	return o
}

func (o *Query) SqlCalcFoundRows() ksql.QueryInterface {
	o.sqlCalcFoundRows = "SQL_CALC_FOUND_ROWS"
	return o
}

func (o *Query) Partitions(names ...string) ksql.QueryInterface {
	o.partitions = append(o.partitions, names...)
	return o
}

func (o *Query) GroupWithRollUp() ksql.QueryInterface {
	o.group.with = "WITH ROLLUP"
	return o
}

func (o *Query) Window(win, as string) ksql.QueryInterface {
	o.windows.Append(&window{name: win, as: as})
	return o
}

func (o *Query) OrderWithRollUp() ksql.QueryInterface {
	o.order.with = "WITH ROLLUP"
	return o
}

func (o *Query) For() ksql.ForInterface {
	return o.forSql
}

func (o *Query) WhereInCall(column string, call func(query ksql.QueryInterface)) ksql.QueryInterface {
	query := NewQuery()
	call(query)
	o.where.InBy(column, query)
	return o
}

func (o *Query) WhereNotInCall(column string, call func(query ksql.QueryInterface)) ksql.QueryInterface {
	query := NewQuery()
	call(query)
	o.where.NotInBy(column, query)
	return o
}
