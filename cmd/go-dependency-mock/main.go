package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/emilien-puget/go-dependency-graph/pkg/config"
	"github.com/emilien-puget/go-dependency-graph/pkg/mocks"
	mocksconfig "github.com/emilien-puget/go-dependency-graph/pkg/mocks/config"
	"github.com/emilien-puget/go-dependency-graph/pkg/parse"
)

func main() {
	project := flag.String("project", "", "the path of the project to inspect")
	generator := flag.String("generator", "mockery", "the name of the generator to use, [mockery], default mockery")
	outOfPackageMocksDirectory := flag.String("mocks-dir", mocksconfig.DefaultOutOfPackageDirectory, "where the mocks will be written")
	skipFolders := flag.String("skip-dirs", mocksconfig.DefaultOutOfPackageDirectory+","+config.VendorDir, "a comma separate list of directory to ignore, default value is the mocks and vendor directory")

	flag.Parse()

	err := run(project, generator, outOfPackageMocksDirectory, skipFolders)
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

func run(project, generator, outOfPackageMocksDirectory, f *string) error {
	err := validateRequiredInput(project, generator)
	if err != nil {
		return err
	}

	c := mocksconfig.Config{}

	if outOfPackageMocksDirectory == nil || *outOfPackageMocksDirectory == "" {
		return errMissingMocksDirectory
	}
	c.OutOfPackageMocksDirectory = filepath.Join(*project, *outOfPackageMocksDirectory)

	var skipDirs []string
	if f != nil || *f != "" {
		skipDirs = strings.Split(*f, ",")
	}

	as, err := parse.Parse(*project, skipDirs)
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

func validateRequiredInput(project, generator *string) error {
	if project == nil || *project == "" {
		return errMissingProject
	}

	if generator == nil || *generator == "" {
		return errMissingGenerator
	}
	return nil
}
