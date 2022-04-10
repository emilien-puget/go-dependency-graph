package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
								Funcs:          []string{"FuncA"},
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
								Funcs:          []string{"FuncA"},
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
