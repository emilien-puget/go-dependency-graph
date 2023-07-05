package mermaid

import (
	"bufio"
	"bytes"
	"fmt"
	"sort"

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

	keys := make([]string, 0, len(s.Packages))
	for k := range s.Packages {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		err := handlePackages(&classBuf, &relationBuf, k, s.Packages[k])
		if err != nil {
			return err
		}
	}
	writer.Write(classBuf.Bytes())
	writer.WriteString("\n")
	writer.Write(relationBuf.Bytes())
	return nil
}

func handlePackages(classBuf, relationBuf *bytes.Buffer, packageName string, services parse.Dependencies) error {
	_, err := fmt.Fprintf(classBuf, "\nnamespace %s {\n", packageName)
	if err != nil {
		return err
	}

	keys := make([]string, 0, len(services))
	for k := range services {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, serviceName := range keys {
		service := services[serviceName]
		err := handleService(classBuf, relationBuf, packageName, serviceName, &service)
		if err != nil {
			return err
		}
	}
	classBuf.WriteString("}")
	return nil
}

func handleService(classBuf, relationBuf *bytes.Buffer, packageName, serviceName string, service *parse.Dependency) error {
	fqdn := packageName + packageSeparator + serviceName

	_, err := fmt.Fprintf(classBuf, "class `%s` {\n", fqdn)
	if err != nil {
		return err
	}
	for _, method := range service.Methods {
		classBuf.WriteString(method)
		classBuf.WriteString("\n")
	}
	classBuf.WriteString("}\n\n")

	for _, deps := range service.Deps {
		sort.Slice(deps, func(i, j int) bool {
			return deps[i].DependencyName < deps[j].DependencyName
		})
		for _, d := range deps {
			s := d.PackageName + packageSeparator + d.DependencyName
			if len(d.Funcs) != 0 {
				sort.Slice(d.Funcs, func(i, j int) bool {
					return d.Funcs[i] < d.Funcs[j]
				})
				for _, fn := range d.Funcs {
					_, err := fmt.Fprintf(relationBuf, "`%s` ..> `%s`: %s\n", fqdn, s, fn)
					if err != nil {
						return err
					}
				}
			} else {
				_, err := fmt.Fprintf(relationBuf, "`%s` ..> `%s`\n", fqdn, s)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
