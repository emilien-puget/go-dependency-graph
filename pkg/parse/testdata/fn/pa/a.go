package pa

// A pa struct.
type A struct{}

func NewA() *A {
	return &A{}
}

func (a A) FuncA(toto string) (titi int, err error) {
	return 0, err
}
