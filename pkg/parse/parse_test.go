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
								Funcs:          map[string]string{"FuncA": "", "FuncB": ""},
							},
						},
						"d": {
							{
								PackageName:    "fn",
								DependencyName: "D",
								VarName:        "d",
								Funcs:          map[string]string{"FuncA": ""},
							},
						},
						"s": {
							{
								PackageName:    "fn",
								DependencyName: "SomeFunc",
								VarName:        "s",
								Funcs:          map[string]string{"s": "(ctx context.Context, err error)"},
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
								Funcs:          map[string]string{"FuncA": ""},
							},
						},
					},
				},
				"C": {
					Comment: "",
					Deps:    map[string][]Dep{},
				},
				"D": {
					Comment: "",
					Deps: map[string][]Dep{
						"a": {
							{
								PackageName:    "pa",
								DependencyName: "A",
								VarName:        "a",
								Funcs:          map[string]string{"FuncA": ""},
							},
						},
					},
				},
			},
			"pa": {
				"A": {
					Comment: "A pa struct.",
					Deps:    map[string][]Dep{},
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
								Funcs:          map[string]string{"FuncA": "(ctx context.Context) (error)", "FuncB": "(context.Context) (err error)"},
							},
						},
						"d": {
							{
								PackageName:    "inter",
								DependencyName: "D",
								VarName:        "d",
								Funcs:          map[string]string{"FuncA": "()"},
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
								Funcs:          map[string]string{"FuncA": "()"},
							},
						},
					},
				},
				"C": {
					Comment: "",
					Deps:    map[string][]Dep{},
				},
				"D": {
					Comment: "",
					Deps: map[string][]Dep{
						"a": {
							{
								PackageName:    "pa",
								DependencyName: "A",
								VarName:        "a",
								Funcs:          map[string]string{"FuncA": "()"},
							},
						},
					},
				},
			},
			"pa": {
				"A": {
					Comment: "A pa struct.",
					Deps:    map[string][]Dep{},
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
				},
				"Greeter": {
					Comment: "Greeter is the type charged with greeting guests.",
					Deps: map[string][]Dep{
						"m": {
							{
								PackageName:    "main",
								DependencyName: "Message",
								VarName:        "m",
								Funcs:          nil,
							},
						},
					},
				},
			},
		},
	}, parse)
}
