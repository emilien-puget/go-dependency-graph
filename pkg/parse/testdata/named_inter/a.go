package inter

type abI interface {
	FuncA()
	FuncB()
}

type adI interface {
	FuncA()
}

type A struct {
	b abI
	d adI
}

func NewA(b *B, d *D) *A {
	return &A{
		b: b,
		d: d,
	}
}
