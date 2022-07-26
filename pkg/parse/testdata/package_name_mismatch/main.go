package package_name_mismatch

import "gopkg.in/yaml.v3"

type A struct {
	a interface {
		Encode(v interface{}) (err error)
	}
}

func NewA(encoder *yaml.Encoder) *A {
	return &A{
		a: encoder,
	}
}
