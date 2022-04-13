package json

import (
	"bufio"
	"encoding/json"

	"github.com/emilien-puget/go-dependency-graph/pkg/parse"
)

// GenerateFromSchema encode the AstSchema into a JSON.
func GenerateFromSchema(writer *bufio.Writer, s parse.AstSchema) error {
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(s)
	if err != nil {
		return err
	}

	return nil
}
