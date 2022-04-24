package parse

import (
	"go/ast"
	"go/token"
	"strings"
)

type structDecl struct {
	doc    string
	fields map[string]field
}

type field struct {
	methods map[string]string
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

func getMethods(f *ast.Field) map[string]string {
	m := make(map[string]string)
	if t, ok := f.Type.(*ast.InterfaceType); ok {
		for _, methods := range t.Methods.List {
			_, ok := m[methods.Names[0].Name]
			if ok {
				continue // methods already defined
			}
			funcType, ok := methods.Type.(*ast.FuncType)
			if !ok {
				continue // not a function
			}
			params := getProto(funcType.Params)
			results := getProto(funcType.Results)
			proto := "(" + strings.Join(params, ", ") + ")"
			if len(results) > 0 {
				proto += "(" + strings.Join(results, ", ") + ")"
			}
			m[methods.Names[0].Name] = proto
		}
	}
	return m
}

func getProto(funcParam *ast.FieldList) (params []string) {
	if funcParam == nil {
		return params
	}
	for _, param := range funcParam.List {
		proto := ""
		if param.Names != nil {
			proto += param.Names[0].Name
		}
		packageName, serviceName := getDepID(param.Type)
		proto += packageName + "." + serviceName
		params = append(params, proto)
	}
	return params
}
