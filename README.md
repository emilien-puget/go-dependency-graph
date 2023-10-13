# go-dependency-graph

A tool to build dependency graph for go programs based on dependency injection functions.

# Install

```
go install github.com/emilien-puget/go-dependency-graph/cmd/go-dependency-graph@latest
```
# Note regarding mermaid

GitHub doesn't support the namespace feature of the mermaidjs class diagram, you can use the
cli https://github.com/mermaid-js/mermaid-cli to generate a svg/png/pdf file.

# Example

## [Simple example with interfaces](./pkg/parse/testdata/inter)

### c4 plantuml component

```puml
@startuml
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Component.puml

title testdata/named_inter

Container_Boundary(testdata/named_inter, "testdata/named_inter") {
Component("A", "A", "", "")
Component("B", "B", "", "")
Component("C", "C", "", "")
Component("D", "D", "", "")

}


Container_Boundary(pa, "pa") {
Component("pa_A", "A", "", "A pa struct.")

}
Rel("A", "B", "FuncA")
Rel("A", "B", "FuncB")
Rel("A", "D", "FuncA")
Rel("B", "C", "FuncA")
Rel("D", "pa_A", "FuncFoo")

@enduml
```

[www.plantuml.com](http://www.plantuml.com/plantuml/uml/ROx1QWCX48RlFeNTKmBDUkcff-owvDH2AVIyJ5Pf11r5HqefVVTg0akIUj33z_qpy-yJGQJiB7imkDYiD3yHXVGiH8Il_jFGAHzpqd7nI1gfNxmJmGBMcLqYPSrHoAVTMqKVho_2GI8T2vgbTy5ZdGbrFoD3LdFIPGW818BJQZPOqen9ZmG6TPn7dr51_DwqWe-yQ-5kot_OUcxJ3Lq9dh_psrwxiQAnxMH5ikscYgOhntvPitU0uWFSTmemtzOQU02UAEQ5-ikwTsrBzxNV8UCo5DF0umsUxjENeFo7Qt0jKit1-tfwhr5bP_y0)

### mermaid class

```mermaid
classDiagram

  namespace testdata_named_inter {
    class `testdata/named_inter/A`
    class `testdata/named_inter/B` {
      FuncA()
      FuncB()
    }

    class `testdata/named_inter/C` {
      FuncA()
    }

    class `testdata/named_inter/D` {
      FuncA()
    }

  }
  namespace testdata_named_inter_pa {
    class `testdata/named_inter/pa/A` {
      FuncFoo(foo string) (bar int, err error)
    }

  }
  `testdata/named_inter/A` ..> `testdata/named_inter/B`: FuncA
  `testdata/named_inter/A` ..> `testdata/named_inter/B`: FuncB
  `testdata/named_inter/A` ..> `testdata/named_inter/D`: FuncA
  `testdata/named_inter/B` ..> `testdata/named_inter/C`: FuncA
  `testdata/named_inter/D` ..> `testdata/named_inter/pa/A`: FuncFoo


```

# Installation

`go install github.com/emilien-puget/go-dependency-graph@latest`

# How to Use

## Generator

`go-dependency-graph --project=<path to project> --result=<result file> --generator=<generator>`

Available generator are as follows

- `c4_plantuml_component`, the default value, more information about that format
  here : https://github.com/plantuml-stdlib/C4-PlantUML
- `mermaid_class`, a class diagram using mermaid
  syntax https://mermaid-js.github.io/mermaid/#/classDiagram?id=class-diagrams

### With the result written to a file

`go-dependency-graph --project=<path to project> --result=<result file>`

### With the result piped

`go-dependency-graph --project=<path to project> > <piped>`

