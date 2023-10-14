package c4

import (
	"bufio"
	"bytes"
	"context"
	"go/token"
	"go/types"
	"testing"

	"github.com/emilien-puget/go-dependency-graph/pkg/parse"
	"github.com/emilien-puget/go-dependency-graph/pkg/parse/struct_decl"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestGenerateUmlFileFromSchema_withParse(t *testing.T) {
	as, err := parse.Parse("../testdata/named_inter", nil)
	require.NoError(t, err)

	file := &bytes.Buffer{}
	buff := bufio.NewWriter(file)
	err = NewGenerator().GenerateFromSchema(context.Background(), buff, as)
	require.NoError(t, err)
	buff.Flush()

	assert.Equal(t, `@startuml
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Component.puml

title testdata/named_inter

Container_Boundary(testdata/named_inter, "testdata/named_inter") {
Component("A", "A", "", "")
Component("B", "B", "", "")
Component("C", "C", "", "")
Component("D", "D", "", "")

}


Container_Boundary(pa, "pa") {
Component("pa_A", "A", "", "A pa struct.")

}
Rel("A", "B", "FuncA")
Rel("A", "B", "FuncB")
Rel("A", "D", "FuncA")
Rel("B", "C", "FuncA")
Rel("D", "pa_A", "FuncFoo")

@enduml`, file.String())
}

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

	err := NewGenerator().GenerateFromSchema(context.Background(), buff, parse.AstSchema{
		ModulePath: "testdata/fn",
		Graph:      graph,
	})
	buff.Flush()
	assert.NoError(t, err)

	assert.Equal(t, `@startuml
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Component.puml

title testdata/fn

Container_Boundary(fn, "fn") {
Component("fn_A", "fn.A", "", "")
Component("fn_B", "fn.B", "", "")
Component("fn_C", "fn.C", "", "")
Component("fn_D", "fn.D", "", "")

}


Container_Boundary(pa, "pa") {
Component("pa_A", "pa.A", "", "A pa struct.")

}
Rel("fn_A", "fn_B", "FuncA")
Rel("fn_A", "fn_B", "FuncB")
Rel("fn_A", "fn_D", "FuncA")
Rel("fn_B", "fn_C", "FuncA")
Rel("fn_D", "pa_A", "FuncFoo")

@enduml`, file.String())
}

func TestGenerateUmlFileFromSchema_ext_dep(t *testing.T) {
	file := &bytes.Buffer{}
	buff := bufio.NewWriter(file)

	graph := parse.NewGraph()
	extA := &parse.Node{
		Name:        "ext_dep.A",
		PackageName: "ext_dep",
		StructName:  "A",
		External:    false,
	}
	graph.AddNode(extA)
	node := &parse.Node{
		Name:        "net/http.Client",
		PackageName: "net/http",
		StructName:  "Client",
		Methods:     nil,
		Doc:         "",
		External:    true,
	}
	graph.AddNode(node)
	graph.AddEdge(extA, &parse.Adj{
		Node: node,
		Func: nil,
	})
	err := NewGenerator().GenerateFromSchema(context.Background(), buff, parse.AstSchema{
		ModulePath: "testdata/ext_dep",
		Graph:      graph,
	})
	buff.Flush()
	assert.NoError(t, err)

	assert.Equal(t, `@startuml
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Component.puml

title testdata/ext_dep

Container_Boundary(ext_dep, "ext_dep") {
Component("ext_dep_A", "ext_dep.A", "", "")

}


Container_Boundary(net/http, "net/http") {
Component("net_http_Client", "net/http.Client", "", "")

}
Component_Ext(net_http_Client, "net_http.Client", "", "")
Rel("ext_dep_A", "net_http_Client", "net/http.Client")

@enduml`, file.String())
}
