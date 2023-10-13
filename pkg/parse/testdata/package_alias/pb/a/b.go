package a

// B pa struct.
type B struct {
	a interface {
		FuncA()
	}
}

func NewB(a *A) *B {
	return &B{
		a: a,
	}
}

func (a B) FuncA() {}
