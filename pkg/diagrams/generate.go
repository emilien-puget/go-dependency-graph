package diagrams

import (
	"bufio"
	"context"
	"errors"

	"github.com/emilien-puget/go-dependency-graph/pkg/diagrams/c4"
	"github.com/emilien-puget/go-dependency-graph/pkg/diagrams/mermaid"
	"github.com/emilien-puget/go-dependency-graph/pkg/parse"
)

type Generator interface {
	GenerateFromSchema(ctx context.Context, writer *bufio.Writer, s parse.AstSchema) error
	GetDefaultResultFileName() string
}

const (
	GeneratorC4PlantumlComponent = "c4_plantuml_component"
	GeneratorMermaidClass        = "mermaid_class"
)

var errUnknownGenerator = errors.New("unknown generator")

func GetGenerator(generator string) (Generator, error) {
	switch generator {
	case GeneratorC4PlantumlComponent:
		return c4.NewGenerator(), nil
	case GeneratorMermaidClass:
		return mermaid.NewGenerator(), nil
	default:
		return nil, errUnknownGenerator
	}
}
