package parse

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
)

// AstSchema is a simpler presentation of the ast of a project.
type AstSchema struct {
	ModulePath string
	Packages   map[string]Dependencies
	Decl       map[string][]string
}

// Dependencies contains all the dependencies of a package.
type Dependencies map[string]Dependency

// Dependency represent a type that has been identified as a dependency.
type Dependency struct {
	Comment string
	Imports []Import
	Deps    map[string][]Dep
}

// Import represent an imported package.
type Import struct {
	Path     string
	External bool
}

// Dep represent one dependency injected.
type Dep struct {
	PackageName    string
	DependencyName string
	VarName        string
	Funcs          map[string]string
	External       bool
}

// Parse parses the project located under pathDir and returns an AstSchema.
func Parse(pathDir string) (AstSchema, error) {
	modulePath, err := getModulePath(pathDir)
	if err != nil {
		return AstSchema{}, err
	}
	as := AstSchema{
		ModulePath: modulePath,
		Packages:   map[string]Dependencies{},
	}

	err = filepath.WalkDir(pathDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			panic(err)
		}
		if d.IsDir() {
			err := parseDir(path, &as)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return AstSchema{}, err
	}

	return as, nil
}

func parseDir(path string, as *AstSchema) error {
	fset := token.NewFileSet()

	dir, err := parser.ParseDir(fset, path, func(info fs.FileInfo) bool {
		n := info.Name()
		if len(n) > 8 && n[len(n)-8:] == "_test.go" {
			return false
		}
		return true
	}, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return err
	}

	for name, p := range dir {
		m := map[string]Dependency{}
		for _, f := range p.Files {
			dependencies := parseFile(f, as.ModulePath)
			for depName, dep := range dependencies {
				m[depName] = dep
			}
		}

		if len(m) != 0 {
			as.Packages[name] = m
		}
	}
	return nil
}

func parseFile(f *ast.File, modulePath string) (dependencies Dependencies) {
	dependencies = make(Dependencies, 0)
	packageName := f.Name.Name
	structs := map[string]structDecl{}
	for _, decl := range f.Decls {
		if d, ok := decl.(*ast.GenDecl); ok {
			name, structDecl := searchStructDecl(d)
			structs[packageName+"."+name] = structDecl
		}
	}

	imports := parseImports(f, modulePath)
	for _, decl := range f.Decls {
		d, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		var deps map[string][]Dep
		name, deps := searchProvider(d, structs, packageName, imports)
		if name == "" {
			continue
		}
		ser := Dependency{}
		ser.Deps = deps
		if len(structs[packageName+"."+name].doc) > 3 {
			ser.Comment = structs[packageName+"."+name].doc[3:]
		}
		dependencies[name] = ser
	}

	return dependencies
}

func parseImports(f *ast.File, modulePath string) map[string]Import {
	imports := make(map[string]Import, 0)
	for _, im := range f.Imports {
		path := strings.Trim(im.Path.Value, "\"")
		split := strings.Split(path, "/")
		importName := split[len(split)-1]
		if im.Name != nil {
			importName = im.Name.Name
		}
		imports[importName] = Import{
			Path:     path,
			External: !strings.Contains(path, modulePath),
		}
	}
	return imports
}
