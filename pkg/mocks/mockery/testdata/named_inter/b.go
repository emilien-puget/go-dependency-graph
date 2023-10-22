package inter

type bcI interface {
	FuncA()
}

type B struct {
	c bcI
}

func NewB(c *C) *B {
	return &B{c: c}
}

func (b B) FuncA() {
}

func (b B) notExported() {
}

func (b B) FuncB() {
}
