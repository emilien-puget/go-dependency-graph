package c4

import (
	"bufio"
	"fmt"
	"strings"

	mymap "github.com/emilien-puget/go-dependency-graph/pkg/map"
	"github.com/emilien-puget/go-dependency-graph/pkg/parse"
)

const packageSeparator = "_"

var replacer *strings.Replacer

// GenerateComponentFromSchema generates a C4 plantuml component.
func GenerateComponentFromSchema(writer *bufio.Writer, s parse.AstSchema) error {
	replacer = strings.NewReplacer(".", "_", "-", "_", "/", "_")

	_, err := writer.WriteString("@startuml\n!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Component.puml\n")
	if err != nil {
		return err
	}
	_, err = writer.WriteString("\ntitle " + s.ModulePath)
	if err != nil {
		return err
	}

	relations := ""
	externalRelations := make(map[string]string)

	for _, packageName := range mymap.OrderedKeys(s.Graph.NodesByPackage) {
		rel, err := handlePackages(writer, packageName, s.Graph.NodesByPackage[packageName], externalRelations, s.Graph)
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

func getServiceID(service *parse.Node) string {
	return "\"" + replacer.Replace(service.PackageName+"."+service.StructName) + "\""
}

func getServiceLabel(service *parse.Node) string {
	return "\"" + service.PackageName + "." + service.StructName + "\""
}

func printExternalRelations(writer *bufio.Writer, externalRelations map[string]string) error {
	for dep, rel := range externalRelations {
		_, err := fmt.Fprintf(writer, "Component_Ext(%s, %q, \"\", \"\")\n", replacer.Replace(dep), dep)
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

func handlePackages(writer *bufio.Writer, packageName string, services []*parse.Node, externalRelations map[string]string, graph *parse.Graph) (string, error) {
	packageUML := fmt.Sprintf("\n\nContainer_Boundary(%s, %q) {\n", packageName, packageName)
	relations := ""
	for _, service := range services {
		serviceLabel := getServiceLabel(service)
		serviceID := getServiceID(service)
		packageUML += fmt.Sprintf("Component(%s, %s, \"\", %q)\n", serviceID, serviceLabel, service.Doc)

		for _, d := range graph.GetAdjacency(service) {
			if d.Node.External {
				externalRelations[strings.ReplaceAll(d.Node.PackageName, "/", packageSeparator)+"."+d.Node.StructName] += getRelation(serviceID, d)
				continue
			}
			relations += getRelation(serviceID, d)
		}
	}
	packageUML += "\n}\n"
	_, err := writer.WriteString(packageUML)
	if err != nil {
		return "", err
	}
	return relations, nil
}

func getRelation(sourceServiceID string, d *parse.Adj) (relations string) {
	if len(d.Func) == 0 {
		return fmt.Sprintf("Rel(%s, %s, %q)\n", sourceServiceID, getServiceID(d.Node), getServiceLabel(d.Node))
	}
	for _, fn := range d.Func {
		relations += fmt.Sprintf("Rel(%s, %s, %q)\n", sourceServiceID, getServiceID(d.Node), fn)
	}
	return relations
}
