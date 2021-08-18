package sql

import (
	"fmt"
	"strings"
)

const (
	selectFormat string = "SELECT %s FROM %s AS %s %s %s %s %s %s %s"
	columnFormat string = "`%s`.`%s`"
	orderFormat  string = "ORDER BY %s"
	groupFormat  string = "GROUP BY %s"
	limitFormat  string = "LIMIT %d,%d"
	joinFormat   string = "%s %s AS %s ON %s"
	leftJoin     string = "LEFT JOIN"
	rightJoin    string = "RIGHT JOIN"
	innerJoin    string = "INNER JOIN"
)

type Select struct {
	table   string
	alias   string
	columns []string
	where   *Where
	orWhere *Where
	limit   int
	offset  int
	orders  []string
	groups  []string
	having  *Having
	joins   []string
}

func NewSelect(table string, alias string) *Select {
	if alias == "" {
		alias = table
	}

	return &Select{
		table: table, alias: alias, columns: make([]string, 0), where: nil, orWhere: nil, limit: 0, offset: 0,
		orders: make([]string, 0), groups: make([]string, 0), having: nil, joins: make([]string, 0),
	}
}

func (s *Select) Columns(columns ...string) *Select {
	for _, column := range columns {
		s.columns = append(s.columns, fmt.Sprintf(columnFormat, s.alias, column))
	}

	return s
}

func (s *Select) Where(where *Where) *Select {
	s.where = where
	return s
}

func (s *Select) OrWhere(where *Where) *Select {
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

func (s *Select) Offset(offset int) *Select {
	if offset < 0 {
		return s
	}

	s.offset = offset
	return s
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

func (s *Select) Args() []interface{} {
	args := make([]interface{}, 0)
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

func (s *Select) join(jt string, table string, alias string, on string, columns ...string) *Select {
	for _, column := range columns {
		s.columns = append(s.columns, fmt.Sprintf(columnFormat, alias, column))
	}

	s.joins = append(s.joins, fmt.Sprintf(joinFormat, jt, formatValue(table), formatValue(alias), on))
	return s
}

func (s *Select) LeftJoin(table string, alias string, on string, columns ...string) *Select {
	return s.join(leftJoin, table, alias, on, columns...)
}

func (s *Select) RightJoin(table string, alias string, on string, columns ...string) *Select {
	return s.join(rightJoin, table, alias, on, columns...)
}

func (s *Select) InnerJoin(table string, alias string, on string, columns ...string) *Select {
	return s.join(innerJoin, table, alias, on, columns...)
}

func (s *Select) getColumns() string {
	if len(s.columns) == 0 {
		return "*"
	}

	return strings.Join(s.columns, ",")
}

func (s *Select) getJoin() string {
	if len(s.joins) == 0 {
		return ""
	}

	return strings.Join(s.joins, " ")
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
			wheres = append(wheres, strings.Replace(s.orWhere.OrPrepare(), "WHERE ", "", 1))
		}
	}

	return strings.Join(wheres, " AND ")
}

func (s *Select) getHaving() string {
	if s.having == nil {
		return ""
	}

	return s.having.Prepare()
}

func (s *Select) getOrder() string {
	if len(s.orders) == 0 {
		return ""
	}

	return fmt.Sprintf(orderFormat, strings.Join(s.orders, ","))
}

func (s *Select) getGroup() string {
	if len(s.groups) == 0 {
		return ""
	}

	return fmt.Sprintf(groupFormat, strings.Join(s.groups, ","))
}

func (s *Select) getLimit() string {
	if s.limit == 0 {
		return ""
	}

	return fmt.Sprintf(limitFormat, s.offset, s.limit)
}

func (s *Select) Prepare() string {
	return fmt.Sprintf(
		selectFormat, s.getColumns(), formatValue(s.table), formatValue(s.alias), s.getJoin(), s.getWhere(), s.getGroup(), s.getHaving(), s.getOrder(), s.getLimit(),
	)
}

func (s *Select) String() string {
	return String(s)
}

func (s *Select) WhereByMap(where map[string]interface{}) *Select {
	if s.where == nil {
		s.where = NewWhere()
	}

	for field, value := range where {
		s.where.Eq(field, value)
	}

	return s
}

func (s *Select) WhereByList(where []string) *Select {
	if s.where == nil {
		s.where = NewWhere()
	}

	for _, value := range where {
		s.where.Statement(value)
	}

	return s
}
