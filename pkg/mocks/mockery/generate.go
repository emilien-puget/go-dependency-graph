package mockery

import (
	"context"
	"fmt"
	"go/types"
	"os"
	"path/filepath"

	"github.com/emilien-puget/go-dependency-graph/pkg/parse"
	"github.com/vektra/mockery/v2/pkg"
)

type Generator struct {
	OutOfPackageMocksDirectory string
}

func NewGenerator(outOfPackageMocksDirectory string) *Generator {
	return &Generator{
		OutOfPackageMocksDirectory: outOfPackageMocksDirectory,
	}
}

func (g Generator) GenerateFromSchema(as parse.AstSchema) error {
	err := os.MkdirAll(g.OutOfPackageMocksDirectory, os.FileMode(0o755))
	if err != nil {
		return fmt.Errorf("os.MkdirAll:%w", err)
	}
	for _, node := range as.Graph.TopologicalSort() {
		if len(node.InboundEdges) == 0 {
			continue
		}

		if node.ActualNamedType == nil {
			continue
		}

		err := g.generateMockForNode(node)
		if err != nil {
			return fmt.Errorf("g.generateMockForNode:%w", err)
		}
	}
	return nil
}

func (g Generator) generateMockForNode(node *parse.Node) error {
	funcs := make([]*types.Func, 0, len(node.Methods))
	for i := range node.Methods {
		funcs = append(funcs, node.Methods[i].TypFuc)
	}
	generator := pkg.NewGenerator(
		context.Background(),
		pkg.GeneratorConfig{
			DisableVersionString: true,
			Exported:             true,
			InPackage:            false,
			KeepTree:             false,
			WithExpecter:         true,
		},
		&pkg.Interface{
			Name:            node.PackageName + node.StructName,
			Pkg:             node.P.Types,
			NamedType:       node.ActualNamedType,
			ActualInterface: types.NewInterfaceType(funcs, nil),
		},
		"mocks",
	)

	err := generator.GenerateAll(context.Background())
	if err != nil {
		return err
	}

	file, err := os.OpenFile(g.determineMockFilePath(node), os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(0o644))
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

func (g Generator) determineMockFilePath(node *parse.Node) string {
	return filepath.Join(g.OutOfPackageMocksDirectory, node.PackageName+node.StructName+".go")
}
