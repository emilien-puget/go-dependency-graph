package c4

import (
	"bufio"
	"fmt"

	"github.com/emilien-puget/go-dependency-graph/pkg/parse"
)

// GenerateComponentFromSchema generates a C4 plantuml component.
func GenerateComponentFromSchema(writer *bufio.Writer, s parse.AstSchema) error {
	_, err := writer.WriteString("@startuml\n!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Component.puml\n")
	if err != nil {
		return err
	}
	_, err = writer.WriteString("\ntitle " + s.ModulePath)
	if err != nil {
		return err
	}

	relations := ""
	for packageName, services := range s.Packages {
		rel, err := handlePackages(writer, packageName, services)
		if err != nil {
			return err
		}
		relations += rel
	}
	_, err = writer.WriteString(relations)
	if err != nil {
		return err
	}
	_, err = writer.WriteString("\n@enduml")
	if err != nil {
		return err
	}
	return nil
}

func handlePackages(writer *bufio.Writer, packageName string, services parse.Dependencies) (string, error) {
	packageUML := fmt.Sprintf("\n\nContainer_Boundary(%s, %q) {\n", packageName, packageName)
	relations := ""
	for serviceName, service := range services {
		fqdn := packageName + "." + serviceName
		packageUML += fmt.Sprintf("Component(%s, %q, \"\", %q)\n", fqdn, fqdn, service.Comment)

		for _, deps := range service.Deps {
			for _, d := range deps {
				if len(d.Funcs) != 0 {
					for _, fn := range d.Funcs {
						relations += fmt.Sprintf("Rel(%s, %s, %s)\n", fqdn, d.PackageName+"."+d.DependencyName, fn)
					}
				} else {
					relations += fmt.Sprintf("Rel(%s, %s, %s)\n", fqdn, d.PackageName+"."+d.DependencyName, d.PackageName+"."+d.DependencyName)
				}
			}
		}
	}
	packageUML += "\n}\n"
	_, err := writer.WriteString(packageUML)
	if err != nil {
		return "", err
	}
	return relations, nil
}
