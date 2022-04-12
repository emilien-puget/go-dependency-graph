package mermaid

import (
	"bufio"
	"fmt"

	"github.com/emilien-puget/go-dependency-graph/pkg/parse"
)

// GenerateClassFromSchema generates a class diagram for mermaid.
func GenerateClassFromSchema(writer *bufio.Writer, s parse.AstSchema) error {
	_, err := writer.WriteString("classDiagram\n")
	if err != nil {
		return err
	}

	for packageName, dependencies := range s.Packages {
		err := handlePackages(writer, packageName, dependencies)
		if err != nil {
			return err
		}
	}
	return nil
}

func handlePackages(writer *bufio.Writer, packageName string, services parse.Dependencies) error {
	for serviceName, service := range services {
		fqdn := packageName + "_" + serviceName

		// TODO : handle packages when mermaid support them https://github.com/mermaid-js/mermaid/issues/1052
		for _, deps := range service.Deps {
			for _, d := range deps {
				if len(d.Funcs) != 0 {
					for _, fn := range d.Funcs {
						_, err := writer.WriteString(fmt.Sprintf("%s <.. %s: %s\n", fqdn, d.PackageName+"_"+d.DependencyName, fn))
						if err != nil {
							return err
						}
					}
				} else {
					_, err := writer.WriteString(fmt.Sprintf("%s <.. %s\n", fqdn, d.PackageName+"_"+d.DependencyName))
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
