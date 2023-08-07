package parse

import (
	"go/ast"
	"go/token"
)

func searchStructDeclDoc(decl *ast.GenDecl) (string, string) {
	if decl.Tok != token.TYPE {
		return "", ""
	}
	spec := decl.Specs[0]
	ts, ok := spec.(*ast.TypeSpec)
	if !ok {
		return "", ""
	}

	doc := ""
	if decl.Doc != nil {
		doc = decl.Doc.List[0].Text
	}
	return ts.Name.Name, doc
}
