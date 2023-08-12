package struct_decl

import (
	"go/ast"
	"go/token"
)

func Search(decl *ast.GenDecl) (structName, doc string) {
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
