package table

import (
	"strings"
	"testing"

	ksql "github.com/kovey/db-go/v3"
	"github.com/stretchr/testify/assert"
)

func TestIndexColumn(t *testing.T) {
	i := &IndexColumns{}
	i.Append(&IndexColumn{Type: Index_Column_Type_Expr, Name: "sum(amount)", Order: ksql.Order_None})
	i.Append(&IndexColumn{Type: Index_Column_Type_Name, Name: "id", Order: ksql.Order_Asc})
	i.Append(&IndexColumn{Type: Index_Column_Type_Pure_Name, Name: "name", Order: ksql.Order_Desc})
	i.Append(&IndexColumn{Type: Index_Column_Type_Name, Name: "nick", Order: ksql.Order_Asc, Length: 10})
	var builder strings.Builder
	i.Build(&builder)
	assert.Equal(t, " ((sum(amount)), `id` ASC, `name` DESC, `nick`(10) ASC)", builder.String())
}
