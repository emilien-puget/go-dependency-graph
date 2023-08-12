package struct_decl

import (
	"fmt"
	"go/types"
	"strings"
)

type Method struct {
	TypFuc *types.Func
}

func (m Method) String() string {
	ret := m.tupleAsString(m.TypFuc.Type().(*types.Signature).Results())
	if ret == "" {
		return fmt.Sprintf("%s(%s)", m.TypFuc.Name(), m.tupleAsString(m.TypFuc.Type().(*types.Signature).Params()))
	}
	return fmt.Sprintf("%s(%s) (%s)", m.TypFuc.Name(), m.tupleAsString(m.TypFuc.Type().(*types.Signature).Params()), ret)
}

func (m Method) tupleAsString(tuple *types.Tuple) string {
	var params []string
	for i := 0; i < tuple.Len(); i++ {
		param := tuple.At(i)
		params = append(params, fmt.Sprintf("%s %s", param.Name(), param.Type()))
	}
	return strings.Join(params, ", ")
}
