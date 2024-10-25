package template

import (
	"go/format"
	"html/template"
	"strings"
)

const (
	main_tpl = `
package main

import (
	"github.com/kovey/db-go/v3/migplug"
	{{range $index, $import := .Imports}}"{{$import}}"{{"\r\n"}}{{end}}
)

type Migrator struct {
}

func (m *Migrator) Register(c migplug.CoreInterface) {
	{{range $index, $migrate := .Migrates}}c.Add(&migrations.{{$migrate}}{}){{"\n"}}{{end}}
}

func Migrate() migplug.PluginInterface {
	return &Migrator{}
}
`
)

type MainTpl struct {
	Imports  []string
	Migrates []string
}

func (m *MainTpl) Parse() ([]byte, error) {
	t := template.Must(template.New("main_tpl").Parse(main_tpl))
	builder := strings.Builder{}
	if err := t.Execute(&builder, m); err != nil {
		return nil, err
	}

	return format.Source([]byte(builder.String()))
}

func (m *MainTpl) Has(name string) bool {
	for _, migrate := range m.Migrates {
		if migrate == name {
			return true
		}
	}

	return false
}
