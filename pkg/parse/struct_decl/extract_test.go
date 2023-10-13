package struct_decl

import (
	"path/filepath"
	"testing"

	"github.com/emilien-puget/go-dependency-graph/pkg/config"
	"github.com/emilien-puget/go-dependency-graph/pkg/parse/package_list"
	"github.com/stretchr/testify/require"
)

func TestExtractExtDep(t *testing.T) {
	t.Parallel()
	pkgs, err := package_list.GetPackagesToParse("../testdata/ext_dep", []string{config.VendorDir})
	require.NoError(t, err)
	got := Extract(pkgs)

	require.Len(t, got, 1)

	packageName := "testdata/ext_dep"
	require.Contains(t, got, packageName)

	require.Contains(t, got[packageName], "A")
	require.NotNil(t, got[packageName]["A"].ActualNamedType)
	require.Len(t, got[packageName]["A"].Methods, 0)
	require.Contains(t, got[packageName]["A"].FilePath, filepath.Join("go-dependency-graph", "pkg", "parse", "testdata", "ext_dep", "a.go"))
}

func TestExtractPackageAlias(t *testing.T) {
	t.Parallel()
	pkgs, err := package_list.GetPackagesToParse("../testdata/package_alias", []string{config.VendorDir})
	require.NoError(t, err)
	got := Extract(pkgs)

	require.Len(t, got, 3)
	require.Contains(t, got, "testdata/package_alias/pa/a")
	require.Contains(t, got, "testdata/package_alias/pb/a")
	require.Contains(t, got, "testdata/package_alias")
}

func TestExtractInter(t *testing.T) {
	t.Parallel()
	pkgs, err := package_list.GetPackagesToParse("../testdata/inter", nil)
	require.NoError(t, err)
	got := Extract(pkgs)

	require.Len(t, got, 2)

	packageNamePa := "testdata/inter/pa"
	require.Contains(t, got, packageNamePa)
	packageNameInter := "testdata/inter"
	require.Contains(t, got, packageNameInter)

	require.Contains(t, got[packageNamePa], "A")

	require.NotNil(t, got[packageNamePa]["A"].ActualNamedType)
	require.Len(t, got[packageNamePa]["A"].Methods, 1)
	require.Contains(t, got[packageNamePa]["A"].FilePath, filepath.Join("go-dependency-graph", "pkg", "parse", "testdata", "inter", "pa", "a.go"))
	require.Equal(t, "FuncFoo(foo string) (bar int, err error)", got[packageNamePa]["A"].Methods[0].String())

	require.Contains(t, got[packageNameInter], "A")
	require.NotNil(t, got[packageNameInter]["A"].ActualNamedType)
	require.Len(t, got[packageNameInter]["A"].Methods, 0)
	require.Contains(t, got[packageNameInter]["A"].FilePath, filepath.Join("go-dependency-graph", "pkg", "parse", "testdata", "inter", "a.go"))

	require.Contains(t, got[packageNameInter], "B")
	require.NotNil(t, got[packageNameInter]["B"].Methods)
	require.NotNil(t, got[packageNameInter]["B"].ActualNamedType)
	require.Len(t, got[packageNameInter]["B"].Methods, 2)
	require.Contains(t, got[packageNameInter]["B"].FilePath, filepath.Join("go-dependency-graph", "pkg", "parse", "testdata", "inter", "b.go"))
	require.Equal(t, "FuncA()", got[packageNameInter]["B"].Methods[0].String())
	require.Equal(t, "FuncB()", got[packageNameInter]["B"].Methods[1].String())

	require.Contains(t, got[packageNameInter], "C")
	require.NotNil(t, got[packageNameInter]["C"].Methods)
	require.NotNil(t, got[packageNameInter]["C"].ActualNamedType)
	require.Len(t, got[packageNameInter]["C"].Methods, 1)
	require.Contains(t, got[packageNameInter]["C"].FilePath, filepath.Join("go-dependency-graph", "pkg", "parse", "testdata", "inter", "c.go"))
	require.Equal(t, "FuncA()", got[packageNameInter]["C"].Methods[0].String())

	require.Contains(t, got[packageNameInter], "D")
	require.NotNil(t, got[packageNameInter]["D"].Methods)
	require.NotNil(t, got[packageNameInter]["D"].ActualNamedType)
	require.Len(t, got[packageNameInter]["D"].Methods, 1)
	require.Contains(t, got[packageNameInter]["D"].FilePath, filepath.Join("go-dependency-graph", "pkg", "parse", "testdata", "inter", "d.go"))
	require.Equal(t, "FuncA()", got[packageNameInter]["D"].Methods[0].String())
}

func TestExtractNamedInter(t *testing.T) {
	t.Parallel()
	pkgs, err := package_list.GetPackagesToParse("../testdata/named_inter", nil)
	require.NoError(t, err)
	got := Extract(pkgs)

	require.Len(t, got, 2)

	packageNamePa := "testdata/named_inter/pa"
	require.Contains(t, got, packageNamePa)
	packageNameInter := "testdata/named_inter"
	require.Contains(t, got, packageNameInter)

	require.Contains(t, got[packageNamePa], "A")

	require.NotNil(t, got[packageNamePa]["A"].ActualNamedType)
	require.Len(t, got[packageNamePa]["A"].Methods, 1)
	require.Equal(t, "FuncFoo(foo string) (bar int, err error)", got[packageNamePa]["A"].Methods[0].String())

	require.Contains(t, got[packageNameInter], "A")
	require.NotNil(t, got[packageNameInter]["A"].ActualNamedType)
	require.Len(t, got[packageNameInter]["A"].Methods, 0)

	require.Contains(t, got[packageNameInter], "B")
	require.NotNil(t, got[packageNameInter]["B"].Methods)
	require.NotNil(t, got[packageNameInter]["B"].ActualNamedType)
	require.Len(t, got[packageNameInter]["B"].Methods, 2)
	require.Equal(t, "FuncA()", got[packageNameInter]["B"].Methods[0].String())
	require.Equal(t, "FuncB()", got[packageNameInter]["B"].Methods[1].String())

	require.Contains(t, got[packageNameInter], "C")
	require.NotNil(t, got[packageNameInter]["C"].Methods)
	require.NotNil(t, got[packageNameInter]["C"].ActualNamedType)
	require.Len(t, got[packageNameInter]["C"].Methods, 1)
	require.Equal(t, "FuncA()", got[packageNameInter]["C"].Methods[0].String())

	require.Contains(t, got[packageNameInter], "D")
	require.NotNil(t, got[packageNameInter]["D"].Methods)
	require.NotNil(t, got[packageNameInter]["D"].ActualNamedType)
	require.Len(t, got[packageNameInter]["D"].Methods, 1)
	require.Equal(t, "FuncA()", got[packageNameInter]["D"].Methods[0].String())
}

func TestExtractFunc(t *testing.T) {
	t.Parallel()
	pkgs, err := package_list.GetPackagesToParse("../testdata/fn", nil)
	require.NoError(t, err)
	got := Extract(pkgs)

	require.Len(t, got, 2)

	packageNamePa := "testdata/fn/pa"
	require.Contains(t, got, packageNamePa)
	packageNameFn := "testdata/fn"
	require.Contains(t, got, packageNameFn)

	require.Contains(t, got[packageNamePa], "A")

	require.NotNil(t, got[packageNamePa]["A"].ActualNamedType)
	require.Len(t, got[packageNamePa]["A"].Methods, 1)
	require.Equal(t, "FuncFoo(foo string) (bar int, err error)", got[packageNamePa]["A"].Methods[0].String())

	require.Contains(t, got[packageNameFn], "A")
	require.NotNil(t, got[packageNameFn]["A"].ActualNamedType)
	require.Len(t, got[packageNameFn]["A"].Methods, 0)

	require.Contains(t, got[packageNameFn], "B")
	require.NotNil(t, got[packageNameFn]["B"].Methods)
	require.NotNil(t, got[packageNameFn]["B"].ActualNamedType)
	require.Len(t, got[packageNameFn]["B"].Methods, 2)
	require.Equal(t, "FuncA()", got[packageNameFn]["B"].Methods[0].String())
	require.Equal(t, "FuncB()", got[packageNameFn]["B"].Methods[1].String())

	require.Contains(t, got[packageNameFn], "C")
	require.NotNil(t, got[packageNameFn]["C"].Methods)
	require.NotNil(t, got[packageNameFn]["C"].ActualNamedType)
	require.Len(t, got[packageNameFn]["C"].Methods, 1)
	require.Equal(t, "FuncA()", got[packageNameFn]["C"].Methods[0].String())

	require.Contains(t, got[packageNameFn], "D")
	require.NotNil(t, got[packageNameFn]["D"].Methods)
	require.NotNil(t, got[packageNameFn]["D"].ActualNamedType)
	require.Len(t, got[packageNameFn]["D"].Methods, 1)
	require.Equal(t, "FuncA()", got[packageNameFn]["D"].Methods[0].String())
}
