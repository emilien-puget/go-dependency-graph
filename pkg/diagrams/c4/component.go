package c4

import (
	"bufio"
	"fmt"
	"strings"

	mymap "github.com/emilien-puget/go-dependency-graph/pkg/map"
	"github.com/emilien-puget/go-dependency-graph/pkg/parse"
)

const umlSeparator = "_"

type Generator struct {
	replacer *strings.Replacer
}

func NewGenerator() *Generator {
	return &Generator{
		replacer: strings.NewReplacer(".", umlSeparator, "-", umlSeparator, "/", umlSeparator),
	}
}

func (g Generator) GetDefaultResultFileName() string {
	return "diag.puml"
}

// GenerateFromSchema generates a C4 plantuml component.
func (g Generator) GenerateFromSchema(writer *bufio.Writer, s parse.AstSchema) error {
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
		rel, err := g.handlePackages(writer, packageName, s.Graph.NodesByPackage[packageName], externalRelations, s.Graph)
		if err != nil {
			return err
		}
		relations += rel
	}
	_, err = writer.WriteString(relations)
	if err != nil {
		return err
	}
	err = g.printExternalRelations(writer, externalRelations)
	if err != nil {
		return err
	}
	_, err = writer.WriteString("\n@enduml")
	if err != nil {
		return err
	}
	return nil
}

func (g Generator) getServiceID(service *parse.Node) string {
	return "\"" + g.replacer.Replace(service.PackageName+"."+service.StructName) + "\""
}

func (g Generator) getServiceLabel(service *parse.Node) string {
	return "\"" + service.PackageName + "." + service.StructName + "\""
}

func (g Generator) printExternalRelations(writer *bufio.Writer, externalRelations map[string]string) error {
	for dep, rel := range externalRelations {
		_, err := fmt.Fprintf(writer, "Component_Ext(%s, %q, \"\", \"\")\n", g.replacer.Replace(dep), dep)
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

func (g Generator) handlePackages(writer *bufio.Writer, packageName string, services []*parse.Node, externalRelations map[string]string, graph *parse.Graph) (string, error) {
	packageUML := fmt.Sprintf("\n\nContainer_Boundary(%s, %q) {\n", packageName, packageName)
	relations := ""
	for _, service := range services {
		serviceLabel := g.getServiceLabel(service)
		serviceID := g.getServiceID(service)
		packageUML += fmt.Sprintf("Component(%s, %s, \"\", %q)\n", serviceID, serviceLabel, service.Doc)

		for _, d := range graph.GetAdjacency(service) {
			if d.Node.External {
				externalRelations[strings.ReplaceAll(d.Node.PackageName, "/", umlSeparator)+"."+d.Node.StructName] += g.getRelation(serviceID, d)
				continue
			}
			relations += g.getRelation(serviceID, d)
		}
	}
	packageUML += "\n}\n"
	_, err := writer.WriteString(packageUML)
	if err != nil {
		return "", err
	}
	return relations, nil
}

func (g Generator) getRelation(sourceServiceID string, d *parse.Adj) (relations string) {
	if len(d.Func) == 0 {
		return fmt.Sprintf("Rel(%s, %s, %s)\n", sourceServiceID, g.getServiceID(d.Node), g.getServiceLabel(d.Node))
	}
	for _, fn := range d.Func {
		relations += fmt.Sprintf("Rel(%s, %s, %q)\n", sourceServiceID, g.getServiceID(d.Node), fn)
	}
	return relations
}
