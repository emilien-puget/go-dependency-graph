package inter

import "testdata/named_inter/pa"

type daI interface {
	FuncFoo(foo string) (bar int, err error)
}
type D struct {
	a daI
}

func NewD(a *pa.A) *D {
	return &D{
		a: a,
	}
}

func (d D) FuncA() {
}
