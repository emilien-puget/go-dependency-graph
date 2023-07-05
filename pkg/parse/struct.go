package parse

import (
	"go/ast"
	"go/token"
)

type StructDecl struct {
	doc     string
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

func searchStructDecl(decl *ast.GenDecl) (string, StructDecl) {
	s := StructDecl{}
	if decl.Tok != token.TYPE {
		return "", s
	}
	spec := decl.Specs[0]
	ts, ok := spec.(*ast.TypeSpec)
	if !ok {
		return "", s
	}

	if decl.Doc != nil {
		s.doc = decl.Doc.List[0].Text
	}
	return ts.Name.Name, s
}
