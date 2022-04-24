package inter

import "context"

type A struct {
	b interface {
		FuncA(ctx context.Context) error
		FuncB(context.Context) (err error)
	}
	d interface {
		FuncA()
	}
}

func NewA(b *B, d *D) *A {
	return &A{
		b: b,
		d: d,
	}
}
