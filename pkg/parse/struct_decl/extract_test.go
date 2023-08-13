package struct_decl

import (
	"path/filepath"
	"testing"

	"github.com/emilien-puget/go-dependency-graph/pkg/parse/package_list"
	"github.com/stretchr/testify/require"
)

func TestExtractExtDep(t *testing.T) {
	t.Parallel()
	pkgs, err := package_list.GetPackagesToParse("../testdata/ext_dep")
	require.NoError(t, err)
	got := Extract(pkgs)

	require.Len(t, got, 1)

	require.Contains(t, got, "ext_dep")

	require.Contains(t, got["ext_dep"], "A")
	require.NotNil(t, got["ext_dep"]["A"].ActualNamedType)
	require.Len(t, got["ext_dep"]["A"].Methods, 0)
	require.Contains(t, got["ext_dep"]["A"].FilePath, filepath.Join("go-dependency-graph", "pkg", "parse", "testdata", "ext_dep", "a.go"))
}

func TestExtractInter(t *testing.T) {
	t.Parallel()
	pkgs, err := package_list.GetPackagesToParse("../testdata/inter")
	require.NoError(t, err)
	got := Extract(pkgs)

	require.Len(t, got, 2)

	require.Contains(t, got, "pa")
	require.Contains(t, got, "inter")

	require.Contains(t, got["pa"], "A")

	require.NotNil(t, got["pa"]["A"].ActualNamedType)
	require.Len(t, got["pa"]["A"].Methods, 1)
	require.Contains(t, got["pa"]["A"].FilePath, filepath.Join("go-dependency-graph", "pkg", "parse", "testdata", "inter", "pa", "a.go"))
	require.Equal(t, "FuncFoo(foo string) (bar int, err error)", got["pa"]["A"].Methods[0].String())

	require.Contains(t, got["inter"], "A")
	require.NotNil(t, got["inter"]["A"].ActualNamedType)
	require.Len(t, got["inter"]["A"].Methods, 0)
	require.Contains(t, got["inter"]["A"].FilePath, filepath.Join("go-dependency-graph", "pkg", "parse", "testdata", "inter", "a.go"))

	require.Contains(t, got["inter"], "B")
	require.NotNil(t, got["inter"]["B"].Methods)
	require.NotNil(t, got["inter"]["B"].ActualNamedType)
	require.Len(t, got["inter"]["B"].Methods, 2)
	require.Contains(t, got["inter"]["B"].FilePath, filepath.Join("go-dependency-graph", "pkg", "parse", "testdata", "inter", "b.go"))
	require.Equal(t, "FuncA()", got["inter"]["B"].Methods[0].String())
	require.Equal(t, "FuncB()", got["inter"]["B"].Methods[1].String())

	require.Contains(t, got["inter"], "C")
	require.NotNil(t, got["inter"]["C"].Methods)
	require.NotNil(t, got["inter"]["C"].ActualNamedType)
	require.Len(t, got["inter"]["C"].Methods, 1)
	require.Contains(t, got["inter"]["C"].FilePath, filepath.Join("go-dependency-graph", "pkg", "parse", "testdata", "inter", "c.go"))
	require.Equal(t, "FuncA()", got["inter"]["C"].Methods[0].String())

	require.Contains(t, got["inter"], "D")
	require.NotNil(t, got["inter"]["D"].Methods)
	require.NotNil(t, got["inter"]["D"].ActualNamedType)
	require.Len(t, got["inter"]["D"].Methods, 1)
	require.Contains(t, got["inter"]["D"].FilePath, filepath.Join("go-dependency-graph", "pkg", "parse", "testdata", "inter", "d.go"))
	require.Equal(t, "FuncA()", got["inter"]["D"].Methods[0].String())
}

func TestExtractNamedInter(t *testing.T) {
	t.Parallel()
	pkgs, err := package_list.GetPackagesToParse("../testdata/named_inter")
	require.NoError(t, err)
	got := Extract(pkgs)

	require.Len(t, got, 2)

	require.Contains(t, got, "pa")
	require.Contains(t, got, "inter")

	require.Contains(t, got["pa"], "A")

	require.NotNil(t, got["pa"]["A"].ActualNamedType)
	require.Len(t, got["pa"]["A"].Methods, 1)
	require.Equal(t, "FuncFoo(foo string) (bar int, err error)", got["pa"]["A"].Methods[0].String())

	require.Contains(t, got["inter"], "A")
	require.NotNil(t, got["inter"]["A"].ActualNamedType)
	require.Len(t, got["inter"]["A"].Methods, 0)

	require.Contains(t, got["inter"], "B")
	require.NotNil(t, got["inter"]["B"].Methods)
	require.NotNil(t, got["inter"]["B"].ActualNamedType)
	require.Len(t, got["inter"]["B"].Methods, 2)
	require.Equal(t, "FuncA()", got["inter"]["B"].Methods[0].String())
	require.Equal(t, "FuncB()", got["inter"]["B"].Methods[1].String())

	require.Contains(t, got["inter"], "C")
	require.NotNil(t, got["inter"]["C"].Methods)
	require.NotNil(t, got["inter"]["C"].ActualNamedType)
	require.Len(t, got["inter"]["C"].Methods, 1)
	require.Equal(t, "FuncA()", got["inter"]["C"].Methods[0].String())

	require.Contains(t, got["inter"], "D")
	require.NotNil(t, got["inter"]["D"].Methods)
	require.NotNil(t, got["inter"]["D"].ActualNamedType)
	require.Len(t, got["inter"]["D"].Methods, 1)
	require.Equal(t, "FuncA()", got["inter"]["D"].Methods[0].String())
}

func TestExtractFunc(t *testing.T) {
	t.Parallel()
	pkgs, err := package_list.GetPackagesToParse("../testdata/fn")
	require.NoError(t, err)
	got := Extract(pkgs)

	require.Len(t, got, 2)

	require.Contains(t, got, "pa")
	require.Contains(t, got, "fn")

	require.Contains(t, got["pa"], "A")

	require.NotNil(t, got["pa"]["A"].ActualNamedType)
	require.Len(t, got["pa"]["A"].Methods, 1)
	require.Equal(t, "FuncFoo(foo string) (bar int, err error)", got["pa"]["A"].Methods[0].String())

	require.Contains(t, got["fn"], "A")
	require.NotNil(t, got["fn"]["A"].ActualNamedType)
	require.Len(t, got["fn"]["A"].Methods, 0)

	require.Contains(t, got["fn"], "B")
	require.NotNil(t, got["fn"]["B"].Methods)
	require.NotNil(t, got["fn"]["B"].ActualNamedType)
	require.Len(t, got["fn"]["B"].Methods, 2)
	require.Equal(t, "FuncA()", got["fn"]["B"].Methods[0].String())
	require.Equal(t, "FuncB()", got["fn"]["B"].Methods[1].String())

	require.Contains(t, got["fn"], "C")
	require.NotNil(t, got["fn"]["C"].Methods)
	require.NotNil(t, got["fn"]["C"].ActualNamedType)
	require.Len(t, got["fn"]["C"].Methods, 1)
	require.Equal(t, "FuncA()", got["fn"]["C"].Methods[0].String())

	require.Contains(t, got["fn"], "D")
	require.NotNil(t, got["fn"]["D"].Methods)
	require.NotNil(t, got["fn"]["D"].ActualNamedType)
	require.Len(t, got["fn"]["D"].Methods, 1)
	require.Equal(t, "FuncA()", got["fn"]["D"].Methods[0].String())
}
