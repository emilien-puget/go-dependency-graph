package inter

import "testdata/inter/pa"

type D struct {
	a interface {
		FuncA()
	}
}

func NewD(a *pa.A) *D {
	return &D{
		a: a,
	}
}

func (d D) FuncA() {
}
