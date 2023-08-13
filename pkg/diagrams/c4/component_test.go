package c4

import (
	"bufio"
	"bytes"
	"go/token"
	"go/types"
	"testing"

	"github.com/emilien-puget/go-dependency-graph/pkg/parse"
	"github.com/emilien-puget/go-dependency-graph/pkg/parse/struct_decl"

	"github.com/stretchr/testify/assert"
)

func TestGenerateUmlFileFromSchema(t *testing.T) {
	file := &bytes.Buffer{}
	buff := bufio.NewWriter(file)

	graph := parse.NewGraph()
	fnA := &parse.Node{
		Name:        "fn.A",
		PackageName: "fn",
		StructName:  "A",
	}
	graph.AddNode(fnA)
	fnB := &parse.Node{
		Name:        "fn.B",
		PackageName: "fn",
		StructName:  "B",
		Methods: []struct_decl.Method{
			{
				TypFuc: types.NewFunc(token.NoPos, nil, "FuncA", &types.Signature{}),
			},
			{
				TypFuc: types.NewFunc(token.NoPos, nil, "FuncB", &types.Signature{}),
			},
		},
	}
	graph.AddNode(fnB)
	fnC := &parse.Node{
		Name:        "fn.C",
		PackageName: "fn",
		StructName:  "C",
		Methods: []struct_decl.Method{
			{
				TypFuc: types.NewFunc(token.NoPos, nil, "FuncA", &types.Signature{}),
			},
		},
	}
	graph.AddNode(fnC)
	fnD := &parse.Node{
		Name:        "fn.D",
		PackageName: "fn",
		StructName:  "D",
		Methods: []struct_decl.Method{
			{
				TypFuc: types.NewFunc(token.NoPos, nil, "FuncA", &types.Signature{}),
			},
		},
	}
	graph.AddNode(fnD)
	paA := &parse.Node{
		Name:        "pa.A",
		PackageName: "pa",
		StructName:  "A",
		Methods: []struct_decl.Method{
			{
				TypFuc: types.NewFunc(
					token.NoPos,
					nil,
					"FuncFoo",
					types.NewSignatureType(
						nil,
						nil,
						nil,
						types.NewTuple(types.NewParam(token.NoPos, nil, "foo", types.Typ[types.String])),
						types.NewTuple(
							types.NewParam(token.NoPos, nil, "bar", types.Typ[types.Int]),
							types.NewParam(token.NoPos, nil, "err", types.Universe.Lookup("error").Type()),
						),
						false,
					)),
			},
		},
		Doc: "A pa struct.",
	}
	graph.AddNode(paA)
	graph.AddEdge(fnA, &parse.Adj{Node: fnB, Func: []string{"FuncA", "FuncB"}})
	graph.AddEdge(fnA, &parse.Adj{Node: fnD, Func: []string{"FuncA"}})
	graph.AddEdge(fnB, &parse.Adj{Node: fnC, Func: []string{"FuncA"}})
	graph.AddEdge(fnD, &parse.Adj{Node: paA, Func: []string{"FuncFoo"}})

	err := GenerateComponentFromSchema(buff, parse.AstSchema{
		ModulePath: "testdata/fn",
		Graph:      graph,
	})
	buff.Flush()
	assert.NoError(t, err)

	assert.Equal(t, `@startuml
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Component.puml

title testdata/fn

Container_Boundary(fn, "fn") {
Component(fn_A, "fn.A", "", "")
Component(fn_B, "fn.B", "", "")
Component(fn_C, "fn.C", "", "")
Component(fn_D, "fn.D", "", "")

}


Container_Boundary(pa, "pa") {
Component(pa_A, "pa.A", "", "A pa struct.")

}
Rel(fn_A, "fn_B", "FuncA")
Rel(fn_A, "fn_B", "FuncB")
Rel(fn_A, "fn_D", "FuncA")
Rel(fn_B, "fn_C", "FuncA")
Rel(fn_D, "pa_A", "FuncFoo")

@enduml`, file.String())
}
