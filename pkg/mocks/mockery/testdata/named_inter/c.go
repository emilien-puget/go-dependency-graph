package inter

type C struct{}

func NewC() *C {
	return &C{}
}

func (c C) FuncA() {
}
