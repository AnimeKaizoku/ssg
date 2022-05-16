package rangeValues

type Integer interface {
	int64 | int | int32 | int16 | int8 | uint64 | uint | uint32 | uint16 | uint8
}

type IntegerRange[T Integer] struct {
	Min T
	Max T
}

type RangeFloat64 struct {
	Min float64
	Max float64
}
