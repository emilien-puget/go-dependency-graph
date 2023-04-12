package mermaid

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/emilien-puget/go-dependency-graph/pkg/parse"
)

const packageSeparator = "_"

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
		fqdn := strings.ReplaceAll(packageName, "/", packageSeparator) + packageSeparator + serviceName

		// TODO : handle packages when mermaid support them https://github.com/mermaid-js/mermaid/issues/1052
		for _, deps := range service.Deps {
			for _, d := range deps {
				if len(d.Funcs) != 0 {
					for _, fn := range d.Funcs {
						_, err := fmt.Fprintf(writer, "%s ..> %s: %s\n", fqdn, strings.ReplaceAll(d.PackageName, "/", packageSeparator)+packageSeparator+d.DependencyName, fn)
						if err != nil {
							return err
						}
					}
				} else {
					_, err := fmt.Fprintf(writer, "%s ..> %s\n", fqdn, strings.ReplaceAll(d.PackageName, "/", packageSeparator)+packageSeparator+d.DependencyName)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
