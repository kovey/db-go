package sql

import (
	"testing"

	ksql "github.com/kovey/db-go/v3"
	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	i := NewIndex().Index("idx_user_id_name")
	i.Type(ksql.Index_Type_FullText).On("user").Column("user_id", 0, ksql.Order_Asc).Column("name", 20, ksql.Order_Desc).Express("user_id - 1", ksql.Order_None).BlockSize("1M")
	assert.Equal(t, "CREATE FULLTEXT INDEX `idx_user_id_name` ON `user` (`user_id` ASC, `name`(20) DESC, (user_id - 1)) KEY_BLOCK_SIZE = 1M", i.Prepare())
	assert.Nil(t, i.Binds())
}

func TestIndexVisible(t *testing.T) {
	i := NewIndex().Index("idx_user_id_name")
	i.Algorithm(ksql.Index_Alg_BTree).On("user").Column("user_id", 0, ksql.Order_Asc).Column("name", 20, ksql.Order_Desc).Express("user_id - 1", ksql.Order_None)
	i.Comment("user index").WithParser("parsers").Visible().EngineAttribute("engine attr").SecondaryEngineAttribute("sec attr").AlgorithmOption(ksql.Index_Alg_Option_Copy).LockOption(ksql.Index_Lock_Option_Default)
	assert.Equal(t, "CREATE INDEX `idx_user_id_name` USING BTREE ON `user` (`user_id` ASC, `name`(20) DESC, (user_id - 1)) WITH PARSER parsers COMMENT 'user index' VISIBLE ENGINE_ATTRIBUTE = 'engine attr' SECONDARY_ENGINE_ATTRIBUTE = 'sec attr' ALGORITHM = COPY LOCK = DEFAULT", i.Prepare())
	assert.Nil(t, i.Binds())
}

func TestIndexOtherInvisible(t *testing.T) {
	i := NewIndex().Index("idx_user_id_name")
	i.Algorithm(ksql.Index_Alg_BTree).On("user").Column("user_id", 0, ksql.Order_Asc).Column("name", 20, ksql.Order_Desc).Express("user_id - 1", ksql.Order_None)
	i.Type(ksql.Index_Type_Normal).Comment("user index").WithParser("parsers").Invisible().EngineAttribute("engine attr").SecondaryEngineAttribute("sec attr").AlgorithmOption(ksql.Index_Alg_Option_Copy).LockOption(ksql.Index_Lock_Option_Default)
	assert.Equal(t, "CREATE INDEX `idx_user_id_name` USING BTREE ON `user` (`user_id` ASC, `name`(20) DESC, (user_id - 1)) WITH PARSER parsers COMMENT 'user index' INVISIBLE ENGINE_ATTRIBUTE = 'engine attr' SECONDARY_ENGINE_ATTRIBUTE = 'sec attr' ALGORITHM = COPY LOCK = DEFAULT", i.Prepare())
	assert.Nil(t, i.Binds())
}
