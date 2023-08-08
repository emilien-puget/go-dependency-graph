package mermaid

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/emilien-puget/go-dependency-graph/pkg/parse"
	"github.com/stretchr/testify/assert"
)

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
		Methods: []string{
			"FuncA()",
			"FuncB()",
		},
	}
	graph.AddNode(fnB)
	fnC := &parse.Node{
		Name:        "fn.C",
		PackageName: "fn",
		StructName:  "C",
		Methods: []string{
			"FuncA()",
		},
	}
	graph.AddNode(fnC)
	fnD := &parse.Node{
		Name:        "fn.D",
		PackageName: "fn",
		StructName:  "D",
		Methods: []string{
			"FuncA()",
		},
	}
	graph.AddNode(fnD)
	paA := &parse.Node{
		Name:        "pa.A",
		PackageName: "pa",
		StructName:  "A",
		Methods: []string{
			"FuncFoo(foo string) (bar int, err error)",
		},
		Doc: "A pa struct.",
	}
	graph.AddNode(paA)
	graph.AddEdge(fnA, &parse.Adj{Node: fnB, Func: []string{"FuncA", "FuncB"}})
	graph.AddEdge(fnA, &parse.Adj{Node: fnD, Func: []string{"FuncA"}})
	graph.AddEdge(fnB, &parse.Adj{Node: fnC, Func: []string{"FuncA"}})
	graph.AddEdge(fnD, &parse.Adj{Node: paA, Func: []string{"FuncFoo"}})
	err := GenerateClassFromSchema(buff, parse.AstSchema{
		ModulePath: "testdata/fn",
		Graph:      graph,
	})
	buff.Flush()
	assert.NoError(t, err)

	assert.Equal(t, "classDiagram\n\nnamespace fn {\nclass `fn/A` {\n}\n\nclass `fn/B` {\nFuncA()\nFuncB()\n}\n\nclass `fn/C` {\nFuncA()\n}\n\nclass `fn/D` {\nFuncA()\n}\n\n}\nnamespace pa {\nclass `pa/A` {\nFuncFoo(foo string) (bar int, err error)\n}\n\n}\n`fn/A` ..> `fn/B`: FuncA\n`fn/A` ..> `fn/B`: FuncB\n`fn/A` ..> `fn/D`: FuncA\n`fn/B` ..> `fn/C`: FuncA\n`fn/D` ..> `pa/A`: FuncFoo\n", file.String())
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
	err := GenerateClassFromSchema(buff, parse.AstSchema{
		ModulePath: "testdata/ext_dep",
		Graph:      graph,
	})
	buff.Flush()
	assert.NoError(t, err)

	assert.Equal(t, "classDiagram\n\nnamespace ext_dep {\nclass `ext_dep/A` {\n}\n\n}\nnamespace net/http {\nclass `net/http/Client` {\n}\n\n}\n`ext_dep/A` ..> `net/http/Client`\n", file.String())
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
		Methods: []string{
			"FuncA()",
			"FuncB()",
		},
	}
	graph.AddNode(interB)
	interC := &parse.Node{
		Name:        "inter.C",
		PackageName: "inter",
		StructName:  "C",
		Methods: []string{
			"FuncA()",
		},
	}
	graph.AddNode(interC)
	interD := &parse.Node{
		Name:        "inter.D",
		PackageName: "inter",
		StructName:  "D",
		Methods: []string{
			"FuncA()",
		},
	}
	graph.AddNode(interD)
	paA := &parse.Node{
		Name:        "pa.A",
		PackageName: "pa",
		StructName:  "A",
		Methods: []string{
			"FuncFoo(foo string) (bar int, err error)",
		},
		Doc: "A pa struct.",
	}
	graph.AddNode(paA)
	graph.AddEdge(interA, &parse.Adj{Node: interB, Func: []string{"FuncA", "FuncB"}})
	graph.AddEdge(interA, &parse.Adj{Node: interD, Func: []string{"FuncA"}})
	graph.AddEdge(interB, &parse.Adj{Node: interC, Func: []string{"FuncA"}})
	graph.AddEdge(interD, &parse.Adj{Node: paA, Func: []string{"FuncFoo"}})
	err := GenerateClassFromSchema(buff, parse.AstSchema{
		ModulePath: "testdata/inter",
		Graph:      graph,
	})
	buff.Flush()
	assert.NoError(t, err)

	assert.Equal(t, "classDiagram\n\nnamespace inter {\nclass `inter/A` {\n}\n\nclass `inter/B` {\nFuncA()\nFuncB()\n}\n\nclass `inter/C` {\nFuncA()\n}\n\nclass `inter/D` {\nFuncA()\n}\n\n}\nnamespace pa {\nclass `pa/A` {\nFuncFoo(foo string) (bar int, err error)\n}\n\n}\n`inter/A` ..> `inter/B`: FuncA\n`inter/A` ..> `inter/B`: FuncB\n`inter/A` ..> `inter/D`: FuncA\n`inter/B` ..> `inter/C`: FuncA\n`inter/D` ..> `pa/A`: FuncFoo\n", file.String())
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
	err := GenerateClassFromSchema(buff, parse.AstSchema{
		ModulePath: "testdata/package_name_mismatch",
		Graph:      graph,
	})
	buff.Flush()
	assert.NoError(t, err)

	assert.Equal(t, "classDiagram\n\nnamespace gopkg.in/yaml.v3 {\nclass `gopkg.in/yaml.v3/Encoder` {\n}\n\n}\nnamespace package_name_mismatch {\nclass `package_name_mismatch/A` {\n}\n\n}\n`package_name_mismatch/A` ..> `gopkg.in/yaml.v3/Encoder`\n", file.String())
}
