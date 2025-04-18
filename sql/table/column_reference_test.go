package table

import (
	"strings"
	"testing"

	ksql "github.com/kovey/db-go/v3"
	"github.com/stretchr/testify/assert"
)

func TestColumnReference(t *testing.T) {
	r := NewColumnReference("user").Column("user_name", 20, ksql.Order_Desc).Column("id", 0, ksql.Order_None).Match(ksql.Reference_Match_Full)
	var builder strings.Builder
	r.Build(&builder)
	assert.Equal(t, " REFERENCES `user` (`user_name`(20) DESC, `id`) MATCH FULL", builder.String())
}

func TestColumnReferenceOn(t *testing.T) {
	r := NewColumnReference("user").Express("round(aa, 2)", ksql.Order_Asc).Match(ksql.Reference_Match_Simple).On(ksql.Reference_On_Opt_UPDATE, ksql.Reference_Option_No_Action)
	var builder strings.Builder
	r.Build(&builder)
	assert.Equal(t, " REFERENCES `user` ((round(aa, 2)) ASC) MATCH SIMPLE ON UPDATE NO ACTION", builder.String())
}
