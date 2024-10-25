package mysql

import "github.com/kovey/db-go/v3/migrate/schema"

type ColumnMetaChanged struct {
	o schema.ColumnInfoInterface
	n schema.ColumnInfoInterface
}

func (c *ColumnMetaChanged) Old() schema.ColumnInfoInterface {
	return c.o
}

func (c *ColumnMetaChanged) New() schema.ColumnInfoInterface {
	return c.n
}

type ColumnChanged struct {
	adds    []schema.ColumnInfoInterface
	changes []schema.ColumnMetaChangedInterface
	deletes []schema.ColumnInfoInterface
}

func (c *ColumnChanged) Adds() []schema.ColumnInfoInterface {
	return c.adds
}

func (c *ColumnChanged) Deletes() []schema.ColumnInfoInterface {
	return c.deletes
}

func (c *ColumnChanged) Changes() []schema.ColumnMetaChangedInterface {
	return c.changes
}

type IndexChanged struct {
	adds    []schema.IndexInfoInterface
	deletes []schema.IndexInfoInterface
}

func (i *IndexChanged) Adds() []schema.IndexInfoInterface {
	return i.adds
}

func (i *IndexChanged) Deletes() []schema.IndexInfoInterface {
	return i.deletes
}

type Changed struct {
	column *ColumnChanged
	index  *IndexChanged
}

func (c *Changed) Column() schema.ColumnChangedInterface {
	return c.column
}

func (c *Changed) Index() schema.IndexChangedInterface {
	return c.index
}
