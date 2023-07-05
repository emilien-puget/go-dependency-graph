package mermaid

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/emilien-puget/go-dependency-graph/pkg/parse"
	"github.com/stretchr/testify/assert"
)

func TestGenerateMermaidClassFromSchema_fn(t *testing.T) {
	file := &bytes.Buffer{}
	buff := bufio.NewWriter(file)
	err := GenerateClassFromSchema(buff, parse.AstSchema{
		ModulePath: "testdata/fn",
		Packages: map[string]parse.Dependencies{
			"fn": {
				"A": {
					Comment: "",
					Deps: map[string][]parse.Dep{
						"b": {
							{
								PackageName:    "fn",
								DependencyName: "B",
								VarName:        "b",
								Funcs:          []string{"FuncA", "FuncB"},
							},
						},
						"d": {
							{
								PackageName:    "fn",
								DependencyName: "D",
								VarName:        "d",
								Funcs:          []string{"FuncA"},
							},
						},
					},
				},
				"B": {
					Comment: "",
					Deps: map[string][]parse.Dep{
						"c": {
							{
								PackageName:    "fn",
								DependencyName: "C",
								VarName:        "c",
								Funcs:          []string{"FuncA"},
							},
						},
					},
					Methods: []string{
						"FuncA()",
						"FuncB()",
					},
				},
				"C": {
					Comment: "",
					Deps:    map[string][]parse.Dep{},
					Methods: []string{
						"FuncA()",
					},
				},
				"D": {
					Comment: "",
					Deps: map[string][]parse.Dep{
						"a": {
							{
								PackageName:    "pa",
								DependencyName: "A",
								VarName:        "a",
								Funcs:          []string{"FuncA"},
							},
						},
					},
					Methods: []string{
						"FuncA()",
					},
				},
			},
			"pa": {
				"A": {
					Comment: "A pa struct.",
					Deps:    map[string][]parse.Dep{},
					Methods: []string{
						"FuncA(toto string) (titi int, err error)",
					},
				},
			},
		},
	})
	buff.Flush()
	assert.NoError(t, err)

	assert.Equal(t, "classDiagram\n\nnamespace fn {\nclass `fn/A` {\n}\n\nclass `fn/B` {\nFuncA()\nFuncB()\n}\n\nclass `fn/C` {\nFuncA()\n}\n\nclass `fn/D` {\nFuncA()\n}\n\n}\nnamespace pa {\nclass `pa/A` {\nFuncA(toto string) (titi int, err error)\n}\n\n}\n`fn/A` ..> `fn/B`: FuncA\n`fn/A` ..> `fn/B`: FuncB\n`fn/A` ..> `fn/D`: FuncA\n`fn/B` ..> `fn/C`: FuncA\n`fn/D` ..> `pa/A`: FuncA\n", file.String())
}

func TestGenerateMermaidClassFromSchema_ext_dep(t *testing.T) {
	file := &bytes.Buffer{}
	buff := bufio.NewWriter(file)
	err := GenerateClassFromSchema(buff, parse.AstSchema{
		ModulePath: "testdata/ext_dep",
		Packages: map[string]parse.Dependencies{
			"ext_dep": {
				"A": {
					Comment: "",
					Deps: map[string][]parse.Dep{
						"b": {
							{
								PackageName:    "net/http",
								DependencyName: "Client",
								VarName:        "b",
								External:       true,
							},
						},
					},
				},
			},
		},
	})
	buff.Flush()
	assert.NoError(t, err)

	assert.Equal(t, "classDiagram\n\nnamespace ext_dep {\nclass `ext_dep/A` {\n}\n\n}\n`ext_dep/A` ..> `net/http/Client`\n", file.String())
}

func TestGenerateMermaidClassFromSchema_inter(t *testing.T) {
	file := &bytes.Buffer{}
	buff := bufio.NewWriter(file)
	err := GenerateClassFromSchema(buff, parse.AstSchema{
		ModulePath: "testdata/inter",
		Packages: map[string]parse.Dependencies{
			"inter": {
				"A": {
					Comment: "",
					Deps: map[string][]parse.Dep{
						"b": {
							{
								PackageName:    "inter",
								DependencyName: "B",
								VarName:        "b",
								Funcs:          []string{"FuncA", "FuncB"},
							},
						},
						"d": {
							{
								PackageName:    "inter",
								DependencyName: "D",
								VarName:        "d",
								Funcs:          []string{"FuncA"},
							},
						},
					},
				},
				"B": {
					Comment: "",
					Deps: map[string][]parse.Dep{
						"c": {
							{
								PackageName:    "inter",
								DependencyName: "C",
								VarName:        "c",
								Funcs:          []string{"FuncA"},
							},
						},
					},
					Methods: []string{
						"FuncA()",
						"FuncB()",
					},
				},
				"C": {
					Comment: "",
					Deps:    map[string][]parse.Dep{},
					Methods: []string{
						"FuncA()",
					},
				},
				"D": {
					Comment: "",
					Deps: map[string][]parse.Dep{
						"a": {
							{
								PackageName:    "pa",
								DependencyName: "A",
								VarName:        "a",
								Funcs:          []string{"FuncA"},
							},
						},
					},
					Methods: []string{
						"FuncA()",
					},
				},
			},
			"pa": {
				"A": {
					Comment: "A pa struct.",
					Deps:    map[string][]parse.Dep{},
					Methods: []string{
						"FuncA()",
					},
				},
			},
		},
	})
	buff.Flush()
	assert.NoError(t, err)

	assert.Equal(t, "classDiagram\n\nnamespace inter {\nclass `inter/A` {\n}\n\nclass `inter/B` {\nFuncA()\nFuncB()\n}\n\nclass `inter/C` {\nFuncA()\n}\n\nclass `inter/D` {\nFuncA()\n}\n\n}\nnamespace pa {\nclass `pa/A` {\nFuncA()\n}\n\n}\n`inter/A` ..> `inter/B`: FuncA\n`inter/A` ..> `inter/B`: FuncB\n`inter/A` ..> `inter/D`: FuncA\n`inter/B` ..> `inter/C`: FuncA\n`inter/D` ..> `pa/A`: FuncA\n", file.String())
}

func TestGenerateMermaidClassFromSchema_package_name_mismatch(t *testing.T) {
	file := &bytes.Buffer{}
	buff := bufio.NewWriter(file)
	err := GenerateClassFromSchema(buff, parse.AstSchema{
		ModulePath: "testdata/package_name_mismatch",
		Packages: map[string]parse.Dependencies{
			"package_name_mismatch": {
				"A": {
					Comment: "",
					Deps: map[string][]parse.Dep{
						"encoder": {
							{
								PackageName:    "gopkg.in/yaml.v3",
								DependencyName: "Encoder",
								VarName:        "encoder",
								External:       true,
								Funcs:          nil,
							},
						},
					},
				},
			},
		},
	})
	buff.Flush()
	assert.NoError(t, err)

	assert.Equal(t, "classDiagram\n\nnamespace package_name_mismatch {\nclass `package_name_mismatch/A` {\n}\n\n}\n`package_name_mismatch/A` ..> `gopkg.in/yaml.v3/Encoder`\n", file.String())
}
