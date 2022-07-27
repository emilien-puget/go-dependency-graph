package ext_dep

import hblep "net/http"

type A struct {
	client *hblep.Client
	c      notAString
	d      string
}

type notAString string

func NewA(b *hblep.Client, c notAString, d string) *A {
	return &A{
		client: b,
		c:      c,
		d:      d,
	}
}
