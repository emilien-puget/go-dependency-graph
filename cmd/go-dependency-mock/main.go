package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/emilien-puget/go-dependency-graph/pkg/mocks"
	"github.com/emilien-puget/go-dependency-graph/pkg/mocks/config"
	"github.com/emilien-puget/go-dependency-graph/pkg/parse"
)

func main() {
	project := flag.String("project", "", "the path of the project to inspect")
	generator := flag.String("generator", "mockery", "the name of the generator to use, [mockery], default mockery")
	inPackage := flag.Bool("in-package", false, "whether the mocks are written in a specific package or in the same package as their struct")
	outOfPackageMocksDirectory := flag.String("mocks-dir", config.DefaultOutOfPackageDirectory, "where the mocks will be written, only used if in-package is false")

	flag.Parse()

	err := run(project, generator, outOfPackageMocksDirectory, *inPackage)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var (
	errMissingProject        = errors.New("project is required")
	errMissingGenerator      = errors.New("generator is required")
	errMissingMocksDirectory = errors.New("mocks-folder is required")
)

func run(project, generator, outOfPackageMocksDirectory *string, inPackage bool) error {
	if project == nil || *project == "" {
		return errMissingProject
	}

	if generator == nil || *generator == "" {
		return errMissingGenerator
	}

	c := config.Config{
		InPackage: inPackage,
	}

	if !inPackage {
		if outOfPackageMocksDirectory == nil || *outOfPackageMocksDirectory == "" {
			return errMissingMocksDirectory
		}
		c.OutOfPackageMocksDirectory = *project + string(filepath.Separator) + *outOfPackageMocksDirectory
	}

	as, err := parse.Parse(*project)
	if err != nil {
		return err
	}

	err = mocks.Generate(
		*generator,
		c,
		as,
	)
	if err != nil {
		return err
	}

	return nil
}
