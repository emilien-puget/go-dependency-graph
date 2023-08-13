package mockery

import (
	"context"
	"go/types"
	"os"
	"path/filepath"

	"github.com/emilien-puget/go-dependency-graph/pkg/mocks/config"
	"github.com/emilien-puget/go-dependency-graph/pkg/parse"
	"github.com/vektra/mockery/v2/pkg"
)

func GenerateFromSchema(c config.Config, as parse.AstSchema) error {
	for _, node := range as.Graph.TopologicalSort() {
		if len(node.InboundEdges) == 0 {
			continue
		}

		if node.ActualNamedType == nil {
			continue
		}

		err := generateMockForNode(c, node)
		if err != nil {
			return err
		}
	}
	return nil
}

func generateMockForNode(c config.Config, node *parse.Node) error {
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
			KeepTree:             c.InPackage,
			PackageName:          determinePackage(c, node),
			WithExpecter:         true,
		},
		&pkg.Interface{
			Name:            node.PackageName + node.StructName,
			Pkg:             node.P.Types,
			NamedType:       node.ActualNamedType,
			ActualInterface: types.NewInterfaceType(funcs, nil),
		},
		determinePackageName(c, node),
	)

	err := generator.GenerateAll(context.Background())
	if err != nil {
		return err
	}

	file, err := os.OpenFile(determineMockFilePath(c, node), os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(0o644))
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

func determinePackage(c config.Config, node *parse.Node) string {
	if c.InPackage {
		return node.PackageName + "_test"
	}
	return node.PackageName
}

func determineMockFilePath(c config.Config, node *parse.Node) string {
	if !c.InPackage {
		return filepath.Join(c.OutOfPackageMocksDirectory, node.PackageName+node.StructName+".go")
	}
	dir, file := filepath.Split(node.FilePath)

	fileName := file[:len(file)-len(filepath.Ext(file))]

	return filepath.Join(dir, fileName+"_mocks_test.go")
}

func determinePackageName(c config.Config, n *parse.Node) string {
	if !c.InPackage {
		return "mocks"
	}

	return n.PackageName + "_test"
}
