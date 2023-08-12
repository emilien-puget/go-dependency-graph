package no_go_mod

type A struct {
	d string
}

func NewA(d string) *A {
	return &A{
		d: d,
	}
}
