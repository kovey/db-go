package table

import (
	"fmt"
	"strings"

	"github.com/kovey/db-go/v3"
)

type Index struct {
	Name    string
	Type    ksql.IndexType
	columns []string
}

func (i *Index) Columns(columns ...string) {
	for _, column := range columns {
		i.columns = append(i.columns, fmt.Sprintf("`%s`", column))
	}
}

func (i *Index) Express() string {
	switch i.Type {
	case ksql.Index_Type_Unique:
		return fmt.Sprintf("UNIQUE KEY `%s` (%s)", i.Name, strings.Join(i.columns, ","))
	case ksql.Index_Type_Primary:
		return fmt.Sprintf("PRIMARY KEY (%s)", i.columns[0])
	case ksql.Index_Type_FullText:
		return fmt.Sprintf("FULLTEXT KEY (%s)", i.columns[0])
	case ksql.Index_Type_Spatial:
		return fmt.Sprintf("SPATIAL KEY (%s)", i.columns[0])
	default:
		return fmt.Sprintf("KEY `%s` (%s)", i.Name, strings.Join(i.columns, ","))
	}
}

func (i *Index) AlterExpress() string {
	switch i.Type {
	case ksql.Index_Type_Unique:
		return fmt.Sprintf("ADD UNIQUE INDEX %s (%s)", i.Name, strings.Join(i.columns, ","))
	case ksql.Index_Type_Primary:
		return fmt.Sprintf("ADD PRIMARY INDEX (%s)", i.columns[0])
	case ksql.Index_Type_FullText:
		return fmt.Sprintf("ADD FULLTEXT INDEX (%s)", i.columns[0])
	case ksql.Index_Type_Spatial:
		return fmt.Sprintf("ADD SPATIAL INDEX (%s)", i.columns[0])
	default:
		return fmt.Sprintf("ADD INDEX %s (%s)", i.Name, strings.Join(i.columns, ","))
	}
}
