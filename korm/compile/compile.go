package compile

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"
)

const (
	korm_gen_flag      = "//go:korm "
	korm_gen_flag_only = "//go:korm"
	method_query       = "Query"
	method_save        = "Save"
	method_delete      = "Delete"
	method_columns     = "Columns"
	method_values      = "Values"
	method_clone       = "Clone"
	struct_model       = "Model"
)

var support_methods = map[string]byte{
	method_query:   1,
	method_save:    1,
	method_delete:  1,
	method_columns: 1,
	method_values:  1,
	method_clone:   1,
}

func isSupport(fnName string) bool {
	_, ok := support_methods[fnName]
	return ok
}

func getMethod(method string) string {
	method = strings.ToLower(method)
	for m := range support_methods {
		if method == strings.ToLower(m) {
			return m
		}
	}

	return ""
}

func isKormStruct(gn *ast.GenDecl, flag *string) bool {
	for i := len(gn.Doc.List) - 1; i >= 0; i-- {
		if strings.HasPrefix(strings.Trim(gn.Doc.List[i].Text, " "), korm_gen_flag) {
			*flag = strings.ReplaceAll(gn.Doc.List[i].Text, korm_gen_flag, "")
			return true
		}

		if strings.Trim(gn.Doc.List[i].Text, " ") == korm_gen_flag_only {
			return true
		}
	}

	return false
}

func _parseStruct(pkg *_packageFile) *_structInfos {
	sts := newStructInfos()
	for file, f := range pkg.files {
		for _, d := range f.Decls {
			if gn, ok := d.(*ast.GenDecl); ok && gn.Doc != nil {
				var flag = ""
				if !isKormStruct(gn, &flag) {
					continue
				}

				var _struct *_structInfo
				for _, s := range gn.Specs {
					if ss, ok := s.(*ast.TypeSpec); ok {
						if tmp, ok := sts.structs[ss.Name.Name]; ok {
							_struct = tmp
							_struct.ts = ss
						} else {
							_struct = &_structInfo{name: ss.Name.Name, file: f, filePath: file, tpl: &templateKorm{Name: ss.Name.Name}, ts: ss, fns: make(map[string]*_funcInfo)}
							sts.structs[_struct.name] = _struct
						}
						break
					}
				}

				if flag != "" {
					_struct.hasGen = true
					_struct.flag = flag
					sts.files[file] = f
				}
				sts.structs[_struct.name] = _struct
			}

			if fn, ok := d.(*ast.FuncDecl); ok && fn.Recv != nil {
				if len(fn.Recv.List) == 0 || !isSupport(fn.Name.Name) {
					continue
				}

				star, ok := fn.Recv.List[0].Type.(*ast.StarExpr)
				if !ok || star.X == nil {
					continue
				}

				indent, ok := star.X.(*ast.Ident)
				if !ok {
					continue
				}

				has := false
				if fn.Doc != nil {
					for i := len(fn.Doc.List) - 1; i >= 0; i-- {
						if strings.Trim(fn.Doc.List[i].Text, " ") == korm_gen_flag_only {
							has = true
							break
						}
					}
				}

				var _struct = sts.structs[indent.Name]
				if _struct == nil {
					_struct = &_structInfo{name: indent.Name, filePath: file, file: f, tpl: &templateKorm{Name: indent.Name}, fnHasGen: has, fns: make(map[string]*_funcInfo)}
					sts.structs[indent.Name] = _struct
				}
				ffn := &_funcInfo{fn: fn, hasGen: has, structName: "self"}
				if len(fn.Recv.List[0].Names) > 0 {
					ffn.recvName = fn.Recv.List[0].Names[0].Name
				}
				_struct.fns[ffn.fn.Name.Name] = ffn
				if has {
					_struct.fnHasGen = true
					if _, ok := sts.files[file]; !ok {
						sts.files[file] = f
					}
				}
			}
		}
	}

	return sts
}

func Compile(projectDir, tempDir string, args []string) error {
	packageInfo, err := parsePackageInfo(projectDir, "")
	if err != nil {
		return err
	}
	if packageInfo.Module.Path == "" {
		return fmt.Errorf("module in %s not found", projectDir)
	}
	files := make([]string, 0, len(args))
	projectName := packageInfo.Module.Path
	packageName := packageInfo.Name
	for i, arg := range args {
		if arg == "-p" && i+1 < len(args) {
			packageName = args[i+1]
		}
		if strings.HasPrefix(arg, "-") {
			continue
		}
		if strings.HasPrefix(arg, projectDir+string(filepath.Separator)) && strings.HasSuffix(arg, ".go") {
			files = args[i:]
			break
		}
	}

	if (packageName != "main" && !strings.HasPrefix(packageName, projectName)) || len(files) == 0 {
		return nil
	}

	fset := token.NewFileSet()
	pkg := newPackageFile(projectDir, packageInfo.Name)
	if err := pkg.ParseFile(fset, files); err != nil {
		return err
	}

	sts := _parseStruct(pkg)
	return sts.replace(fset, tempDir, args)
}
