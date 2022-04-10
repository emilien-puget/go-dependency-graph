package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/emilien-puget/go-dependency-graph/pkg/generator/c4"
	"github.com/emilien-puget/go-dependency-graph/pkg/generator/json"
	"github.com/emilien-puget/go-dependency-graph/pkg/parse"
)

func main() {
	project := flag.String("project", "", "the path of the project to inspect")
	path := flag.String("result", "", "the path of the generated file, not used if stdout is piped")
	generator := flag.String("generator", "c4_plantuml_component", "the name of the generator to use, [c4_plantuml_component, json], default c4_plantuml_component")
	flag.Parse()

	err := run(project, path, generator)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var (
	errMissingResult  = errors.New("result is required")
	errMissingProject = errors.New("project is required")
)

func run(project, path, generator *string) error {
	if project == nil || *project == "" {
		return errMissingProject
	}

	writer, err := getWriter(path)
	if err != nil {
		return err
	}
	defer writer.Flush()

	as, err := parse.Parse(*project)
	if err != nil {
		return err
	}

	switch *generator {
	case "c4_plantuml_component":
		err = c4.GenerateC4ComponentUmlFromSchema(writer, as)
		if err != nil {
			return err
		}
	case "json":
		err := json.GenerateJSONFromSchema(writer, as)
		if err != nil {
			return err
		}
	default:
		_, _ = fmt.Fprintln(os.Stderr, "unknown generator")
		return nil
	}
	return nil
}

func getWriter(path *string) (*bufio.Writer, error) {
	o, _ := os.Stdout.Stat()
	if (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
		if path == nil || *path == "" {
			return nil, errMissingResult
		}
		file, err := os.Create(*path)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		return bufio.NewWriter(file), nil
	}
	return bufio.NewWriter(os.Stdout), nil
}
