package parse

import (
	"fmt"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

type StructDecl struct {
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

type InterDecl struct {
	methods []string
	implems map[string]map[string]*StructDecl
}

// ExtractTypes extracts types information from a package.
func ExtractTypes(pkg *packages.Package) (map[string]*StructDecl, map[string]*InterDecl) {
	classes := make(map[string]*StructDecl)
	inter := make(map[string]*InterDecl)

	// Iterate through all types in the package.
	for _, typ := range pkg.TypesInfo.Defs {
		readTypeObject(typ, classes, inter)
	}

	return classes, inter
}

func readTypeObject(typ types.Object, classes map[string]*StructDecl, inters map[string]*InterDecl) {
	if typ == nil {
		return
	}

	tp, ok := typ.Type().(*types.Named)
	if !ok {
		return
	}

	switch s := tp.Underlying().(type) {
	case *types.Interface:
		inter := &InterDecl{
			implems: make(map[string]map[string]*StructDecl),
		}
		for i := 0; i < s.NumMethods(); i++ {
			inter.methods = append(inter.methods, getFuncAsString(s.Method(i)))
		}
		inters[tp.Obj().Name()] = inter
	case *types.Struct:
		class := &StructDecl{}
		class.fields = make(map[string]field)
		for i := 0; i < s.NumFields(); i++ {
			f := s.Field(i)

			switch p := f.Type().(type) {
			case *types.Signature:
				class.fields[f.Name()] = field{
					kind: fieldKindFunc,
					fn:   p.String(),
				}
			case *types.Interface:
				var methods []string
				for i := 0; i < p.NumMethods(); i++ {
					methods = append(methods, p.Method(i).Name())
				}
				class.fields[f.Name()] = field{
					kind:    fieldKindInterface,
					methods: methods,
				}
			}
		}

		// Iterate through all methods of the class.
		for i := 0; i < tp.NumMethods(); i++ {
			class.methods = append(class.methods, getFuncAsString(tp.Method(i)))
		}

		classes[tp.Obj().Name()] = class
	}
}

func getFuncAsString(method *types.Func) string {
	ret := tupleAsString(method.Type().(*types.Signature).Results())
	if ret == "" {
		return fmt.Sprintf("%s(%s)", method.Name(), tupleAsString(method.Type().(*types.Signature).Params()))
	}
	return fmt.Sprintf("%s(%s) (%s)", method.Name(), tupleAsString(method.Type().(*types.Signature).Params()), ret)
}

func tupleAsString(tuple *types.Tuple) string {
	var params []string
	for i := 0; i < tuple.Len(); i++ {
		param := tuple.At(i)
		params = append(params, fmt.Sprintf("%s %s", param.Name(), param.Type()))
	}
	return strings.Join(params, ", ")
}
