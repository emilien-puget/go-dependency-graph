package inter

type daI interface {
	FuncA()
}
type D struct {
	a daI
}

func NewD(a daI) *D {
	return &D{
		a: a,
	}
}

func (d D) FuncA() {
}
