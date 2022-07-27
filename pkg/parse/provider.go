package parse

import (
	"go/ast"
	"go/types"
)

func searchProvider(funcdecl *ast.FuncDecl, structs map[string]structDecl, packageName string, imports map[string]Import, typesInfo *types.Info) (name string, deps map[string][]Dep) {
	if funcdecl.Name.Name[:3] != "New" {
		return name, deps
	}
	name = searchDependencyName(funcdecl)
	s, ok := structs[packageName+"."+name]
	if !ok {
		return "", nil
	}

	deps = searchDependencies(funcdecl, packageName, imports, typesInfo)

	searchDependenciesAssignment(funcdecl, deps, s)
	return name, deps
}

// searchDependencyName search the created dependency as the first variable returned.
func searchDependencyName(funcdecl *ast.FuncDecl) string {
	results := funcdecl.Type.Results
	if results == nil {
		return ""
	}
	switch t := funcdecl.Type.Results.List[0].Type.(type) { // get the type of the dependency.
	case *ast.StarExpr: // dependency returned as a pointer.
		ident, ok := t.X.(*ast.Ident)
		if !ok {
			return ""
		}
		return ident.Name
	case *ast.Ident: // dependency returned as a value.
		return t.Name
	}
	return ""
}

// searchDependencies returns the dependency found in the provider type declaration.
func searchDependencies(funcdecl *ast.FuncDecl, name string, imports map[string]Import, info *types.Info) (deps map[string][]Dep) {
	deps = map[string][]Dep{}
	for _, param := range funcdecl.Type.Params.List {
		if !checkDepsMethods(info.TypeOf(param.Type)) { // ignore dependencies without methods
			continue
		}
		packageName, serviceName := getDepID(param.Type)
		if serviceName == "" {
			continue
		}
		imp := imports[packageName]
		external := imp.External
		if packageName == "" {
			packageName = name
		}
		varName := ""
		for _, name := range param.Names {
			varName = name.String()
		}
		if external {
			packageName = imp.Path
		}
		deps[varName] = append(deps[varName], Dep{
			VarName:        varName,
			PackageName:    packageName,
			DependencyName: serviceName,
			External:       external,
		})
	}
	return deps
}

func checkDepsMethods(t types.Type) bool {
	ptrType, ok := t.(*types.Pointer)
	if ok {
		t = ptrType.Elem()
	}
	namedType, ok := t.(*types.Named)
	if !ok {
		return false
	}

	if namedType.NumMethods() > 0 {
		return true
	}
	return false
}

func getDepID(dep ast.Expr) (packageName, serviceName string) {
	if depStar, ok := dep.(*ast.StarExpr); ok {
		dep = depStar.X
	}

	switch p := dep.(type) {
	case *ast.SelectorExpr:
		ident, ok := p.X.(*ast.Ident)
		if !ok {
			return "", ""
		}
		return ident.Name, p.Sel.Name
	case *ast.Ident:
		return "", p.Name
	}

	return "", ""
}

// searchDependenciesAssignment parse the provider function to search for a return.
// this return is then parsed to look for injected functions to complete the deps found in the previous step.
func searchDependenciesAssignment(funcdecl *ast.FuncDecl, deps map[string][]Dep, s structDecl) {
	if funcdecl.Body != nil {
		for _, stmt := range funcdecl.Body.List {
			retStmt, ok := stmt.(*ast.ReturnStmt)
			if !ok {
				continue
			}
			uExpr, ok := retStmt.Results[0].(*ast.UnaryExpr)
			if !ok {
				continue
			}
			cpLit, ok := uExpr.X.(*ast.CompositeLit)
			if !ok {
				continue
			}
			for _, elt := range cpLit.Elts {
				kvExpr, ok := elt.(*ast.KeyValueExpr)
				if !ok {
					continue
				}
				setDepsFunc(kvExpr, deps, s)
			}
		}
	}
}

func setDepsFunc(kvExpr *ast.KeyValueExpr, deps map[string][]Dep, s structDecl) {
	switch value := kvExpr.Value.(type) {
	case *ast.SelectorExpr:
		x, ok := value.X.(*ast.Ident)
		if !ok {
			return
		}
		if _, ok := deps[x.String()]; ok {
			for i := range deps[x.String()] {
				deps[x.String()][i].Funcs = append(deps[x.String()][i].Funcs, value.Sel.String())
			}
		}
	case *ast.Ident:
		if _, ok := deps[value.Name]; ok {
			for i := range deps[value.Name] {
				deps[value.Name][i].Funcs = s.fields[value.Name].methods
			}
		}
	}
}
