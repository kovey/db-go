package rows

import "database/sql"

type Rows[T any] struct {
	rows []T
}

func NewRows[T any]() *Rows[T] {
	return &Rows[T]{rows: make([]T, 0)}
}

func (r *Rows[T]) All() []T {
	return r.rows
}

func (r *Rows[T]) Scan(rows *sql.Rows, model T) error {
	cols, err := rows.Columns()
	if err != nil {
		return err
	}
	for rows.Next() {
		row := NewRow(model)
		if err := row.ScanByRows(rows, cols); err != nil {
			return err
		}

		r.rows = append(r.rows, row.Model)
	}

	return nil
}
