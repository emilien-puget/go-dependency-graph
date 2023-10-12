package package_alias

import (
	"testdata/package_alias/pa/a"
	pba "testdata/package_alias/pb/a"
)

type A struct {
	paA interface {
		FuncFoo(foo string) (bar int, err error)
	}
	pbA interface {
		FuncA()
	}
}

func NewA(paA *a.A, paB *pba.A) *A {
	return &A{
		paA: paA,
		pbA: paB,
	}
}
