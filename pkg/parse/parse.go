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
	Graph      *Graph
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
		Graph:      NewGraph(),
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
	pkgs, err := packages.Load(cfg, dirs...)
	for i := range pkgs {
		pkgType := ExtractTypes(pkgs[i])
		if err != nil {
			return AstSchema{}, err
		}
		types[pkgs[i].Name] = pkgType
	}
	for i := range pkgs {
		parsePackage(pkgs[i], &as, types)
		if err != nil {
			return AstSchema{}, err
		}
	}

	return as, nil
}

func parsePackage(p *packages.Package, as *AstSchema, types map[string]map[string]*StructDecl) {
	for _, f := range p.Syntax {
		parseFile(f, p, as.ModulePath, types, as.Graph)
	}
}

func parseFile(f *ast.File, p *packages.Package, modulePath string, types map[string]map[string]*StructDecl, graph *Graph) {
	packageName := f.Name.Name

	structDoc := map[string]string{}
	for _, decl := range f.Decls {
		if d, ok := decl.(*ast.GenDecl); ok {
			name, structDecl := searchStructDecl(d)
			structDoc[packageName+"."+name] = structDecl
		}
	}

	imports := parseImports(f, modulePath, p.Imports)
	for _, decl := range f.Decls {
		d, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		name, deps, sDecl := searchProvider(d, packageName, imports, p.TypesInfo, types)
		if name == "" {
			continue
		}
		newNode := &Node{
			Name:        packageName + "." + name,
			PackageName: packageName,
			StructName:  name,
			Methods:     sDecl.methods,
		}
		if len(structDoc[packageName+"."+name]) > 3 {
			newNode.Doc = structDoc[packageName+"."+name][3:]
		}
		graph.AddNode(newNode)

		for s := range deps {
			for i2 := range deps[s] {
				graph.AddEdge(newNode, &Adj{
					Node: &Node{
						Name:        deps[s][i2].PackageName + "." + deps[s][i2].DependencyName,
						PackageName: deps[s][i2].PackageName,
						StructName:  deps[s][i2].DependencyName,
						External:    deps[s][i2].External,
					},
					Func: deps[s][i2].Funcs,
				})
			}
		}
	}
	return
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
