package fn

import "context"

type A struct {
	bFuncA func(ctx context.Context) error
	bFuncB func(context.Context) (err error)
	dFuncA func()
	s      SomeFunc
}

type SomeFunc func(ctx context.Context, err error)

func NewA(b *B, d *D, s SomeFunc) *A {
	return &A{
		bFuncA: b.FuncA,
		bFuncB: b.FuncB,
		dFuncA: d.FuncA,
		s:      s,
	}
}
