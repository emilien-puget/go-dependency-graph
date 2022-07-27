package parse

import (
	"go/ast"
	"io/fs"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

// AstSchema is a simpler presentation of the ast of a project.
type AstSchema struct {
	ModulePath string
	Packages   map[string]Dependencies
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
	Funcs          []string
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

	cfg := &packages.Config{Dir: pathDir, Mode: packages.NeedName | packages.NeedImports | packages.NeedDeps | packages.NeedExportsFile | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo, Tests: false}
	var dirs []string
	err = filepath.WalkDir(pathDir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".go") {
			dir, _ := filepath.Split(p)
			dirs = append(dirs, dir)
			return nil
		}

		return nil
	})
	if err != nil {
		return AstSchema{}, err
	}
	pkgs, err := packages.Load(cfg, dirs...)
	for i := range pkgs {
		parsePackage(pkgs[i], &as)
		if err != nil {
			return AstSchema{}, err
		}
	}

	return as, nil
}

func parsePackage(p *packages.Package, as *AstSchema) {
	for _, f := range p.Syntax {
		dependencies := parseFile(f, p, as.ModulePath)
		for depName, dep := range dependencies {
			if _, ok := as.Packages[p.Name]; !ok {
				as.Packages[p.Name] = make(Dependencies)
			}
			as.Packages[p.Name][depName] = dep
		}
	}
}

func parseFile(f *ast.File, p *packages.Package, modulePath string) (dependencies Dependencies) {
	dependencies = make(Dependencies, 0)
	packageName := f.Name.Name
	structs := map[string]structDecl{}
	for _, decl := range f.Decls {
		if d, ok := decl.(*ast.GenDecl); ok {
			name, structDecl := searchStructDecl(d)
			structs[packageName+"."+name] = structDecl
		}
	}

	imports := parseImports(f, modulePath, p.Imports)
	for _, decl := range f.Decls {
		d, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		var deps map[string][]Dep
		name, deps := searchProvider(d, structs, packageName, imports, p.TypesInfo)
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

func parseImports(f *ast.File, modulePath string, imp map[string]*packages.Package) map[string]Import {
	imports := make(map[string]Import, 0)
	for _, im := range f.Imports {
		p := strings.Trim(im.Path.Value, "\"")
		importName := imp[p].Name
		if im.Name != nil {
			importName = im.Name.Name
		}
		imports[importName] = Import{
			Path:     p,
			External: !strings.Contains(p, modulePath),
		}
	}
	return imports
}
