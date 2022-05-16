package rangeValues

import (
	"math"
	"math/rand"
)

//---------------------------------------------------------

func (r *IntegerRange[T]) IsInRange(value T) bool {
	return r.Min <= value && r.Max >= value
}

func (r *IntegerRange[T]) IsValueInRange(value *IntegerRange[T]) bool {
	if value == nil {
		return false
	}

	return value == r ||
		(r.Min <= value.Min && r.Max >= value.Max)
}

func (r *IntegerRange[T]) GetRandom() T {
	return T(rand.Int63n(int64(r.Max-r.Min)) + int64(r.Min))
}

//---------------------------------------------------------

func (r *RangeFloat64) IsInRange(value float64) bool {
	if math.IsNaN(value) || r.IsNaN() {
		return false
	}

	return r.Min <= value && r.Max >= value
}

func (r *RangeFloat64) IsNaN() bool {
	return math.IsNaN(r.Min) || math.IsNaN(r.Max)
}

func (r *RangeFloat64) IsValueInRange(value *RangeFloat64) bool {
	if value == nil || r.IsNaN() {
		return false
	}

	return value == r ||
		(r.Min <= value.Min && r.Max >= value.Max)
}

//func (r *RangeFloat64) GetRandom() float64 {
//	return 0
//}

//---------------------------------------------------------
//---------------------------------------------------------
//---------------------------------------------------------
//---------------------------------------------------------
