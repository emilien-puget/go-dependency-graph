package inter

import "testdata/inter/pa"

type D struct {
	a interface {
		FuncFoo(foo string) (bar int, err error)
	}
}

func NewD(a *pa.A) *D {
	return &D{
		a: a,
	}
}

func (d D) FuncA() {
}
