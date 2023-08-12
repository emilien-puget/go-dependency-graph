package struct_decl

import (
	"testing"

	"github.com/emilien-puget/go-dependency-graph/pkg/parse/package_list"
	"github.com/stretchr/testify/require"
)

func TestSearch(t *testing.T) {
	t.Parallel()
	pkgs, err := package_list.GetPackagesToParse("../testdata/fn")
	require.NoError(t, err)

	var docs []map[string]string
	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			doc := GetStructDoc(file, file.Name.Name)
			docs = append(docs, doc)
		}
	}

	require.Equal(
		t,
		[]map[string]string{
			{"pa.A": "// A pa struct."},
			{"fn.A": ""},
			{"fn.B": ""},
			{"fn.C": ""},
			{"fn.D": ""},
			{},
		},
		docs,
	)
}
