package ext_dep

import hblep "net/http"

type A struct {
	client *hblep.Client
	c      string
}

func NewA(b *hblep.Client, c string) *A {
	return &A{
		client: b,
		c:      c,
	}
}
