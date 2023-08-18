package mocks

import (
	"errors"
	"os"

	"github.com/emilien-puget/go-dependency-graph/pkg/mocks/config"
	"github.com/emilien-puget/go-dependency-graph/pkg/mocks/mockery"
	"github.com/emilien-puget/go-dependency-graph/pkg/parse"
)

type generateFromSchema func(c config.Config, s parse.AstSchema) error

const (
	GeneratorMockery = "mockery"
)

var generators = map[string]generateFromSchema{
	GeneratorMockery: mockery.GenerateFromSchema,
}

var errUnknownGenerator = errors.New("unknown generator")

func Generate(generator string, c config.Config, as parse.AstSchema) error {
	gen, ok := generators[generator]
	if !ok {
		return errUnknownGenerator
	}

	err := os.MkdirAll(c.OutOfPackageMocksDirectory, os.FileMode(0o755))
	if err != nil {
		return err
	}
	err = gen(
		c,
		as,
	)
	if err != nil {
		return err
	}
	return nil
}
