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

	for _, k := range orderedKeys(s.Packages) {
		err := handlePackages(&classBuf, &relationBuf, k, s.Packages[k])
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

func handlePackages(classBuf, relationBuf *bytes.Buffer, packageName string, services parse.Dependencies) error {
	_, err := fmt.Fprintf(classBuf, "\nnamespace %s {\n", packageName)
	if err != nil {
		return err
	}

	for _, serviceName := range orderedKeys(services) {
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

	for _, d := range orderedKeys(service.Deps) {
		err := handleDeps(service.Deps[d], relationBuf, serviceFqdn)
		if err != nil {
			return err
		}
	}
	return nil
}

func handleDeps(deps []parse.Dep, relationBuf *bytes.Buffer, serviceFqdn string) error {
	sort.SliceStable(deps, func(i, j int) bool {
		return deps[i].PackageName+deps[i].DependencyName < deps[j].PackageName+deps[j].DependencyName
	})
	for _, d := range deps {
		s := d.PackageName + packageSeparator + d.DependencyName
		if len(d.Funcs) != 0 {
			sort.Slice(d.Funcs, func(i, j int) bool {
				return d.Funcs[i] < d.Funcs[j]
			})
			for _, fn := range d.Funcs {
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
	}
	return nil
}

func orderedKeys[v any](tab map[string]v) []string {
	keys := make([]string, 0, len(tab))
	for k := range tab {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
