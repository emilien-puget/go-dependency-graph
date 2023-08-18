package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/emilien-puget/go-dependency-graph/pkg/config"
	"github.com/emilien-puget/go-dependency-graph/pkg/diagrams"
	mocksconfig "github.com/emilien-puget/go-dependency-graph/pkg/mocks/config"
	"github.com/emilien-puget/go-dependency-graph/pkg/parse"
	"github.com/emilien-puget/go-dependency-graph/pkg/writer"
)

func main() {
	project := flag.String("project", "", "the path of the project to inspect")
	path := flag.String("result", "", "the path of the generated file, not used if stdout is piped")
	generator := flag.String("generator", "c4_plantuml_component", "the name of the generator to use, [c4_plantuml_component, mermaid_class], default c4_plantuml_component")
	skipFolders := flag.String("skip-dirs", mocksconfig.DefaultOutOfPackageDirectory+","+config.VendorDir, "a comma separate list of directory to ignore, default value is the mocks and vendor directory")
	flag.Parse()

	err := run(project, path, generator, skipFolders)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var (
	errMissingProject   = errors.New("project is required")
	errMissingGenerator = errors.New("generator is required")
)

func run(project, path, generator, f *string) error {
	err := validateRequiredInput(project, generator)
	if err != nil {
		return err
	}

	w, closer, err := writer.GetWriter(path)
	if err != nil {
		return err
	}
	defer closer()

	var skipDirs []string
	if f != nil || *f != "" {
		skipDirs = strings.Split(*f, ",")
	}

	as, err := parse.Parse(*project, skipDirs)
	if err != nil {
		return err
	}

	err = diagrams.Generate(*generator, w, as)
	if err != nil {
		return err
	}
	return nil
}

func validateRequiredInput(project, generator *string) error {
	if project == nil || *project == "" {
		return errMissingProject
	}

	if generator == nil || *generator == "" {
		return errMissingGenerator
	}
	return nil
}
