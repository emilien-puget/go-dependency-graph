package struct_decl

import (
	"go/types"

	"golang.org/x/tools/go/packages"
)

type Decl struct {
	Fields          map[string]Field
	Methods         []Method
	ActualNamedType *types.Named
}

const (
	fieldKindInterface = "interface"
	fieldKindFunc      = "func"
)

type Field struct {
	Kind    string
	Methods []string
	Fn      string
}

// Extract extracts struct declaration from packages.
func Extract(pkgs []*packages.Package) map[string]map[string]*Decl {
	declarations := make(map[string]map[string]*Decl)
	for i := range pkgs {
		pkgType := extractTypes(pkgs[i])
		declarations[pkgs[i].Name] = pkgType
	}
	return declarations
}

func extractTypes(pkg *packages.Package) map[string]*Decl {
	classes := make(map[string]*Decl)

	// Iterate through all types in the package.
	for _, typ := range pkg.TypesInfo.Defs {
		readTypeObject(typ, classes)
	}

	return classes
}

func readTypeObject(typ types.Object, classes map[string]*Decl) {
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
	name := tp.Obj().Name()

	class := &Decl{}
	class.ActualNamedType = tp
	class.Fields = make(map[string]Field)
	for i := 0; i < s.NumFields(); i++ {
		f := s.Field(i)

		switch p := f.Type().(type) {
		case *types.Signature: // struct Field is a func.
			class.Fields[f.Name()] = Field{
				Kind: fieldKindFunc,
				Fn:   p.String(),
			}
		case *types.Interface: // struct Field is an anonymous interface.
			readInterface(p, class, f)
		case *types.Named: // struct Field is a named type
			ni, ok := p.Underlying().(*types.Interface) // the named type is an interface
			if !ok {
				continue
			}
			readInterface(ni, class, f)
		}
	}

	// Iterate through all methods of the class.
	for i := 0; i < tp.NumMethods(); i++ {
		class.Methods = append(class.Methods, Method{TypFuc: tp.Method(i)})
	}

	classes[name] = class
}

func readInterface(p *types.Interface, class *Decl, f *types.Var) {
	var methods []string
	for i := 0; i < p.NumMethods(); i++ {
		methods = append(methods, p.Method(i).Name())
	}
	class.Fields[f.Name()] = Field{
		Kind:    fieldKindInterface,
		Methods: methods,
	}
}
