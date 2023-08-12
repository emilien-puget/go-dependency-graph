package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse_ext_dep(t *testing.T) {
	t.Parallel()
	parse, err := Parse("testdata/ext_dep")
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
	assert.Equal(t, graph.GetNodesSortedByName(), parse.Graph.GetNodesSortedByName())
}

func TestParse_fn(t *testing.T) {
	t.Parallel()
	parse, err := Parse("testdata/fn")
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
		Methods: []string{
			"FuncA()",
			"FuncB()",
		},
	}
	graph.AddNode(fnB)
	fnC := &Node{
		Name:        "fn.C",
		PackageName: "fn",
		StructName:  "C",
		Methods: []string{
			"FuncA()",
		},
	}
	graph.AddNode(fnC)
	fnD := &Node{
		Name:        "fn.D",
		PackageName: "fn",
		StructName:  "D",
		Methods: []string{
			"FuncA()",
		},
	}
	graph.AddNode(fnD)
	paA := &Node{
		Name:        "pa.A",
		PackageName: "pa",
		StructName:  "A",
		Methods: []string{
			"FuncFoo(foo string) (bar int, err error)",
		},
		Doc: "A pa struct.",
	}
	graph.AddNode(paA)
	graph.AddEdge(fnA, &Adj{Node: fnB, Func: []string{"FuncA", "FuncB"}})
	graph.AddEdge(fnA, &Adj{Node: fnD, Func: []string{"FuncA"}})
	graph.AddEdge(fnB, &Adj{Node: fnC, Func: []string{"FuncA"}})
	graph.AddEdge(fnD, &Adj{Node: paA, Func: []string{"FuncFoo"}})

	assert.Equal(t, graph.GetNodesSortedByName(), parse.Graph.GetNodesSortedByName())

	assert.Equal(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("fn.A")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("fn.A")))
	assert.Equal(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("fn.B")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("fn.B")))
	assert.Equal(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("fn.C")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("fn.C")))
	assert.Equal(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("fn.D")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("fn.D")))
	assert.Equal(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("pa.A")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("pa.A")))
}

func TestParse_named_inter(t *testing.T) {
	t.Parallel()
	parse, err := Parse("testdata/named_inter")
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
		Methods: []string{
			"FuncA()",
			"FuncB()",
		},
	}
	graph.AddNode(interB)
	interC := &Node{
		Name:        "inter.C",
		PackageName: "inter",
		StructName:  "C",
		Methods: []string{
			"FuncA()",
		},
	}
	graph.AddNode(interC)
	interD := &Node{
		Name:        "inter.D",
		PackageName: "inter",
		StructName:  "D",
		Methods: []string{
			"FuncA()",
		},
	}
	graph.AddNode(interD)
	paA := &Node{
		Name:        "pa.A",
		PackageName: "pa",
		StructName:  "A",
		Methods: []string{
			"FuncFoo(foo string) (bar int, err error)",
		},
		Doc: "A pa struct.",
	}
	graph.AddNode(paA)
	graph.AddEdge(interA, &Adj{Node: interB, Func: []string{"FuncA", "FuncB"}})
	graph.AddEdge(interA, &Adj{Node: interD, Func: []string{"FuncA"}})
	graph.AddEdge(interB, &Adj{Node: interC, Func: []string{"FuncA"}})
	graph.AddEdge(interD, &Adj{Node: paA, Func: []string{"FuncFoo"}})

	assert.Equal(t, graph.GetNodesSortedByName(), parse.Graph.GetNodesSortedByName())

	assert.Equal(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("inter.A")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("inter.A")))
	assert.Equal(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("inter.B")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("inter.B")))
	assert.Equal(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("inter.C")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("inter.C")))
	assert.Equal(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("inter.D")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("inter.D")))
	assert.Equal(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("pa.A")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("pa.A")))
}

func TestParse_inter(t *testing.T) {
	t.Parallel()
	parse, err := Parse("testdata/inter")
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
		Methods: []string{
			"FuncA()",
			"FuncB()",
		},
	}
	graph.AddNode(interB)
	interC := &Node{
		Name:        "inter.C",
		PackageName: "inter",
		StructName:  "C",
		Methods: []string{
			"FuncA()",
		},
	}
	graph.AddNode(interC)
	interD := &Node{
		Name:        "inter.D",
		PackageName: "inter",
		StructName:  "D",
		Methods: []string{
			"FuncA()",
		},
	}
	graph.AddNode(interD)
	paA := &Node{
		Name:        "pa.A",
		PackageName: "pa",
		StructName:  "A",
		Methods: []string{
			"FuncFoo(foo string) (bar int, err error)",
		},
		Doc: "A pa struct.",
	}
	graph.AddNode(paA)
	graph.AddEdge(interA, &Adj{Node: interB, Func: []string{"FuncA", "FuncB"}})
	graph.AddEdge(interA, &Adj{Node: interD, Func: []string{"FuncA"}})
	graph.AddEdge(interB, &Adj{Node: interC, Func: []string{"FuncA"}})
	graph.AddEdge(interD, &Adj{Node: paA, Func: []string{"FuncFoo"}})

	assert.Equal(t, graph.GetNodesSortedByName(), parse.Graph.GetNodesSortedByName())

	assert.Equal(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("inter.A")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("inter.A")))
	assert.Equal(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("inter.B")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("inter.B")))
	assert.Equal(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("inter.C")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("inter.C")))
	assert.Equal(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("inter.D")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("inter.D")))
	assert.Equal(t, graph.GetAdjacenciesSortedByName(graph.GetNodeByName("pa.A")), parse.Graph.GetAdjacenciesSortedByName(parse.Graph.GetNodeByName("pa.A")))
}

func TestParse_wire_sample(t *testing.T) {
	t.Parallel()
	parse, err := Parse("testdata/wire_sample")
	assert.NoError(t, err)

	graph := NewGraph()
	mainGreeter := &Node{
		Name:        "main.Greeter",
		PackageName: "main",
		StructName:  "Greeter",
		Methods: []string{
			"Greet() ( testdata/wire_sample.Message)",
		},
		Doc: "Greeter is the type charged with greeting guests.",
	}
	graph.AddNode(mainGreeter)
	mainEvent := &Node{
		Name:        "main.Event",
		PackageName: "main",
		StructName:  "Event",
		Methods: []string{
			"Start()",
		},
		Doc: "Event is a gathering with greeters.",
	}
	graph.AddNode(mainEvent)
	graph.AddEdge(mainEvent, &Adj{Node: mainGreeter})

	assert.Equal(t, graph.GetNodesSortedByName(), parse.Graph.GetNodesSortedByName())
}

func TestParse_package_name_mismatch(t *testing.T) {
	t.Parallel()
	parse, err := Parse("testdata/package_name_mismatch")
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
	graph.AddEdge(mainGreeter, &Adj{Node: mainEvent})

	assert.Equal(t, graph.GetNodesSortedByName(), parse.Graph.GetNodesSortedByName())
}
