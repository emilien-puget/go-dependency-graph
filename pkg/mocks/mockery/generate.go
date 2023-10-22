package mockery

import (
	"context"
	"fmt"
	"go/types"
	"os"
	"path/filepath"
	"strings"

	"github.com/emilien-puget/go-dependency-graph/pkg/parse"
	"github.com/vektra/mockery/v2/pkg"
)

type Generator struct {
	OutOfPackageMocksDirectory string
	Replacer                   *strings.Replacer
}

func NewGenerator(outOfPackageMocksDirectory string) *Generator {
	return &Generator{
		Replacer:                   strings.NewReplacer("/", "_"),
		OutOfPackageMocksDirectory: outOfPackageMocksDirectory,
	}
}

func (g Generator) GenerateFromSchema(ctx context.Context, as parse.AstSchema) error {
	err := os.MkdirAll(g.OutOfPackageMocksDirectory, os.FileMode(0o755))
	if err != nil {
		return fmt.Errorf("os.MkdirAll:%w", err)
	}
	for _, node := range as.Graph.TopologicalSort() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if len(node.InboundEdges) == 0 {
			continue
		}

		if node.ActualNamedType == nil {
			continue
		}

		err := g.generateMockForNode(ctx, as.ModulePath, node)
		if err != nil {
			return fmt.Errorf("g.generateMockForNode:%w", err)
		}
	}
	return nil
}

func (g Generator) generateMockForNode(ctx context.Context, path string, node *parse.Node) error {
	funcs := make([]*types.Func, 0, len(node.Methods))
	for i := range node.Methods {
		if node.Methods[i].TypFuc.Exported() {
			funcs = append(funcs, node.Methods[i].TypFuc)
		}
	}
	name := strings.TrimPrefix(node.PackageName+node.StructName, path)
	name = strings.TrimPrefix(name, "/")
	name = g.Replacer.Replace(name)
	generator := pkg.NewGenerator(
		ctx,
		pkg.GeneratorConfig{
			DisableVersionString: true,
			Exported:             true,
			InPackage:            false,
			KeepTree:             false,
			WithExpecter:         true,
		},
		&pkg.Interface{
			Name:            name,
			Pkg:             node.P.Types,
			NamedType:       node.ActualNamedType,
			ActualInterface: types.NewInterfaceType(funcs, nil),
		},
		"mocks",
	)

	err := generator.GenerateAll(ctx)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(g.determineMockFilePath(name), os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(0o644))
	if err != nil {
		return err
	}
	defer file.Close()
	err = generator.Write(file)
	if err != nil {
		return err
	}

	return nil
}

func (g Generator) determineMockFilePath(node string) string {
	return filepath.Join(g.OutOfPackageMocksDirectory, node+".go")
}
