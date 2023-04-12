package c4

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/emilien-puget/go-dependency-graph/pkg/parse"
)

const packageSeparator = "_"

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
	_, err = writer.WriteString("\nAddElementTag(\"external\", $bgColor=\"#8CDE42FF\")")
	if err != nil {
		return err
	}

	relations := ""
	externalRelations := make(map[string]string)
	for packageName, services := range s.Packages {
		rel, err := handlePackages(writer, packageName, services, externalRelations)
		if err != nil {
			return err
		}
		relations += rel
	}
	_, err = writer.WriteString(relations)
	if err != nil {
		return err
	}
	err = printExternalRelations(writer, externalRelations)
	if err != nil {
		return err
	}
	_, err = writer.WriteString("\n@enduml")
	if err != nil {
		return err
	}
	return nil
}

func printExternalRelations(writer *bufio.Writer, externalRelations map[string]string) error {
	for dep, rel := range externalRelations {
		_, err := fmt.Fprintf(writer, "Component(%s, %q, \"\", \"\", $tags=\"external\")\n", strings.ReplaceAll(dep, "/", packageSeparator), dep)
		if err != nil {
			return err
		}
		_, err = writer.WriteString(rel)
		if err != nil {
			return err
		}
	}
	return nil
}

func handlePackages(writer *bufio.Writer, packageName string, services parse.Dependencies, externalRelations map[string]string) (string, error) {
	packageUML := fmt.Sprintf("\n\nContainer_Boundary(%s, %q) {\n", packageName, packageName)
	relations := ""
	for serviceName, service := range services {
		fqdn := packageName + "." + serviceName
		packageUML += fmt.Sprintf("Component(%s, %q, \"\", %q)\n", fqdn, fqdn, service.Comment)

		for _, deps := range service.Deps {
			for _, d := range deps {
				if d.External {
					externalRelations[strings.ReplaceAll(d.PackageName, "/", packageSeparator)+"."+d.DependencyName] += getRelation(d, fqdn)
					continue
				}
				relations += getRelation(d, fqdn)
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

func getRelation(d parse.Dep, fqdn string) (relations string) {
	if len(d.Funcs) == 0 {
		return fmt.Sprintf("Rel(%s, %q, %q)\n", fqdn, strings.ReplaceAll(d.PackageName, "/", packageSeparator)+"."+d.DependencyName, d.PackageName+"."+d.DependencyName)
	}
	for _, fn := range d.Funcs {
		relations += fmt.Sprintf("Rel(%s, %q, %q)\n", fqdn, strings.ReplaceAll(d.PackageName, "/", packageSeparator)+"."+d.DependencyName, fn)
	}
	return relations
}
