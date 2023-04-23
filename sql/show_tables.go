package sql

const (
	showTablesFormat = "SHOW TABLES"
)

type ShowTables struct {
}

func NewShowTables() *ShowTables {
	return &ShowTables{}
}

func (d *ShowTables) Args() []any {
	return []any{}
}

func (d *ShowTables) Prepare() string {
	return showTablesFormat
}

func (d *ShowTables) String() string {
	return String(d)
}
