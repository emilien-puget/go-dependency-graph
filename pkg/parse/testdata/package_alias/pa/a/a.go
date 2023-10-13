package a

// A pa struct.
type A struct{}

func NewA() *A {
	return &A{}
}

func (a A) FuncFoo(foo string) (bar int, err error) {
	return bar, err
}
