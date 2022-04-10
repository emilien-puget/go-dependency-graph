package fn

import "testdata/fn/pa"

type D struct {
	PaAFuncA func()
}

func NewD(a *pa.A) *D {
	return &D{
		PaAFuncA: a.FuncA,
	}
}

func (d D) FuncA() {
}
