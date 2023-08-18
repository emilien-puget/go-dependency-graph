package parse

import (
	"fmt"
	"go/token"
	"go/types"
	"testing"

	"github.com/emilien-puget/go-dependency-graph/pkg/parse/struct_decl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_error(t *testing.T) {
	tests := map[string]struct {
		pathDir string
		wantErr assert.ErrorAssertionFunc
	}{
		"no_go_mod": {
			pathDir: "testdata/no_go_mod",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrGoModNotFound)
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := Parse(tt.pathDir, nil)
			if !tt.wantErr(t, err, fmt.Sprintf("Parse(%v)", tt.pathDir)) {
				return
			}
		})
	}
}

func TestParse_ext_dep(t *testing.T) {
	t.Parallel()
	parse, err := Parse("testdata/ext_dep", nil)
	assert.NoError(t, err)

	graph := NewGraph()
	extA := &Node{
		Name:        "ext_dep.A",
		PackageName: "ext_dep",
		StructName:  "A",
		External:    false,
	}
	graph.AddNode(extA)
	node := &Node{
		Name:        "net/http.Client",
		PackageName: "net/http",
		StructName:  "Client",
		Methods:     nil,
		Doc:         "",
		External:    true,
	}
	graph.AddNode(node)
	graph.AddEdge(extA, &Adj{
		Node: node,
		Func: nil,
	})
	assertNodes(t, graph.GetNodesSortedByName(), parse.Graph.GetNodesSortedByName())

	require.NotNil(t, parse.Graph.GetNodeByName("ext_dep.A").ActualNamedType)
	require.NotNil(t, parse.Graph.GetNodeByName("ext_dep.A").P)

	require.Nil(t, parse.Graph.GetNodeByName("net/http.Client").ActualNamedType)
	require.Nil(t, parse.Graph.GetNodeByName("net/http.Client").P)
}

func TestParse_fn(t *testing.T) {
	t.Parallel()
	parse, err := Parse("testdata/fn", nil)
	assert.NoError(t, err)

	graph := NewGraph()
	fnA := &Node{
		Name:        "fn.A",
		PackageName: "fn",
		StructName:  "A",
	}
	graph.AddNode(fnA)
	fnB := &Node{
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
	fnC := &Node{
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
	fnD := &Node{
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
	paA := &Node{
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
					),
				),
			},
		},
		Doc: "A pa struct.",
	}
	graph.AddNode(paA)
	graph.AddEdge(fnA, &Adj{Node: fnB, Func: []string{"FuncA", "FuncB"}})
	graph.AddEdge(fnA, &Adj{Node: fnD, Func: []string{"FuncA"}})
	graph.AddEdge(fnB, &Adj{Node: fnC, Func: []string{"FuncA"}})
	graph.AddEdge(fnD, &Adj{Node: paA, Func: []string{"FuncFoo"}})

	assertNodes(t, graph.GetNodesSortedByName(), parse.Graph.GetNodesSortedByName())

	assertAdj(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("fn.A")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("fn.A")))
	assertAdj(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("fn.B")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("fn.B")))
	assertAdj(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("fn.C")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("fn.C")))
	assertAdj(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("fn.D")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("fn.D")))
	assertAdj(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("pa.A")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("pa.A")))
}

func TestParse_named_inter(t *testing.T) {
	t.Parallel()
	parse, err := Parse("testdata/named_inter", nil)
	assert.NoError(t, err)

	graph := NewGraph()
	interA := &Node{
		Name:        "inter.A",
		PackageName: "inter",
		StructName:  "A",
	}
	graph.AddNode(interA)
	interB := &Node{
		Name:        "inter.B",
		PackageName: "inter",
		StructName:  "B",
		Methods: []struct_decl.Method{
			{
				TypFuc: types.NewFunc(token.NoPos, nil, "FuncA", &types.Signature{}),
			}, {
				TypFuc: types.NewFunc(token.NoPos, nil, "FuncB", &types.Signature{}),
			},
		},
	}
	graph.AddNode(interB)
	interC := &Node{
		Name:        "inter.C",
		PackageName: "inter",
		StructName:  "C",
		Methods: []struct_decl.Method{
			{
				TypFuc: types.NewFunc(token.NoPos, nil, "FuncA", &types.Signature{}),
			},
		},
	}
	graph.AddNode(interC)
	interD := &Node{
		Name:        "inter.D",
		PackageName: "inter",
		StructName:  "D",
		Methods: []struct_decl.Method{
			{
				TypFuc: types.NewFunc(token.NoPos, nil, "FuncA", &types.Signature{}),
			},
		},
	}
	graph.AddNode(interD)
	paA := &Node{
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
					),
				),
			},
		},
		Doc: "A pa struct.",
	}
	graph.AddNode(paA)
	graph.AddEdge(interA, &Adj{Node: interB, Func: []string{"FuncA", "FuncB"}})
	graph.AddEdge(interA, &Adj{Node: interD, Func: []string{"FuncA"}})
	graph.AddEdge(interB, &Adj{Node: interC, Func: []string{"FuncA"}})
	graph.AddEdge(interD, &Adj{Node: paA, Func: []string{"FuncFoo"}})

	assertNodes(t, graph.GetNodesSortedByName(), parse.Graph.GetNodesSortedByName())

	assertAdj(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("inter.A")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("inter.A")))
	assertAdj(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("inter.B")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("inter.B")))
	assertAdj(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("inter.C")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("inter.C")))
	assertAdj(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("inter.D")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("inter.D")))
	assertAdj(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("pa.A")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("pa.A")))
}

func TestParse_inter(t *testing.T) {
	t.Parallel()
	parse, err := Parse("testdata/inter", nil)
	assert.NoError(t, err)

	graph := NewGraph()
	interA := &Node{
		Name:        "inter.A",
		PackageName: "inter",
		StructName:  "A",
	}
	graph.AddNode(interA)
	interB := &Node{
		Name:        "inter.B",
		PackageName: "inter",
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
	graph.AddNode(interB)
	interC := &Node{
		Name:        "inter.C",
		PackageName: "inter",
		StructName:  "C",
		Methods: []struct_decl.Method{
			{
				TypFuc: types.NewFunc(token.NoPos, nil, "FuncA", &types.Signature{}),
			},
		},
	}
	graph.AddNode(interC)
	interD := &Node{
		Name:        "inter.D",
		PackageName: "inter",
		StructName:  "D",
		Methods: []struct_decl.Method{
			{
				TypFuc: types.NewFunc(token.NoPos, nil, "FuncA", &types.Signature{}),
			},
		},
	}
	graph.AddNode(interD)
	paA := &Node{
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
					),
				),
			},
		},
		Doc: "A pa struct.",
	}
	graph.AddNode(paA)
	graph.AddEdge(interA, &Adj{Node: interB, Func: []string{"FuncA", "FuncB"}})
	graph.AddEdge(interA, &Adj{Node: interD, Func: []string{"FuncA"}})
	graph.AddEdge(interB, &Adj{Node: interC, Func: []string{"FuncA"}})
	graph.AddEdge(interD, &Adj{Node: paA, Func: []string{"FuncFoo"}})

	assertNodes(t, graph.GetNodesSortedByName(), parse.Graph.GetNodesSortedByName())

	assertAdj(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("inter.A")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("inter.A")))
	assertAdj(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("inter.A")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("inter.A")))
	assertAdj(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("inter.B")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("inter.B")))
	assertAdj(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("inter.C")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("inter.C")))
	assertAdj(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("inter.D")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("inter.D")))
	assertAdj(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("pa.A")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("pa.A")))
}

func TestParse_wire_sample(t *testing.T) {
	t.Parallel()

	parse, err := Parse("testdata/wire_sample", nil)
	assert.NoError(t, err)

	graph := NewGraph()
	mainGreeter := &Node{
		Name:        "main.Greeter",
		PackageName: "main",
		StructName:  "Greeter",
		Doc:         "Greeter is the type charged with greeting guests.",
	}
	graph.AddNode(mainGreeter)
	mainEvent := &Node{
		Name:        "main.Event",
		PackageName: "main",
		StructName:  "Event",
		Doc:         "Event is a gathering with greeters.",
	}
	graph.AddNode(mainEvent)
	graph.AddEdge(mainEvent, &Adj{Node: mainGreeter})

	assertNodes(t, graph.GetNodesSortedByName(), parse.Graph.GetNodesSortedByName())
	assertAdj(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("main.Greeter")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("main.Greeter")))
	assertAdj(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("main.Event")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("main.Event")))
}

func TestParse_package_name_mismatch(t *testing.T) {
	t.Parallel()
	parse, err := Parse("testdata/package_name_mismatch", nil)
	assert.NoError(t, err)

	graph := NewGraph()
	mainGreeter := &Node{
		Name:        "package_name_mismatch.A",
		PackageName: "package_name_mismatch",
		StructName:  "A",
	}
	graph.AddNode(mainGreeter)
	mainEvent := &Node{
		Name:        "gopkg.in/yaml.v3.Encoder",
		PackageName: "gopkg.in/yaml.v3",
		StructName:  "Encoder",
		External:    true,
	}
	graph.AddNode(mainEvent)
	graph.AddEdge(mainGreeter, &Adj{Node: mainEvent, Func: []string{"Encode"}})

	assertNodes(t, graph.GetNodesSortedByName(), parse.Graph.GetNodesSortedByName())
	assertAdj(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("package_name_mismatch.A")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("package_name_mismatch.A")))
}

func assertAdj(t *testing.T, expectedAdj, gotAdj []*Adj) {
	for i := range expectedAdj {
		require.Equal(t, expectedAdj[i].Func, gotAdj[i].Func)
		assertNode(t, expectedAdj[i].Node, gotAdj[i].Node)

	}
}

func assertNodes(t *testing.T, expectedNodes, gotNodes []*Node) {
	require.Equal(t, len(expectedNodes), len(gotNodes))
	for i := range gotNodes {
		assertNode(t, expectedNodes[i], gotNodes[i])
	}
}

func assertNode(t *testing.T, expected, got *Node) {
	require.Equal(t, expected.Name, got.Name)
	require.Equal(t, expected.StructName, got.StructName)
	require.Equal(t, expected.PackageName, got.PackageName)
	require.Equal(t, expected.Doc, got.Doc)

	for i2 := range expected.Methods {
		require.Equal(t, expected.Methods[i2].String(), got.Methods[i2].String())
	}
}
