package sum

import (
	"matrixbase/pkg/container/types"
	"matrixbase/pkg/container/vector"
	"matrixbase/pkg/encoding"
	"matrixbase/pkg/sql/colexec/aggregation"
	"matrixbase/pkg/vectorize/sum"
	"matrixbase/pkg/vm/mempool"
	"matrixbase/pkg/vm/process"
)

func NewFloat(typ types.Type) *floatSum {
	return &floatSum{typ: typ}
}

func (a *floatSum) Reset() {
	a.cnt = 0
	a.sum = 0
}

func (a *floatSum) Type() types.Type {
	return a.typ
}

func (a *floatSum) Dup() aggregation.Aggregation {
	return &floatSum{typ: a.typ}
}

func (a *floatSum) Fill(sels []int64, vec *vector.Vector) error {
	if n := len(sels); n > 0 {
		switch vec.Typ.Oid {
		case types.T_float32:
			a.sum += float64(sum.Float32SumSels(vec.Col.([]float32), sels))
		case types.T_float64:
			a.sum += sum.Float64SumSels(vec.Col.([]float64), sels)
		}
		a.cnt += int64(n - vec.Nsp.FilterCount(sels))
	} else {
		switch vec.Typ.Oid {
		case types.T_float32:
			a.sum += float64(sum.Float32Sum(vec.Col.([]float32)))
		case types.T_float64:
			a.sum += sum.Float64Sum(vec.Col.([]float64))
		}
		a.cnt += int64(vec.Length() - vec.Nsp.Length())
	}
	return nil
}

func (a *floatSum) Eval() interface{} {
	if a.cnt == 0 {
		return nil
	}
	return a.sum
}

func (a *floatSum) EvalCopy(proc *process.Process) (*vector.Vector, error) {
	data, err := proc.Alloc(8)
	if err != nil {
		return nil, err
	}
	vec := vector.New(a.typ)
	copy(data[mempool.CountSize:], encoding.EncodeFloat64(a.sum))
	vec.Data = data
	vec.Col = encoding.DecodeFloat64Slice(data[mempool.CountSize : mempool.CountSize+8])
	if a.cnt == 0 {
		vec.Nsp.Add(0)
	}
	return vec, nil
}