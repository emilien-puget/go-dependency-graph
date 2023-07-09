package parse

import (
	"fmt"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

// ExtractTypes extracts types information from a package.
func ExtractTypes(pkg *packages.Package) map[string]*StructDecl {
	classes := make(map[string]*StructDecl)

	// Iterate through all types in the package.
	for _, typ := range pkg.TypesInfo.Defs {
		readTypeObject(typ, classes)
	}

	return classes
}

func readTypeObject(typ types.Object, classes map[string]*StructDecl) {
	if typ == nil {
		return
	}

	tp, ok := typ.Type().(*types.Named)
	if !ok {
		return
	}
	s, ok := tp.Underlying().(*types.Struct)
	if !ok {
		return
	}

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