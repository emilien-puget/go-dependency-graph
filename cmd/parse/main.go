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

	if project == nil || *project == "" {
		_, _ = fmt.Fprintln(os.Stderr, "project is required")
		os.Exit(1)
	}

	writer, err := getWriter(path)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer writer.Flush()

	as, err := parse.Parse(*project)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch *generator {
	case "c4_plantuml_component":
		err = c4.GenerateC4ComponentUmlFromSchema(writer, as)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "json":
		err := json.GenerateJSONFromSchema(writer, as)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		_, _ = fmt.Fprintln(os.Stderr, "unknown generator")
		os.Exit(1)
	}
}

func getWriter(path *string) (*bufio.Writer, error) {
	o, _ := os.Stdout.Stat()
	if (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
		if path == nil || *path == "" {
			return nil, errors.New("result is required")
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
