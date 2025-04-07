package table

import (
	"testing"

	ksql "github.com/kovey/db-go/v3"
	"github.com/stretchr/testify/assert"
)

func TestIndexUnique(t *testing.T) {
	i := &Index{Name: "idx_test", Type: ksql.Index_Type_Unique}
	i.Columns("user_id", "type")
	assert.Equal(t, "UNIQUE KEY `idx_test` (`user_id`,`type`)", i.Express())
	assert.Equal(t, "ADD UNIQUE INDEX `idx_test` (`user_id`,`type`)", i.AlterExpress())
	assert.Equal(t, ksql.Index_Type_Unique, i.Type)
}

func TestIndexNormal(t *testing.T) {
	i := &Index{Name: "idx_test", Type: ksql.Index_Type_Normal}
	i.Columns("user_id")
	assert.Equal(t, "KEY `idx_test` (`user_id`)", i.Express())
	assert.Equal(t, "ADD INDEX `idx_test` (`user_id`)", i.AlterExpress())
	assert.Equal(t, ksql.Index_Type_Normal, i.Type)
}

func TestIndexPrimary(t *testing.T) {
	i := &Index{Name: "idx_test", Type: ksql.Index_Type_Primary}
	i.Columns("user_id")
	assert.Equal(t, "PRIMARY KEY (`user_id`)", i.Express())
	assert.Equal(t, "ADD PRIMARY INDEX (`user_id`)", i.AlterExpress())
	assert.Equal(t, ksql.Index_Type_Primary, i.Type)
}

func TestIndexSpatial(t *testing.T) {
	i := &Index{Name: "idx_test", Type: ksql.Index_Type_Spatial}
	i.Columns("user_id")
	assert.Equal(t, "SPATIAL KEY (`user_id`)", i.Express())
	assert.Equal(t, "ADD SPATIAL INDEX (`user_id`)", i.AlterExpress())
	assert.Equal(t, ksql.Index_Type_Spatial, i.Type)
}

func TestIndexFullText(t *testing.T) {
	i := &Index{Name: "idx_test", Type: ksql.Index_Type_FullText}
	i.Columns("user_id")
	assert.Equal(t, "FULLTEXT KEY (`user_id`)", i.Express())
	assert.Equal(t, "ADD FULLTEXT INDEX (`user_id`)", i.AlterExpress())
	assert.Equal(t, ksql.Index_Type_FullText, i.Type)
}
