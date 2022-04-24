package inter

import "context"

type B struct {
	c interface {
		FuncA()
	}
}

func NewB(c *C) *B {
	return &B{c: c}
}

func (b B) FuncA(ctx context.Context) error {
	return nil
}

func (b B) FuncB(_ context.Context) (err error) {
	return nil
}
