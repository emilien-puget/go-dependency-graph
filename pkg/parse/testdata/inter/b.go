package inter

type B struct {
	c interface {
		FuncA()
	}
}

func NewB(c *C) *B {
	return &B{c: c}
}

func (b B) FuncA() {
}

func (b B) FuncB() {
}
