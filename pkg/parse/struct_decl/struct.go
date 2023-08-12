package struct_decl

import (
	"go/ast"
	"go/token"
)

func GetStructDoc(f *ast.File, packageName string) map[string]string {
	structDoc := map[string]string{}
	for _, decl := range f.Decls {
		if d, ok := decl.(*ast.GenDecl); ok {
			name, doc := search(d)
			if name == "" && doc == "" {
				continue
			}
			structDoc[packageName+"."+name] = doc
		}
	}
	return structDoc
}

func search(decl *ast.GenDecl) (structName, doc string) {
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
