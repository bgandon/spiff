package dynaml

import (
	"fmt"

	"github.com/cloudfoundry-incubator/spiff/yaml"
)

type AdditionExpr struct {
	A Expression
	B Expression
}

func (e AdditionExpr) Evaluate(binding Binding) (yaml.Node, EvaluationInfo, bool) {
	a, infoa, ok := e.A.Evaluate(binding)
	if !ok {
		return nil, infoa, false
	}

	b, infob, ok := e.B.Evaluate(binding)
	info := infoa.Join(infob)
	if !ok {
		return nil, info, false
	}

	aint, ok := a.Value().(int64)
	if !ok {
		return nil, info, false
	}

	bint, ok := b.Value().(int64)
	if !ok {
		return nil, info, false
	}

	return node(aint + bint), info, true
}

func (e AdditionExpr) String() string {
	return fmt.Sprintf("%s + %s", e.A, e.B)
}
