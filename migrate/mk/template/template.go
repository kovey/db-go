package template

import (
	"go/format"
	"html/template"
	"strings"
)

const (
	tpl_migrate = `
package {{.Package}}

// ksql-tool
// migrator file: {{.Name}}
// tool version:  {{.ToolVersion}}
// created time:  {{.CreateTime}}

import (
	"context"
	"github.com/kovey/db-go/v3/db"
)

type {{.Name}} struct {
}

func (self *{{.Name}}) Up(ctx context.Context) error {
	// TODO Code
	return nil
}

func (self *{{.Name}}) Down(ctx context.Context) error {
	// TODO Code
	return nil
}

func (self *{{.Name}}) Id() uint64 {
	return {{.Id}}
}

func (self *{{.Name}}) Name() string {
	return "{{.Name}}"
}

func (self *{{.Name}}) Version() string {
	return "{{.Version}}"
}
`
)

type MigrateTemplate struct {
	Name        string
	Package     string
	Id          uint64
	Version     string
	ToolVersion string
	CreateTime  string
}

func (m *MigrateTemplate) Parse() ([]byte, error) {
	t := template.Must(template.New("migrate").Parse(tpl_migrate))
	builder := strings.Builder{}
	if err := t.Execute(&builder, m); err != nil {
		return nil, err
	}

	return format.Source([]byte(builder.String()))
}
