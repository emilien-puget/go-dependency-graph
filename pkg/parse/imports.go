package parse

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/packages"
)

// importDecl represent an imported package.
type importDecl struct {
	Path     string
	External bool
}

func parseImports(f *ast.File, modulePath string, imp map[string]*packages.Package) map[string]importDecl {
	imports := make(map[string]importDecl)
	for _, im := range f.Imports {
		p := strings.Trim(im.Path.Value, "\"")
		importName := imp[p].Name
		if im.Name != nil {
			importName = im.Name.Name
		}
		imports[importName] = importDecl{
			Path:     p,
			External: !strings.Contains(p, modulePath),
		}
	}
	return imports
}
