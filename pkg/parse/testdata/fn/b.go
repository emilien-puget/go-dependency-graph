package fn

type B struct {
	cFuncA func()
}

func NewB(c *C) *B {
	return &B{
		cFuncA: c.FuncA,
	}
}

func (b B) FuncA() {
}

func (b B) FuncB() {
}
