package compile

import (
	"fmt"
	"go/format"
	"html/template"
	"strings"
)

const (
	import_context = "context"
	import_model   = "github.com/kovey/db-go/v3/model"
	template_korm  = `package compile
import(
	"context"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/model"
)

type {{.Name}} struct {
{{range .Columns}}{{.Name}} {{.Type}}{{"\r\n"}}{{end}}
}

func (self *{{.Name}}) Columns() []string {
	return []string{ {{.ReturnColumns | safe}} }
}

func (self *{{.Name}}) Values() []any {
	return []any{ {{.ReturnValues | safe}} }
}

func (self *{{.Name}}) Save(ctx context.Context) error {
	return self.Model.Save(ctx, self)
}

func (self *{{.Name}}) Delete(ctx context.Context) error {
	return self.Model.Delete(ctx, self)
}

func (self *{{.Name}}) Query() ksql.BuilderInterface[*{{.Name}}] {
	return model.Row(self)
}
`
)

type columnInfo struct {
	Name string
	Type string
}

type templateKorm struct {
	Name          string
	Columns       []*columnInfo
	ReturnValues  string
	ReturnColumns string
}

func (t *templateKorm) init(columns []*_column) {
	var cc []string
	var vals []string
	for _, column := range columns {
		t.Columns = append(t.Columns, &columnInfo{Name: column.name, Type: "int"})
		cc = append(cc, fmt.Sprintf(`"%s"`, column.tag))
		vals = append(vals, fmt.Sprintf("&self.%s", column.name))
	}

	t.ReturnColumns = strings.Join(cc, ",")
	t.ReturnValues = strings.Join(vals, ",")
}

func (tk *templateKorm) Parse() ([]byte, error) {
	t := template.Must(template.New("main_tpl").Funcs(template.FuncMap{"safe": func(tag string) template.HTML {
		return template.HTML(tag)
	}}).Parse(template_korm))
	builder := strings.Builder{}
	if err := t.Execute(&builder, tk); err != nil {
		return nil, err
	}

	return format.Source([]byte(builder.String()))
}
