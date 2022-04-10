package fn

type A struct {
	bFuncA func()
	bFuncB func()
	dFuncA func()
}

func NewA(b *B, d *D) *A {
	return &A{
		bFuncA: b.FuncA,
		bFuncB: b.FuncB,
		dFuncA: d.FuncA,
	}
}
