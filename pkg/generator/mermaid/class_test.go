package mermaid

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/emilien-puget/go-dependency-graph/pkg/parse"
	"github.com/stretchr/testify/assert"
)

func TestGenerateMermaidClassFromSchema(t *testing.T) {
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
				},
				"C": {
					Comment: "",
					Deps:    map[string][]parse.Dep{},
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
				},
			},
			"pa": {
				"A": {
					Comment: "A pa struct.",
					Deps:    map[string][]parse.Dep{},
				},
			},
		},
	})
	buff.Flush()
	assert.NoError(t, err)

	// TODO : properly test the output
	// assert.Equal(t, "classDiagram\nfn_A <.. fn_B: FuncA\nfn_A <.. fn_B: FuncB\nfn_A <.. fn_D: FuncA\nfn_B <.. fn_C: FuncA\nfn_D <.. pa_A: FuncA\n", file.String())
}
