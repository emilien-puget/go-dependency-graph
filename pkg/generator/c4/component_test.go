package c4

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/emilien-puget/go-dependency-graph/pkg/parse"

	"github.com/stretchr/testify/assert"
)

func TestGenerateUmlFileFromSchema(t *testing.T) {
	file := &bytes.Buffer{}
	buff := bufio.NewWriter(file)
	err := GenerateC4ComponentUmlFromSchema(buff, parse.AstSchema{
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
	// assert.Equal(t, "@startuml\n!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Component.puml\n\ntitle testdata/fn\n\nContainer_Boundary(fn, \"fn\") {\nComponent(fn.A, \"fn.A\", \"\", \"\")\nComponent(fn.B, \"fn.B\", \"\", \"\")\nComponent(fn.C, \"fn.C\", \"\", \"\")\nComponent(fn.D, \"fn.D\", \"\", \"\")\n\n}\n\n\nContainer_Boundary(pa, \"pa\") {\nComponent(pa.A, \"pa.A\", \"\", \"A pa struct.\")\n\n}\nRel(fn.A, fn.B, FuncA)\nRel(fn.A, fn.B, FuncB)\nRel(fn.A, fn.D, FuncA)\nRel(fn.B, fn.C, FuncA)\nRel(fn.D, pa.A, FuncA)\n\n@enduml", file.String())
}
