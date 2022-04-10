package c4

import (
	"bufio"
	"fmt"

	"github.com/emilien-puget/go-dependency-graph/pkg/parse"
)

func GenerateC4ComponentUmlFromSchema(writer *bufio.Writer, s parse.AstSchema) error {
	_, err := writer.Write([]byte("@startuml\n!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Component.puml\n"))
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte("\ntitle " + s.ModulePath))
	if err != nil {
		return err
	}

	relations := ""
	for packageName, services := range s.Packages {
		packageUML := fmt.Sprintf("\n\nContainer_Boundary(%s, \"%s\") {\n", packageName, packageName)

		for serviceName, service := range services {
			fqdn := packageName + "." + serviceName
			packageUML += fmt.Sprintf("Component(%s, \"%s\", \"\", \"%s\")\n", fqdn, fqdn, service.Comment)

			for _, deps := range service.Deps {
				for _, d := range deps {
					if len(d.Funcs) != 0 {
						for _, fn := range d.Funcs {
							relations += fmt.Sprintf("Rel(%s, %s, %s)\n", fqdn, d.PackageName+"."+d.ServiceName, fn)
						}
					} else {
						relations += fmt.Sprintf("Rel(%s, %s, %s)\n", fqdn, d.PackageName+"."+d.ServiceName, d.PackageName+"."+d.ServiceName)
					}
				}
			}
		}
		packageUML += "\n}\n"
		_, err := writer.Write([]byte(packageUML))
		if err != nil {
			return err
		}
	}
	_, err = writer.Write([]byte(relations))
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte("\n@enduml"))
	if err != nil {
		return err
	}
	return nil
}
