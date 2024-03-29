package parse

import (
	"fmt"
	"go/ast"
	"path/filepath"

	"github.com/emilien-puget/go-dependency-graph/pkg/parse/package_list"
	"github.com/emilien-puget/go-dependency-graph/pkg/parse/struct_decl"
	"golang.org/x/tools/go/packages"
)

// AstSchema is a simpler presentation of the ast of a project.
type AstSchema struct {
	ModulePath string
	Graph      *Graph
}

// Parse parses the project located under pathDir and returns an AstSchema.
func Parse(pathDir string, skipDirs []string) (AstSchema, error) {
	pathDir, err := filepath.Abs(pathDir)
	if err != nil {
		return AstSchema{}, fmt.Errorf("filepath.Abs:%w", err)
	}
	modulePath, err := getModulePath(pathDir)
	if err != nil {
		return AstSchema{}, fmt.Errorf("getModulePath:%w", err)
	}
	as := AstSchema{
		ModulePath: modulePath,
		Graph:      NewGraph(),
	}

	pkgs, err := package_list.GetPackagesToParse(pathDir, skipDirs)
	if err != nil {
		return AstSchema{}, fmt.Errorf("package_list.GetPackagesToParse:%w", err)
	}

	types := struct_decl.Extract(pkgs)
	if err != nil {
		return AstSchema{}, fmt.Errorf("struct_decl.Extract:%w", err)
	}

	parsePackages(pkgs, &as, types)

	return as, nil
}

func parsePackages(pkgs []*packages.Package, schema *AstSchema, types map[string]map[string]*struct_decl.Decl) {
	for i := range pkgs {
		parsePackage(pkgs[i], schema, types)
	}
}

func parsePackage(p *packages.Package, schema *AstSchema, types map[string]map[string]*struct_decl.Decl) {
	for _, f := range p.Syntax {
		parseFile(f, p, schema.ModulePath, types, schema.Graph)
	}
}

func parseFile(f *ast.File, p *packages.Package, modulePath string, types map[string]map[string]*struct_decl.Decl, graph *Graph) {
	packageName := p.ID

	structDoc := struct_decl.GetStructDoc(f, packageName)

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
			Name:            packageName + "." + name,
			PackageName:     packageName,
			StructName:      name,
			Methods:         sDecl.Methods,
			ActualNamedType: sDecl.ActualNamedType,
			P:               p,
			FilePath:        sDecl.FilePath,
		}
		if len(structDoc[packageName+"."+name]) > 3 {
			newNode.Doc = structDoc[packageName+"."+name][3:]
		}
		graph.AddNode(newNode)

		for s := range deps {
			for i2 := range deps[s] {
				adjNode := &Node{
					Name:        deps[s][i2].PackageName + "." + deps[s][i2].DependencyName,
					PackageName: deps[s][i2].PackageName,
					StructName:  deps[s][i2].DependencyName,
					External:    deps[s][i2].External,
				}
				graph.AddNode(adjNode)
				graph.AddEdge(newNode, &Adj{
					Node: adjNode,
					Func: deps[s][i2].Funcs,
				})
			}
		}
	}
}
