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

func TestParse_package_alias(t *testing.T) {
	t.Parallel()
	parse, err := Parse("testdata/package_alias", nil)
	assert.NoError(t, err)

	graph := NewGraph()
	PaaA := &Node{
		Name:        "testdata/package_alias/pa/a.A",
		PackageName: "testdata/package_alias/pa/a",
		StructName:  "A",
		Doc:         "A pa struct.",
		External:    false,
	}
	graph.AddNode(PaaA)
	pbaA := &Node{
		Name:        "testdata/package_alias/pb/a.A",
		PackageName: "testdata/package_alias/pb/a",
		StructName:  "A",
		Doc:         "A pa struct.",
		External:    false,
	}
	graph.AddNode(pbaA)
	pbaB := &Node{
		Name:        "testdata/package_alias/pb/a.B",
		PackageName: "testdata/package_alias/pb/a",
		StructName:  "B",
		Doc:         "B pa struct.",
		External:    false,
	}
	graph.AddNode(pbaB)
	a := &Node{
		Name:        "testdata/package_alias.A",
		PackageName: "testdata/package_alias",
		StructName:  "A",
		External:    false,
	}
	graph.AddNode(a)
	graph.AddEdge(a, &Adj{Node: PaaA})
	graph.AddEdge(a, &Adj{Node: pbaA})
	graph.AddEdge(pbaB, &Adj{Node: pbaA})
	assertNodes(t, graph.GetNodesSortedByName(), parse.Graph.GetNodesSortedByName())
}

func TestParse_ext_dep(t *testing.T) {
	t.Parallel()
	parse, err := Parse("testdata/ext_dep", nil)
	assert.NoError(t, err)

	graph := NewGraph()
	extA := &Node{
		Name:        "testdata/ext_dep.A",
		PackageName: "testdata/ext_dep",
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

	require.NotNil(t, parse.Graph.GetNodeByName("testdata/ext_dep.A").ActualNamedType)
	require.NotNil(t, parse.Graph.GetNodeByName("testdata/ext_dep.A").P)

	require.Nil(t, parse.Graph.GetNodeByName("net/http.Client").ActualNamedType)
	require.Nil(t, parse.Graph.GetNodeByName("net/http.Client").P)
}

func TestParse_fn(t *testing.T) {
	t.Parallel()
	parse, err := Parse("testdata/fn", nil)
	assert.NoError(t, err)

	graph := NewGraph()
	fnA := &Node{
		Name:        "testdata/fn.A",
		PackageName: "testdata/fn",
		StructName:  "A",
	}
	graph.AddNode(fnA)
	fnB := &Node{
		Name:        "testdata/fn.B",
		PackageName: "testdata/fn",
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
		Name:        "testdata/fn.C",
		PackageName: "testdata/fn",
		StructName:  "C",
		Methods: []struct_decl.Method{
			{
				TypFuc: types.NewFunc(token.NoPos, nil, "FuncA", &types.Signature{}),
			},
		},
	}
	graph.AddNode(fnC)
	fnD := &Node{
		Name:        "testdata/fn.D",
		PackageName: "testdata/fn",
		StructName:  "D",
		Methods: []struct_decl.Method{
			{
				TypFuc: types.NewFunc(token.NoPos, nil, "FuncA", &types.Signature{}),
			},
		},
	}
	graph.AddNode(fnD)
	paA := &Node{
		Name:        "testdata/fn/pa.A",
		PackageName: "testdata/fn/pa",
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

	assertAdj(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("testdata/fn.A")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("testdata/fn.A")))
	assertAdj(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("testdata/fn.B")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("testdata/fn.B")))
	assertAdj(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("testdata/fn.C")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("testdata/fn.C")))
	assertAdj(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("testdata/fn.D")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("testdata/fn.D")))
	assertAdj(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("testdata/fn/pa.A")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("testdata/fn/pa.A")))
}

func TestParse_named_inter(t *testing.T) {
	t.Parallel()
	parse, err := Parse("testdata/named_inter", nil)
	assert.NoError(t, err)

	graph := NewGraph()
	interA := &Node{
		Name:        "testdata/named_inter.A",
		PackageName: "testdata/named_inter",
		StructName:  "A",
	}
	graph.AddNode(interA)
	interB := &Node{
		Name:        "testdata/named_inter.B",
		PackageName: "testdata/named_inter",
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
		Name:        "testdata/named_inter.C",
		PackageName: "testdata/named_inter",
		StructName:  "C",
		Methods: []struct_decl.Method{
			{
				TypFuc: types.NewFunc(token.NoPos, nil, "FuncA", &types.Signature{}),
			},
		},
	}
	graph.AddNode(interC)
	interD := &Node{
		Name:        "testdata/named_inter.D",
		PackageName: "testdata/named_inter",
		StructName:  "D",
		Methods: []struct_decl.Method{
			{
				TypFuc: types.NewFunc(token.NoPos, nil, "FuncA", &types.Signature{}),
			},
		},
	}
	graph.AddNode(interD)
	paA := &Node{
		Name:        "testdata/named_inter/pa.A",
		PackageName: "testdata/named_inter/pa",
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

	assertAdj(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("testdata/named_inter.A")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("testdata/named_inter.A")))
	assertAdj(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("testdata/named_inter.B")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("testdata/named_inter.B")))
	assertAdj(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("testdata/named_inter.C")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("testdata/named_inter.C")))
	assertAdj(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("testdata/named_inter.D")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("testdata/named_inter.D")))
	assertAdj(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("testdata/named_inter/pa.A")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("testdata/named_inter/pa.A")))
}

func TestParse_inter(t *testing.T) {
	t.Parallel()
	parse, err := Parse("testdata/inter", nil)
	assert.NoError(t, err)

	graph := NewGraph()
	interA := &Node{
		Name:        "testdata/inter.A",
		PackageName: "testdata/inter",
		StructName:  "A",
	}
	graph.AddNode(interA)
	interB := &Node{
		Name:        "testdata/inter.B",
		PackageName: "testdata/inter",
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
		Name:        "testdata/inter.C",
		PackageName: "testdata/inter",
		StructName:  "C",
		Methods: []struct_decl.Method{
			{
				TypFuc: types.NewFunc(token.NoPos, nil, "FuncA", &types.Signature{}),
			},
		},
	}
	graph.AddNode(interC)
	interD := &Node{
		Name:        "testdata/inter.D",
		PackageName: "testdata/inter",
		StructName:  "D",
		Methods: []struct_decl.Method{
			{
				TypFuc: types.NewFunc(token.NoPos, nil, "FuncA", &types.Signature{}),
			},
		},
	}
	graph.AddNode(interD)
	paA := &Node{
		Name:        "testdata/inter/pa.A",
		PackageName: "testdata/inter/pa",
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

	assertAdj(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("testdata/inter.A")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("testdata/inter.A")))
	assertAdj(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("testdata/inter.A")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("testdata/inter.A")))
	assertAdj(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("testdata/inter.B")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("testdata/inter.B")))
	assertAdj(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("testdata/inter.C")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("testdata/inter.C")))
	assertAdj(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("testdata/inter.D")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("testdata/inter.D")))
	assertAdj(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("testdata/inter/pa.A")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("testdata/inter/pa.A")))
}

func TestParse_wire_sample(t *testing.T) {
	t.Parallel()

	parse, err := Parse("testdata/wire_sample", nil)
	assert.NoError(t, err)

	graph := NewGraph()
	mainGreeter := &Node{
		Name:        "testdata/wire_sample.Greeter",
		PackageName: "testdata/wire_sample",
		StructName:  "Greeter",
		Doc:         "Greeter is the type charged with greeting guests.",
	}
	graph.AddNode(mainGreeter)
	mainEvent := &Node{
		Name:        "testdata/wire_sample.Event",
		PackageName: "testdata/wire_sample",
		StructName:  "Event",
		Doc:         "Event is a gathering with greeters.",
	}
	graph.AddNode(mainEvent)
	graph.AddEdge(mainEvent, &Adj{Node: mainGreeter})

	assertNodes(t, graph.GetNodesSortedByName(), parse.Graph.GetNodesSortedByName())
	assertAdj(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("testdata/wire_sample.Greeter")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("testdata/wire_sample.Greeter")))
	assertAdj(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("testdata/wire_sample.Event")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("testdata/wire_sample.Event")))
}

func TestParse_package_name_mismatch(t *testing.T) {
	t.Parallel()
	parse, err := Parse("testdata/package_name_mismatch", nil)
	assert.NoError(t, err)

	graph := NewGraph()
	mainGreeter := &Node{
		Name:        "testdata/package_name_mismatch.A",
		PackageName: "testdata/package_name_mismatch",
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
