package mermaid

import (
	"bufio"
	"bytes"
	"fmt"
	"sort"

	mymap "github.com/emilien-puget/go-dependency-graph/pkg/map"
	"github.com/emilien-puget/go-dependency-graph/pkg/parse"
)

const (
	packageSeparator = "/"
)

// GenerateClassFromSchema generates a class diagram for mermaid.
func GenerateClassFromSchema(writer *bufio.Writer, s parse.AstSchema) error {
	_, err := writer.WriteString("classDiagram\n")
	if err != nil {
		return err
	}

	var classBuf bytes.Buffer
	var relationBuf bytes.Buffer

	for _, k := range mymap.OrderedKeys(s.Graph.NodesByPackage) {
		err := handlePackages(&classBuf, &relationBuf, k, s.Graph.NodesByPackage[k], s.Graph)
		if err != nil {
			return err
		}
	}
	_, err = writer.Write(classBuf.Bytes())
	if err != nil {
		return err
	}
	_, err = writer.WriteString("\n")
	if err != nil {
		return err
	}
	_, err = writer.Write(relationBuf.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func handlePackages(classBuf, relationBuf *bytes.Buffer, packageName string, services []*parse.Node, graph *parse.Graph) error {
	_, err := fmt.Fprintf(classBuf, "\nnamespace %s {\n", packageName)
	if err != nil {
		return err
	}

	for i := range services {
		err := handleService(classBuf, relationBuf, packageName, services[i].StructName, services[i], graph)
		if err != nil {
			return err
		}
	}
	classBuf.WriteString("}")
	return nil
}

func handleService(classBuf, relationBuf *bytes.Buffer, packageName, serviceName string, service *parse.Node, graph *parse.Graph) error {
	serviceFqdn := packageName + packageSeparator + serviceName

	_, err := fmt.Fprintf(classBuf, "class `%s` {\n", serviceFqdn)
	if err != nil {
		return err
	}
	for _, method := range service.Methods {
		classBuf.WriteString(method)
		classBuf.WriteString("\n")
	}
	classBuf.WriteString("}\n\n")

	for _, d := range graph.GetAdjacenciesSortedByName(service) {
		err := handleDeps(d, relationBuf, serviceFqdn)
		if err != nil {
			return err
		}
	}
	return nil
}

func handleDeps(deps *parse.Adj, relationBuf *bytes.Buffer, serviceFqdn string) error {
	s := deps.Node.PackageName + packageSeparator + deps.Node.StructName
	if len(deps.Func) != 0 {
		sort.Slice(deps.Func, func(i, j int) bool {
			return deps.Func[i] < deps.Func[j]
		})
		for _, fn := range deps.Func {
			_, err := fmt.Fprintf(relationBuf, "`%s` ..> `%s`: %s\n", serviceFqdn, s, fn)
			if err != nil {
				return err
			}
		}
	} else {
		_, err := fmt.Fprintf(relationBuf, "`%s` ..> `%s`\n", serviceFqdn, s)
		if err != nil {
			return err
		}
	}

	return nil
}
