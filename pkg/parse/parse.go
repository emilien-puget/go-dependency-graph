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
	Methods      []string
	Comment      string
	Imports      []Import
	Deps         map[string][]Dep
	ExternalDeps map[string][]Dep
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

	cfg := &packages.Config{Dir: pathDir, Mode: packages.NeedName | packages.NeedImports | packages.NeedDeps | packages.NeedExportFile | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo, Tests: false}
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
	types := make(map[string]map[string]*StructDecl)
	typesInter := make(map[string]map[string]*InterDecl)
	pkgs, err := packages.Load(cfg, dirs...)
	for i := range pkgs {
		pkgType, pkgInter := ExtractTypes(pkgs[i])
		if err != nil {
			return AstSchema{}, err
		}
		types[pkgs[i].Name] = pkgType
		typesInter[pkgs[i].Name] = pkgInter
	}

	assoInterStruct(typesInter, types)
	for i := range pkgs {
		parsePackage(pkgs[i], &as, types)
		if err != nil {
			return AstSchema{}, err
		}
	}

	return as, nil
}
func assoInterStruct(interds map[string]map[string]*InterDecl, structds map[string]map[string]*StructDecl) {
	for s := range interds {
		for interd := range interds[s] {

			for s2 := range structds {
				for s3 := range structds[s2] {
					funcName(interds, structds, s2, s3, s, interd)
				}
			}

		}
	}
}

func funcName(interds map[string]map[string]*InterDecl, structds map[string]map[string]*StructDecl, s2 string, s3 string, s string, interd string) {
	if len(structds[s2][s3].methods) == len(interds[s][interd].methods) {
		for i := range structds[s2][s3].methods {
			if structds[s2][s3].methods[i] != interds[s][interd].methods[i] {
				return
			}
		}
		_, ok := interds[s][interd].implems[s2]
		if !ok {
			interds[s][interd].implems[s2] = make(map[string]*StructDecl)
		}
		interds[s][interd].implems[s2][s3] = structds[s2][s3]
	}
}

func parsePackage(p *packages.Package, as *AstSchema, types map[string]map[string]*StructDecl) {
	for _, f := range p.Syntax {
		dependencies := parseFile(f, p, as.ModulePath, types)
		for depName, dep := range dependencies {
			if _, ok := as.Packages[p.Name]; !ok {
				as.Packages[p.Name] = make(Dependencies)
			}
			as.Packages[p.Name][depName] = dep
		}
	}
}

func parseFile(f *ast.File, p *packages.Package, modulePath string, types map[string]map[string]*StructDecl) (dependencies Dependencies) {
	dependencies = make(Dependencies, 0)
	packageName := f.Name.Name
	structs := map[string]string{}
	for _, decl := range f.Decls {
		if d, ok := decl.(*ast.GenDecl); ok {
			name, structDecl := searchStructDeclDoc(d)
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
		name, deps, sDecl := searchProvider(d, packageName, imports, p.TypesInfo, types)
		if name == "" {
			continue
		}
		ser := Dependency{}
		ser.Deps = deps
		ser.Methods = sDecl.methods
		if len(structs[packageName+"."+name]) > 3 {
			ser.Comment = structs[packageName+"."+name][3:]
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
