package fn

import "context"

type B struct {
	cFuncA func()
}

func NewB(c *C) *B {
	return &B{
		cFuncA: c.FuncA,
	}
}

func (b B) FuncA(ctx context.Context) error {
	return nil
}

func (b B) FuncB(_ context.Context) (err error) {
	return nil
}
