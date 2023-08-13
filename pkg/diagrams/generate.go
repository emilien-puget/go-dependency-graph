package diagrams

import (
	"bufio"
	"errors"

	"github.com/emilien-puget/go-dependency-graph/pkg/diagrams/c4"
	"github.com/emilien-puget/go-dependency-graph/pkg/diagrams/mermaid"
	"github.com/emilien-puget/go-dependency-graph/pkg/parse"
)

type generateFromSchema func(writer *bufio.Writer, s parse.AstSchema) error

const (
	GeneratorC4PlantumlComponent = "c4_plantuml_component"
	GeneratorMermaidClass        = "mermaid_class"
)

var generators = map[string]generateFromSchema{
	GeneratorC4PlantumlComponent: c4.GenerateComponentFromSchema,
	GeneratorMermaidClass:        mermaid.GenerateClassFromSchema,
}

var errUnknownGenerator = errors.New("unknown generator")

func Generate(generator string, writer *bufio.Writer, as parse.AstSchema) error {
	gen, ok := generators[generator]
	if !ok {
		return errUnknownGenerator
	}

	err := gen(writer, as)
	if err != nil {
		return err
	}
	return nil
}
