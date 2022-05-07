package rangeValues

import (
	"math"
	"math/rand"
)

//---------------------------------------------------------

func (r *RangeInt32) IsInRange(value int32) bool {
	return r.Min <= value && r.Max >= value
}

func (r *RangeInt32) IsValueInRange(value *RangeInt32) bool {
	if value == nil {
		return false
	}

	return value == r ||
		(r.Min <= value.Min && r.Max >= value.Max)
}

func (r *RangeInt32) GetRandom() int32 {
	return rand.Int31n(r.Max-r.Min) + r.Min
}

//---------------------------------------------------------

func (r *RangeInt) IsInRange(value int) bool {
	return r.Min <= value && r.Max >= value
}

func (r *RangeInt) IsValueInRange(value *RangeInt) bool {
	if value == nil {
		return false
	}

	return value == r ||
		(r.Min <= value.Min && r.Max >= value.Max)
}

func (r *RangeInt) GetRandom() int {
	return rand.Intn(r.Max-r.Min) + r.Min
}

//---------------------------------------------------------

func (r *RangeInt64) IsInRange(value int64) bool {
	return r.Min <= value && r.Max >= value
}

func (r *RangeInt64) IsValueInRange(value *RangeInt64) bool {
	if value == nil {
		return false
	}

	return value == r ||
		(r.Min <= value.Min && r.Max >= value.Max)
}

func (r *RangeInt64) GetRandom() int64 {
	return rand.Int63n(r.Max-r.Min) + r.Min
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
