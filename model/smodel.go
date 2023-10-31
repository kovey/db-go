package model

import (
	"context"

	"github.com/kovey/db-go/v2/itf"
	"github.com/kovey/db-go/v2/sql/meta"
	"github.com/kovey/db-go/v2/table"
)

type ModelShardingInterface interface {
	itf.ModelInterface
	Empty() bool
}

type BaseSharding[T ModelShardingInterface] struct {
	Table     table.TableShardingInterface[T]
	primaryId *PrimaryId
	isInsert  bool
	isEmpty   bool
}

func NewBaseSharding[T ModelShardingInterface](tb table.TableShardingInterface[T], primaryId *PrimaryId) *BaseSharding[T] {
	return &BaseSharding[T]{Table: tb, primaryId: primaryId, isInsert: true, isEmpty: false}
}

func (b *BaseSharding[T]) NoAutoInc() {
	b.primaryId.IsAutoInc = false
}

func (b *BaseSharding[T]) Save(key any, model T) error {
	return b.SaveCtx(context.Background(), key, model)
}

func (b *BaseSharding[T]) SaveCtx(ctx context.Context, key any, model T) error {
	columns := model.Columns()
	fields := model.Fields()
	values := model.Values()
	var primary any
	data := meta.NewData()
	for index, column := range columns {
		if column == b.primaryId.Name {
			primary = fields[index]
			b.primaryId.Parse(values[index])
			if b.primaryId.Null() {
				continue
			}

			continue
		}

		data[column] = values[index]
	}

	if !b.isInsert {
		where := meta.NewWhere()
		where.Add(b.primaryId.Name, b.primaryId.Value())
		_, err := b.Table.UpdateCtx(ctx, key, data, where)
		return err
	}

	id, err := b.Table.InsertCtx(ctx, key, data)
	if err != nil {
		return err
	}

	if id <= 0 {
		return nil
	}

	switch tmp := primary.(type) {
	case *int:
		*tmp = int(id)
	case *int8:
		*tmp = int8(id)
	case *int16:
		*tmp = int16(id)
	case *int32:
		*tmp = int32(id)
	case *int64:
		*tmp = int64(id)
	case *uint:
		*tmp = uint(id)
	case *uint8:
		*tmp = uint8(id)
	case *uint16:
		*tmp = uint16(id)
	case *uint32:
		*tmp = uint32(id)
	case *uint64:
		*tmp = uint64(id)
	}
	return nil
}

func (b *BaseSharding[T]) Delete(key any, model T) error {
	return b.DeleteCtx(context.Background(), key, model)
}

func (b *BaseSharding[T]) DeleteCtx(ctx context.Context, key any, model T) error {
	where := meta.NewWhere()
	columns := model.Columns()
	values := model.Values()
	for index, column := range columns {
		if column == b.primaryId.Name {
			where.Add(b.primaryId.Name, values[index])
			break
		}
	}

	_, err := b.Table.DeleteCtx(ctx, key, where)
	return err
}

func (b *BaseSharding[T]) FetchRow(key any, where meta.Where, model T) error {
	return b.FetchRowCtx(context.Background(), key, where, model)
}

func (b *BaseSharding[T]) FetchRowCtx(ctx context.Context, key any, where meta.Where, model T) error {
	b.isEmpty = false
	err := b.Table.FetchRowCtx(ctx, key, where, model)
	if err != nil {
		return err
	}

	b.isInsert = false
	return nil
}

func (b *BaseSharding[T]) LockRow(key any, where meta.Where, model T) error {
	return b.LockRowCtx(context.Background(), key, where, model)
}

func (b *BaseSharding[T]) LockRowCtx(ctx context.Context, key any, where meta.Where, model T) error {
	b.isEmpty = false
	err := b.Table.LockRowCtx(ctx, key, where, model)
	if err != nil {
		return err
	}

	b.isInsert = false
	return nil
}

func (b *BaseSharding[T]) Empty() bool {
	return b.isEmpty
}

func (b *BaseSharding[T]) SetEmpty() {
	b.isEmpty = true
}
