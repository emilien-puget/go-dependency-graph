package parse

import (
	"go/ast"
	"go/token"
)

type structDecl struct {
	fields  map[string]field
	methods []string
}

const (
	fieldKindInterface = "interface"
	fieldKindFunc      = "func"
)

type field struct {
	kind    string
	methods []string
	fn      string
}

func searchStructDecl(decl *ast.GenDecl) (structName, doc string) {
	if decl.Tok != token.TYPE {
		return "", ""
	}
	spec := decl.Specs[0]
	ts, ok := spec.(*ast.TypeSpec)
	if !ok {
		return "", ""
	}

	if decl.Doc != nil {
		doc = decl.Doc.List[0].Text
	}
	return ts.Name.Name, doc
}
