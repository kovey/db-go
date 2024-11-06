package orm

import (
	"go/format"
	"html/template"
	"strings"
)

const (
	tpl_model = `
package {{.Package}}
// Code generated by ksql-tool. 
// Do'nt Edit!!!
// Do'nt Edit!!!
// Do'nt Edit!!!
// {{.Comment}}
// from database: {{.DbName}}
// table:         {{.Table}}
// orm version:   {{.Version}}
// created time:  {{.CreateTime}}

import(
	"context"

	"github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/model" {{range .Imports}}{{"\n"}}"{{.}}"{{end}}
)

const (
	Table_{{.Name}} = "{{.Table}}" // {{.Comment}}{{range .Consts}}{{"\n"}}Table_{{.Table}}_{{.Column}} = "{{.Name}}" // {{.Comment}}{{end}}
)

type {{.Name}} struct {
	*model.Model {{.ModelTag | safe}} // model {{range .Fields}}{{"\n"}}{{.Name}} {{if .CanNull}}*{{end}}{{.Type}} {{.Tag | safe}} // {{.Comment}}{{end}}
}

func New{{.Name}}() *{{.Name}} {
	return &{{.Name}}{Model: model.NewModel(Table_{{.Name}}, "{{.PrimaryId}}", model.{{.PrimaryType}})}
}

func (self *{{.Name}}) Save(ctx context.Context) error {
	return self.Model.Save(ctx, self)
}

func (self *{{.Name}}) Clone() ksql.RowInterface {
	return New{{.Name}}()
}

func (self *{{.Name}}) Values() []any {
	return []any{ {{.Values | safe}} }
}

func (self *{{.Name}}) Columns() []string {
	return []string{ {{.Columns | safe}} }
}

func (self *{{.Name}}) Delete(ctx context.Context) error {
	return self.Model.Delete(ctx, self)
}
`
)

type field struct {
	Name    string
	Type    string
	Comment string
	Tag     string
	CanNull bool
}

type constInfo struct {
	Table   string
	Column  string
	Name    string
	Comment string
}

type modelTpl struct {
	Imports     []string
	Package     string
	Name        string
	Fields      []field
	Table       string
	PrimaryId   string
	PrimaryType string
	Values      string
	Columns     string
	Comment     string
	Version     string
	CreateTime  string
	ModelTag    string
	HasSql      bool
	Consts      []constInfo
	DbName      string
}

func (m *modelTpl) Parse() ([]byte, error) {
	t := template.Must(template.New("main_tpl").Funcs(template.FuncMap{"safe": func(tag string) template.HTML {
		return template.HTML(tag)
	}}).Parse(tpl_model))
	builder := strings.Builder{}
	if err := t.Execute(&builder, m); err != nil {
		return nil, err
	}

	return format.Source([]byte(builder.String()))
}
