package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse_ext_dep(t *testing.T) {
	t.Parallel()
	parse, err := Parse("testdata/ext_dep")
	assert.NoError(t, err)

	assert.Equal(t, AstSchema{
		ModulePath: "testdata/ext_dep",
		Packages: map[string]Dependencies{
			"ext_dep": {
				"A": {
					Comment: "",
					Deps: map[string][]Dep{
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
	}, parse)
}

func TestParse_fn(t *testing.T) {
	t.Parallel()
	parse, err := Parse("testdata/fn")
	assert.NoError(t, err)

	assert.Equal(t, AstSchema{
		ModulePath: "testdata/fn",
		Packages: map[string]Dependencies{
			"fn": {
				"A": {
					Comment: "",
					Deps: map[string][]Dep{
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
					Deps: map[string][]Dep{
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
					Deps:    map[string][]Dep{},
					Methods: []string{
						"FuncA()",
					},
				},
				"D": {
					Comment: "",
					Deps: map[string][]Dep{
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
					Deps:    map[string][]Dep{},
					Methods: []string{
						"FuncA(toto string) (titi int, err error)",
					},
				},
			},
		},
	}, parse)
}

func TestParse_inter(t *testing.T) {
	t.Parallel()
	parse, err := Parse("testdata/inter")
	assert.NoError(t, err)

	assert.Equal(t, AstSchema{
		ModulePath: "testdata/inter",
		Packages: map[string]Dependencies{
			"inter": {
				"A": {
					Comment: "",
					Deps: map[string][]Dep{
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
					Deps: map[string][]Dep{
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
					Deps:    map[string][]Dep{},
					Methods: []string{
						"FuncA()",
					},
				},
				"D": {
					Comment: "",
					Deps: map[string][]Dep{
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
					Deps:    map[string][]Dep{},
					Methods: []string{
						"FuncA()",
					},
				},
			},
		},
	}, parse)
}

func TestParse_wire_sample(t *testing.T) {
	t.Parallel()
	parse, err := Parse("testdata/wire_sample")
	assert.NoError(t, err)

	assert.Equal(t, AstSchema{
		ModulePath: "testdata/wire_sample",
		Packages: map[string]Dependencies{
			"main": {
				"Event": {
					Comment: "Event is a gathering with greeters.",
					Deps: map[string][]Dep{
						"g": {
							{
								PackageName:    "main",
								DependencyName: "Greeter",
								VarName:        "g",
								Funcs:          nil,
							},
						},
					},
					Methods: []string{
						"Start()",
					},
				},
				"Greeter": {
					Comment: "Greeter is the type charged with greeting guests.",
					Deps:    map[string][]Dep{},
					Methods: []string{
						"Greet() ( testdata/wire_sample.Message)",
					},
				},
			},
		},
	}, parse)
}

func TestParse_package_name_mismatch(t *testing.T) {
	t.Parallel()
	parse, err := Parse("testdata/package_name_mismatch")
	assert.NoError(t, err)

	assert.Equal(t, AstSchema{
		ModulePath: "testdata/package_name_mismatch",
		Packages: map[string]Dependencies{
			"package_name_mismatch": {
				"A": {
					Comment: "",
					Deps: map[string][]Dep{
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
	}, parse)
}
