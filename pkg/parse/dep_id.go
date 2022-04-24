package parse

import "go/ast"

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
