package inter

type A struct {
	b interface {
		FuncA()
		FuncB()
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
