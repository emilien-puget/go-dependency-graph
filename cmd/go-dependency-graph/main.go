package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/emilien-puget/go-dependency-graph/pkg/config"
	"github.com/emilien-puget/go-dependency-graph/pkg/diagrams"
	"github.com/emilien-puget/go-dependency-graph/pkg/mocks"
	mocksconfig "github.com/emilien-puget/go-dependency-graph/pkg/mocks/config"
	"github.com/emilien-puget/go-dependency-graph/pkg/parse"
	"github.com/emilien-puget/go-dependency-graph/pkg/writer"
	"golang.org/x/sync/errgroup"
)

func main() {
	project := flag.String("project", "", "the path of the project to inspect, default is current dir")
	diagEnable := flag.Bool("generate-diag", true, "enable diagram generation, default is true")
	mocksEnable := flag.Bool("generate-mocks", true, "enable mocks generation, default is true")
	diagResult := flag.String("diag-result", "", "the path of the generated file, not used if stdout is piped, default is in the project dir")
	diagGenerator := flag.String("diag-generator", "c4_plantuml_component", "the name of the generator to use, [c4_plantuml_component, mermaid_class], default c4_plantuml_component")
	mockGenerator := flag.String("mock-generator", "mockery", "the name of the generator to use, [mockery], default mockery")
	mockResult := flag.String("mock-result", mocksconfig.DefaultOutOfPackageDirectory, "where the mocks will be written")
	skipFolders := flag.String("skip-dirs", mocksconfig.DefaultOutOfPackageDirectory+","+config.VendorDir, "a comma separate list of directory to ignore, default value is the mocks and vendor directory")
	flag.Parse()

	err := run(project, diagEnable, mocksEnable, diagResult, diagGenerator, mockGenerator, mockResult, skipFolders)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var (
	errMissingMockGenerator    = errors.New("mock-generator is required")
	errMissingDiagramGenerator = errors.New("diag-generator is required")
	errMissingMockResult       = errors.New("mock-result is required")
)

func run(project *string, diagEnable, mocksEnable *bool, diagResult, diagGeneratorType, mockGeneratorType, mockResult, skipFolders *string) error {
	err := validateRequiredInput(diagEnable, mocksEnable, diagGeneratorType, mockGeneratorType, mockResult)
	if err != nil {
		return fmt.Errorf("validateRequiredInput: %w", err)
	}

	if project == nil || *project == "" {
		dir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("os.Getwd: %w", err)
		}
		project = &dir
	}

	as, err := getAst(project, skipFolders)
	if err != nil {
		return fmt.Errorf("getAst: %w", err)
	}

	err = runGenerators(as, project, diagEnable, mocksEnable, diagResult, diagGeneratorType, mockGeneratorType, mockResult)
	if err != nil {
		return fmt.Errorf("runGenerators: %w", err)
	}
	return nil
}

func runGenerators(as parse.AstSchema, project *string, diagEnable, mocksEnable *bool, diagResult, diagGeneratorType, mockGeneratorType, mockResult *string) (err error) {
	group, ctx := errgroup.WithContext(context.Background())
	if *diagEnable {
		group.Go(func() error {
			err = generateDiag(ctx, project, diagResult, diagGeneratorType, as)
			if err != nil {
				return fmt.Errorf("generateDiag: %w", err)
			}
			return nil
		})
	}

	if *mocksEnable {
		group.Go(func() error {
			err = generateMock(ctx, project, mockResult, mockGeneratorType, as)
			if err != nil {
				return fmt.Errorf("generateMock: %w", err)
			}
			return nil
		})
	}
	err = group.Wait()
	if err != nil {
		return fmt.Errorf("group wait: %w", err)
	}
	return nil
}

func generateDiag(ctx context.Context, project, diagResult, diagGeneratorType *string, as parse.AstSchema) error {
	diagGenerator, err := diagrams.GetGenerator(*diagGeneratorType)
	if err != nil {
		return fmt.Errorf("diagrams.GetGenerator:%w", err)
	}
	if diagResult == nil || *diagResult == "" {
		join := filepath.Join(*project, diagGenerator.GetDefaultResultFileName())
		diagResult = &join
	}

	w, closer, err := writer.GetWriter(diagResult)
	if err != nil {
		return fmt.Errorf("writer.GetWriter:%w", err)
	}
	defer closer()

	err = diagGenerator.GenerateFromSchema(ctx, w, as)
	if err != nil {
		return fmt.Errorf("diagGenerator.GenerateFromSchema:%w", err)
	}
	return nil
}

func generateMock(ctx context.Context, project, mockResult, mockGeneratorType *string, as parse.AstSchema) error {
	c := mocksconfig.Config{
		OutOfPackageMocksDirectory: filepath.Join(*project, *mockResult),
	}

	mockGenerator, err := mocks.GetGenerator(*mockGeneratorType, c)
	if err != nil {
		return fmt.Errorf("mocks.GetGenerator:%w", err)
	}
	err = mockGenerator.GenerateFromSchema(ctx, as)
	if err != nil {
		return fmt.Errorf("mockGenerator.GenerateFromSchema:%w", err)
	}
	return nil
}

func validateRequiredInput(diagEnable, mocksEnable *bool, diagGeneratorType, mockGeneratorType, mockResult *string) error {
	if *diagEnable {
		if diagGeneratorType == nil || *diagGeneratorType == "" {
			return errMissingDiagramGenerator
		}
		if mockGeneratorType == nil || *mockGeneratorType == "" {
			return errMissingMockGenerator
		}
	}
	if *mocksEnable {
		if mockResult == nil || *mockResult == "" {
			return errMissingMockResult
		}
	}
	return nil
}

func getAst(project, skipFolders *string) (parse.AstSchema, error) {
	var skipDirs []string
	if skipFolders != nil || *skipFolders != "" {
		skipDirs = strings.Split(*skipFolders, ",")
	}

	as, err := parse.Parse(*project, skipDirs)
	if err != nil {
		return parse.AstSchema{}, fmt.Errorf("parse.Parse:%w", err)
	}
	return as, nil
}
