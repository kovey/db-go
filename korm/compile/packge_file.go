package compile

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
)

type _packageFile struct {
	pkg   *types.Package
	files map[string]*ast.File
}

func newPackageFile(path, name string) *_packageFile {
	return &_packageFile{pkg: types.NewPackage(path, name), files: make(map[string]*ast.File)}
}

func (p *_packageFile) ParseFile(fset *token.FileSet, files []string) error {
	for _, file := range files {
		f, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			return err
		}
		p.files[file] = f
	}

	return nil
}
