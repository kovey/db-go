package sql

import (
	"fmt"
	"strings"

	"github.com/kovey/db-go/v2/sql/meta"
)

const (
	selectFormat = "SELECT %s FROM %s AS %s %s %s %s %s %s %s %s"
	orderFormat  = "ORDER BY %s"
	groupFormat  = "GROUP BY %s"
	limitFormat  = "LIMIT %d,%d"
	joinFormat   = "%s %s AS %s ON %s"
	leftJoin     = "LEFT JOIN"
	rightJoin    = "RIGHT JOIN"
	innerJoin    = "INNER JOIN"
	subFormat    = "(%s)"
	and          = " AND "
	forUpdate    = "FOR UPDATE"
	where        = "WHERE "
	space        = " "
	comma        = ","
	emptyStr     = ""
	dot          = "."
	underline    = "_"
	star         = "*"
)

type Select struct {
	table     string
	alias     string
	columns   []string
	where     WhereInterface
	orWhere   WhereInterface
	limit     int
	offset    int
	orders    []string
	groups    []string
	having    *Having
	joins     []string
	forUpdate string
	sub       *Select
	joinArgs  []any
}

func NewSelectSub(sub *Select, alias string) *Select {
	return &Select{
		sub: sub, alias: alias, columns: make([]string, 0), where: nil, orWhere: nil, limit: 0, offset: 0,
		orders: make([]string, 0), groups: make([]string, 0), having: nil, joins: make([]string, 0), forUpdate: emptyStr,
	}
}

func NewSelect(table string, alias string) *Select {
	if alias == emptyStr {
		if strings.Contains(table, dot) {
			alias = strings.ReplaceAll(table, dot, underline)
		} else {
			alias = table
		}
	}

	return &Select{
		table: table, alias: alias, columns: make([]string, 0), where: nil, orWhere: nil, limit: 0, offset: 0,
		orders: make([]string, 0), groups: make([]string, 0), having: nil, joins: make([]string, 0), forUpdate: emptyStr,
	}
}

func (s *Select) Columns(columns ...string) *Select {
	for _, col := range columns {
		column := meta.NewColumn(col)
		column.SetTable(s.alias)
		s.columns = append(s.columns, column.String())
	}

	return s
}

func (s *Select) ColMeta(columns ...*meta.Column) *Select {
	for _, column := range columns {
		column.SetTable(s.alias)
		s.columns = append(s.columns, column.String())
	}

	return s
}

func (s *Select) SetColumns(columns []string) *Select {
	s.columns = make([]string, len(columns))
	for index, col := range columns {
		column := meta.NewColumn(col)
		column.SetTable(s.alias)
		s.columns[index] = column.String()
	}

	return s
}

func (s *Select) SetColMeta(columns []*meta.Column) *Select {
	s.columns = make([]string, len(columns))
	for index, column := range columns {
		column.SetTable(s.alias)
		s.columns[index] = column.String()
	}

	return s
}

func (s *Select) CaseWhen(caseWhens ...*meta.CaseWhen) *Select {
	for _, caseWhen := range caseWhens {
		s.columns = append(s.columns, caseWhen.String())
	}

	return s
}

func (s *Select) Where(where WhereInterface) *Select {
	s.where = where
	return s
}

func (s *Select) OrWhere(where WhereInterface) *Select {
	s.orWhere = where
	return s
}

func (s *Select) Having(having *Having) *Select {
	s.having = having
	return s
}

func (s *Select) Limit(size int) *Select {
	if size < 0 {
		return s
	}

	s.limit = size
	return s
}

func (s *Select) GetLimit() int {
	return s.limit
}

func (s *Select) Offset(offset int) *Select {
	if offset < 0 {
		return s
	}

	s.offset = offset
	return s
}

func (s *Select) GetOffset() int {
	return s.offset
}

func (s *Select) Order(orders ...string) *Select {
	for _, order := range orders {
		s.orders = append(s.orders, formatOrder(order))
	}

	return s
}

func (s *Select) Group(groups ...string) *Select {
	for _, group := range groups {
		s.groups = append(s.groups, formatValue(group))
	}

	return s
}

func (s *Select) Args() []any {
	args := make([]any, 0)
	if s.sub != nil {
		args = append(args, s.sub.Args()...)
	}

	if len(s.joinArgs) > 0 {
		args = append(args, s.joinArgs...)
	}

	if s.where != nil {
		args = append(args, s.where.Args()...)
	}

	if s.orWhere != nil {
		args = append(args, s.orWhere.Args()...)
	}

	if s.having != nil {
		args = append(args, s.having.Args()...)
	}

	return args
}

func (s *Select) join(jt string, join *Join) *Select {
	s.columns = append(s.columns, join.columns...)
	s.joins = append(s.joins, fmt.Sprintf(joinFormat, jt, join.tableName(), formatValue(join.alias), join.on))
	args := join.args()
	if len(args) > 0 {
		s.joinArgs = append(s.joinArgs, args...)
	}

	return s
}

func (s *Select) LeftJoinSub(sub *Select, alias, on string, columns ...string) *Select {
	return s.LeftJoinWith(NewJoinSub(sub, alias, on, columns...))
}

func (s *Select) LeftJoin(table, alias, on string, columns ...string) *Select {
	return s.LeftJoinWith(NewJoin(table, alias, on, columns...))
}

func (s *Select) LeftJoinWith(join *Join) *Select {
	return s.join(leftJoin, join)
}

func (s *Select) RightJoinSub(sub *Select, alias, on string, columns ...string) *Select {
	return s.RightJoinWith(NewJoinSub(sub, alias, on, columns...))
}

func (s *Select) RightJoin(table, alias, on string, columns ...string) *Select {
	return s.RightJoinWith(NewJoin(table, alias, on, columns...))
}

func (s *Select) RightJoinWith(join *Join) *Select {
	return s.join(rightJoin, join)
}

func (s *Select) InnerJoinSub(sub *Select, alias, on string, columns ...string) *Select {
	return s.InnerJoinWith(NewJoinSub(sub, alias, on, columns...))
}

func (s *Select) InnerJoin(table, alias, on string, columns ...string) *Select {
	return s.InnerJoinWith(NewJoin(table, alias, on, columns...))
}

func (s *Select) InnerJoinWith(join *Join) *Select {
	return s.join(innerJoin, join)
}

func (s *Select) getColumns() string {
	if len(s.columns) == 0 {
		return star
	}

	return strings.Join(s.columns, comma)
}

func (s *Select) getJoin() string {
	if len(s.joins) == 0 {
		return emptyStr
	}

	return strings.Join(s.joins, space)
}

func (s *Select) getWhere() string {
	wheres := make([]string, 0)
	if s.where != nil {
		wheres = append(wheres, s.where.Prepare())
	}

	if s.orWhere != nil {
		if len(wheres) == 0 {
			wheres = append(wheres, s.orWhere.OrPrepare())
		} else {
			wheres = append(wheres, strings.Replace(s.orWhere.OrPrepare(), where, emptyStr, 1))
		}
	}

	return strings.Join(wheres, and)
}

func (s *Select) getHaving() string {
	if s.having == nil {
		return emptyStr
	}

	return s.having.Prepare()
}

func (s *Select) getOrder() string {
	if len(s.orders) == 0 {
		return emptyStr
	}

	return fmt.Sprintf(orderFormat, strings.Join(s.orders, comma))
}

func (s *Select) getGroup() string {
	if len(s.groups) == 0 {
		return emptyStr
	}

	return fmt.Sprintf(groupFormat, strings.Join(s.groups, comma))
}

func (s *Select) getLimit() string {
	if s.limit == 0 {
		return emptyStr
	}

	return fmt.Sprintf(limitFormat, s.offset, s.limit)
}

func (s *Select) Prepare() string {
	if s.sub == nil {
		return strings.Trim(fmt.Sprintf(
			selectFormat, s.getColumns(), formatValue(s.table), formatValue(s.alias), s.getJoin(), s.getWhere(), s.getGroup(), s.getHaving(),
			s.getOrder(), s.getLimit(), s.forUpdate,
		), space)
	}

	return strings.Trim(fmt.Sprintf(
		selectFormat, s.getColumns(), fmt.Sprintf(subFormat, s.sub.Prepare()), formatValue(s.alias), s.getJoin(), s.getWhere(), s.getGroup(), s.getHaving(),
		s.getOrder(), s.getLimit(), s.forUpdate,
	), space)
}

func (s *Select) String() string {
	return String(s)
}

func (s *Select) WhereByMap(where meta.Where) *Select {
	if s.where == nil {
		s.where = NewWhere()
	}

	for field, value := range where {
		s.where.Eq(field, value)
	}

	return s
}

func (s *Select) WhereByList(where meta.List) *Select {
	if s.where == nil {
		s.where = NewWhere()
	}

	for _, value := range where {
		s.where.Statement(value)
	}

	return s
}

func (s *Select) ForUpdate() *Select {
	s.forUpdate = forUpdate
	return s
}
