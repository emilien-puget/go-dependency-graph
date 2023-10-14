package mermaid

import (
	"bufio"
	"bytes"
	"go/token"
	"go/types"
	"testing"

	"github.com/emilien-puget/go-dependency-graph/pkg/parse"
	"github.com/emilien-puget/go-dependency-graph/pkg/parse/struct_decl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateMermaidClassFromSchema_withParse(t *testing.T) {
	as, err := parse.Parse("../testdata/named_inter", nil)
	require.NoError(t, err)

	file := &bytes.Buffer{}
	buff := bufio.NewWriter(file)
	err = NewGenerator().GenerateFromSchema(buff, as)
	require.NoError(t, err)
	buff.Flush()

	assert.Equal(t, "classDiagram\n\nnamespace testdata_named_inter {\nclass `testdata/named_inter/A`\nclass `testdata/named_inter/B` {\nFuncA()\nFuncB()\n}\n\nclass `testdata/named_inter/C` {\nFuncA()\n}\n\nclass `testdata/named_inter/D` {\nFuncA()\n}\n\n}\nnamespace testdata_named_inter_pa {\nclass `testdata/named_inter/pa/A` {\nFuncFoo(foo string) (bar int, err error)\n}\n\n}\n`testdata/named_inter/A` ..> `testdata/named_inter/B`: FuncA\n`testdata/named_inter/A` ..> `testdata/named_inter/B`: FuncB\n`testdata/named_inter/A` ..> `testdata/named_inter/D`: FuncA\n`testdata/named_inter/B` ..> `testdata/named_inter/C`: FuncA\n`testdata/named_inter/D` ..> `testdata/named_inter/pa/A`: FuncFoo\n", file.String())
}

func TestGenerateMermaidClassFromSchema_fn(t *testing.T) {
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
	err := NewGenerator().GenerateFromSchema(nil, buff, parse.AstSchema{
		ModulePath: "testdata/fn",
		Graph:      graph,
	})
	buff.Flush()
	assert.NoError(t, err)

	assert.Equal(t, "classDiagram\n\nnamespace fn {\nclass `fn/A`\nclass `fn/B` {\nFuncA()\nFuncB()\n}\n\nclass `fn/C` {\nFuncA()\n}\n\nclass `fn/D` {\nFuncA()\n}\n\n}\nnamespace pa {\nclass `pa/A` {\nFuncFoo(foo string) (bar int, err error)\n}\n\n}\n`fn/A` ..> `fn/B`: FuncA\n`fn/A` ..> `fn/B`: FuncB\n`fn/A` ..> `fn/D`: FuncA\n`fn/B` ..> `fn/C`: FuncA\n`fn/D` ..> `pa/A`: FuncFoo\n", file.String())
}

func TestGenerateMermaidClassFromSchema_ext_dep(t *testing.T) {
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
	err := NewGenerator().GenerateFromSchema(nil, buff, parse.AstSchema{
		ModulePath: "testdata/ext_dep",
		Graph:      graph,
	})
	buff.Flush()
	assert.NoError(t, err)

	assert.Equal(t, "classDiagram\n\nnamespace ext_dep {\nclass `ext_dep/A`\n}\nnamespace net_http {\nclass `net/http/Client`\n}\n`ext_dep/A` ..> `net/http/Client`\n", file.String())
}

func TestGenerateMermaidClassFromSchema_inter(t *testing.T) {
	file := &bytes.Buffer{}
	buff := bufio.NewWriter(file)

	graph := parse.NewGraph()
	interA := &parse.Node{
		Name:        "inter.A",
		PackageName: "inter",
		StructName:  "A",
	}
	graph.AddNode(interA)
	interB := &parse.Node{
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
	interC := &parse.Node{
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
	interD := &parse.Node{
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
	graph.AddEdge(interA, &parse.Adj{Node: interB, Func: []string{"FuncA", "FuncB"}})
	graph.AddEdge(interA, &parse.Adj{Node: interD, Func: []string{"FuncA"}})
	graph.AddEdge(interB, &parse.Adj{Node: interC, Func: []string{"FuncA"}})
	graph.AddEdge(interD, &parse.Adj{Node: paA, Func: []string{"FuncFoo"}})
	err := NewGenerator().GenerateFromSchema(nil, buff, parse.AstSchema{
		ModulePath: "testdata/inter",
		Graph:      graph,
	})
	buff.Flush()
	assert.NoError(t, err)

	assert.Equal(t, "classDiagram\n\nnamespace inter {\nclass `inter/A`\nclass `inter/B` {\nFuncA()\nFuncB()\n}\n\nclass `inter/C` {\nFuncA()\n}\n\nclass `inter/D` {\nFuncA()\n}\n\n}\nnamespace pa {\nclass `pa/A` {\nFuncFoo(foo string) (bar int, err error)\n}\n\n}\n`inter/A` ..> `inter/B`: FuncA\n`inter/A` ..> `inter/B`: FuncB\n`inter/A` ..> `inter/D`: FuncA\n`inter/B` ..> `inter/C`: FuncA\n`inter/D` ..> `pa/A`: FuncFoo\n", file.String())
}

func TestGenerateMermaidClassFromSchema_package_name_mismatch(t *testing.T) {
	file := &bytes.Buffer{}
	buff := bufio.NewWriter(file)

	graph := parse.NewGraph()
	mainGreeter := &parse.Node{
		Name:        "package_name_mismatch.A",
		PackageName: "package_name_mismatch",
		StructName:  "A",
	}
	graph.AddNode(mainGreeter)
	mainEvent := &parse.Node{
		Name:        "gopkg.in/yaml.v3.Encoder",
		PackageName: "gopkg.in/yaml.v3",
		StructName:  "Encoder",
		External:    true,
	}
	graph.AddNode(mainEvent)
	graph.AddEdge(mainGreeter, &parse.Adj{Node: mainEvent})
	err := NewGenerator().GenerateFromSchema(nil, buff, parse.AstSchema{
		ModulePath: "testdata/package_name_mismatch",
		Graph:      graph,
	})
	buff.Flush()
	assert.NoError(t, err)

	assert.Equal(t, "classDiagram\n\nnamespace gopkg_in_yaml_v3 {\nclass `gopkg.in/yaml.v3/Encoder`\n}\nnamespace package_name_mismatch {\nclass `package_name_mismatch/A`\n}\n`package_name_mismatch/A` ..> `gopkg.in/yaml.v3/Encoder`\n", file.String())
}
