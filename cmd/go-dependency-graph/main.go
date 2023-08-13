package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/emilien-puget/go-dependency-graph/pkg/diagrams"
	"github.com/emilien-puget/go-dependency-graph/pkg/parse"
	writer "github.com/emilien-puget/go-dependency-graph/pkg/writer"
)

func main() {
	project := flag.String("project", "", "the path of the project to inspect")
	path := flag.String("result", "", "the path of the generated file, not used if stdout is piped")
	generator := flag.String("generator", "c4_plantuml_component", "the name of the generator to use, [c4_plantuml_component, mermaid_class], default c4_plantuml_component")
	flag.Parse()

	err := run(project, path, generator)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var (
	errMissingProject   = errors.New("project is required")
	errMissingGenerator = errors.New("generator is required")
)

func run(project, path, generator *string) error {
	if project == nil || *project == "" {
		return errMissingProject
	}

	if generator == nil || *generator == "" {
		return errMissingGenerator
	}

	w, closer, err := writer.GetWriter(path)
	if err != nil {
		return err
	}
	defer closer()

	as, err := parse.Parse(*project)
	if err != nil {
		return err
	}

	err = diagrams.Generate(*generator, w, as)
	if err != nil {
		return err
	}
	return nil
}
