package compile

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type _column struct {
	name string
	tag  string
}

func getTag(value string) string {
	info := strings.Split(strings.Trim(value, "`"), " ")
	for _, tag := range info {
		dbTag := strings.Split(tag, ":")
		if len(dbTag) != 2 || dbTag[0] != "db" {
			continue
		}

		return strings.ReplaceAll(dbTag[1], `"`, "")
	}

	return ""
}

func getColumns(stInfo *_structInfo) []*_column {
	var columns []*_column
	st, ok := stInfo.ts.Type.(*ast.StructType)
	if !ok || st.Fields == nil {
		return columns
	}

	for _, field := range st.Fields.List {
		if len(field.Names) == 0 || field.Tag == nil {
			if star, ok := field.Type.(*ast.StarExpr); ok {
				if sel, ok := star.X.(*ast.SelectorExpr); ok {
					stInfo.hasModel = sel.Sel.Name == struct_model
				}
			}
			continue
		}

		columns = append(columns, &_column{name: field.Names[0].Name, tag: getTag(field.Tag.Value)})
	}
	return columns
}

type _funcInfo struct {
	fn         *ast.FuncDecl
	hasGen     bool
	structName string
	recvName   string
}

func (f *_funcInfo) replace(fn *ast.FuncDecl) {
	f.fn.Body.List = fn.Body.List
}

type _structInfo struct {
	ts         *ast.TypeSpec
	name       string
	hasGen     bool
	flag       string
	fnHasGen   bool
	fns        map[string]*_funcInfo
	file       *ast.File
	filePath   string
	hasChanged bool
	tpl        *templateKorm
	tplFuncs   map[string]*ast.FuncDecl
	imports    []*ast.ImportSpec
	hasContext bool
	hasModel   bool
	hasKsql    bool
}

func (s *_structInfo) initColumns(columns []*_column) error {
	if s.tplFuncs == nil {
		s.tplFuncs = make(map[string]*ast.FuncDecl)
	}
	s.tpl.init(columns)
	buf, err := s.tpl.Parse()
	if err != nil {
		return err
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", buf, 0)
	if err != nil {
		return err
	}

	for _, d := range file.Decls {
		if fn, ok := d.(*ast.FuncDecl); ok && fn.Recv != nil {
			s.tplFuncs[fn.Name.Name] = fn
		}
	}

	s.imports = file.Imports
	return nil
}

func (s *_structInfo) add(fn *ast.FuncDecl) {
	s.fns[fn.Name.Name] = &_funcInfo{structName: s.name, recvName: "self", fn: fn, hasGen: true}
	s.hasChanged = true
	s.file.Decls = append(s.file.Decls, fn)
	switch fn.Name.Name {
	case method_delete, method_save:
		s.hasContext = true
	case method_clone, method_query:
		s.hasKsql = true
	}
}

func (s *_structInfo) fill() {
	if s.flag == "*" {
		for method := range support_methods {
			if _, ok := s.fns[method]; ok {
				continue
			}

			s.add(s.tplFuncs[method])
		}
		return
	}

	for _, flag := range strings.Split(s.flag, ",") {
		method := getMethod(flag)
		if method == "" {
			continue
		}

		if _, ok := s.fns[method]; ok {
			continue
		}

		s.add(s.tplFuncs[method])
	}
}

func (s *_structInfo) replace() error {
	if !s.hasGen && !s.fnHasGen || s.ts == nil {
		return nil
	}

	if err := s.initColumns(getColumns(s)); err != nil {
		return err
	}

	for _, fn := range s.fns {
		if !fn.hasGen {
			continue
		}

		fn.replace(s.tplFuncs[fn.fn.Name.Name])
		s.hasChanged = true
	}

	if !s.hasGen {
		return nil
	}

	s.fill()
	if s.hasChanged {
		ims := s.file.Decls[0].(*ast.GenDecl)
		for _, im := range s.imports {
			has := false
			for _, imm := range s.file.Imports {
				if im.Path.Value == imm.Path.Value {
					has = true
					break
				}
			}
			if has {
				continue
			}

			if s.hasContext && im.Path.Value == import_context {
				ims.Specs = append(ims.Specs, im)
			}

			if s.hasModel && im.Path.Value == import_model {
				ims.Specs = append(ims.Specs, im)
			}

			if s.hasKsql && im.Path.Value == import_ksql {
				ims.Specs = append(ims.Specs, im)
			}
		}
	}
	return nil
}

type _structInfos struct {
	structs    map[string]*_structInfo
	files      map[string]*ast.File
	printerCfg *printer.Config
}

func newStructInfos() *_structInfos {
	return &_structInfos{structs: make(map[string]*_structInfo), files: make(map[string]*ast.File), printerCfg: &printer.Config{Tabwidth: 8, Mode: printer.SourcePos}}
}

func (s *_structInfos) hasChanged(filePath string) bool {
	for _, st := range s.structs {
		if st.filePath == filePath {
			return st.hasChanged
		}
	}

	return false
}

func (s *_structInfos) replace(fset *token.FileSet, tempDir string, args []string) error {
	for _, st := range s.structs {
		if err := st.replace(); err != nil {
			return err
		}
	}

	for file, f := range s.files {
		if !s.hasChanged(file) {
			continue
		}
		originPath := file
		tgDir := path.Join(tempDir, os.Getenv("TOOLEXEC_IMPORTPATH"))
		buffer := bytes.NewBuffer(nil)
		if err := s.printerCfg.Fprint(buffer, fset, f); err != nil {
			return err
		}

		_ = os.MkdirAll(tgDir, 0777)
		tmpEntryFile := path.Join(tgDir, filepath.Base(originPath))
		if err := os.WriteFile(tmpEntryFile, buffer.Bytes(), 0777); err != nil {
			return err
		}

		for i := range args {
			if args[i] == originPath {
				args[i] = tmpEntryFile
				break
			}
		}
	}

	return nil
}
