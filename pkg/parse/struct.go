package parse

import (
	"go/ast"
	"go/token"
)

type structDecl struct {
	doc    string
	fields map[string]field
}

type field struct {
	methods []string
}

func searchStructDecl(decl *ast.GenDecl) (string, structDecl) {
	s := structDecl{}
	if decl.Tok != token.TYPE {
		return "", s
	}
	spec := decl.Specs[0]
	ts, ok := spec.(*ast.TypeSpec)
	if !ok {
		return "", s
	}
	st, ok := ts.Type.(*ast.StructType)
	if !ok {
		return "", s
	}

	fields := map[string]field{}
	for _, f := range st.Fields.List {
		fields[f.Names[0].Name] = field{
			methods: getMethods(f),
		}
	}
	s.fields = fields

	if decl.Doc != nil {
		s.doc = decl.Doc.List[0].Text
	}
	return ts.Name.Name, s
}

func getMethods(f *ast.Field) []string {
	var m []string
	switch t := f.Type.(type) {
	case *ast.InterfaceType:
		for _, methods := range t.Methods.List {
			m = append(m, methods.Names[0].Name)
		}
	}
	return m
}
