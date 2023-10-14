package mocks

import (
	"context"
	"errors"

	"github.com/emilien-puget/go-dependency-graph/pkg/mocks/config"
	"github.com/emilien-puget/go-dependency-graph/pkg/mocks/mockery"
	"github.com/emilien-puget/go-dependency-graph/pkg/parse"
)

const (
	GeneratorMockery = "mockery"
)

type Generator interface {
	GenerateFromSchema(ctx context.Context, as parse.AstSchema) error
}

var errUnknownGenerator = errors.New("unknown generator")

func GetGenerator(generator string, c config.Config) (Generator, error) {
	switch generator {
	case GeneratorMockery:
		return mockery.NewGenerator(c.OutOfPackageMocksDirectory), nil
	default:
		return nil, errUnknownGenerator
	}
}
