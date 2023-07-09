package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/emilien-puget/go-dependency-graph/pkg/generator/c4"
	"github.com/emilien-puget/go-dependency-graph/pkg/generator/json"
	"github.com/emilien-puget/go-dependency-graph/pkg/generator/mermaid"
	"github.com/emilien-puget/go-dependency-graph/pkg/parse"
)

func main() {
	project := flag.String("project", "", "the path of the project to inspect")
	path := flag.String("result", "", "the path of the generated file, not used if stdout is piped")
	generator := flag.String("generator", "c4_plantuml_component", "the name of the generator to use, [c4_plantuml_component, mermaid_class, json], default c4_plantuml_component")
	flag.Parse()

	err := run(project, path, generator)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var (
	errMissingResult    = errors.New("result is required")
	errMissingProject   = errors.New("project is required")
	errUnknownGenerator = errors.New("unknown generator")
)

type generateFromSchema func(writer *bufio.Writer, s parse.AstSchema) error

var generators = map[string]generateFromSchema{
	"c4_plantuml_component": c4.GenerateComponentFromSchema,
	"json":                  json.GenerateFromSchema,
	"mermaid_class":         mermaid.GenerateClassFromSchema,
}

func run(project, path, generator *string) error {
	if project == nil || *project == "" {
		return errMissingProject
	}

	gen, ok := generators[*generator]
	if !ok {
		return errUnknownGenerator
	}

	writer, closer, err := getWriter(path)
	if err != nil {
		return err
	}
	defer closer()

	as, err := parse.Parse(*project)
	if err != nil {
		return err
	}

	err = gen(writer, as)
	if err != nil {
		return err
	}
	return nil
}

func getWriter(path *string) (*bufio.Writer, func(), error) {
	o, _ := os.Stdout.Stat()
	if (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
		if path == nil || *path == "" {
			return nil, nil, errMissingResult
		}
		file, err := os.Create(*path)
		if err != nil {
			return nil, nil, err
		}
		writer := bufio.NewWriter(file)
		return writer, func() {
			_ = writer.Flush()
			_ = file.Close()
		}, nil
	}
	writer := bufio.NewWriter(os.Stdout)
	return writer, func() {
		_ = writer.Flush()
	}, nil
}
