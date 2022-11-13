package rangeValues

type Integer interface {
	~int64 | ~int | ~int32 | ~int16 | ~int8 | ~uint64 | ~uint | ~uint32 | ~uint16 | ~uint8
}

type IntegerRange[T Integer] struct {
	Min T
	Max T
}

// IntContainer is a container for an integer with some additional flags.
// With standard methods, it will be parsed as "1234:flag1:flag2:flag3"
// as in, first part has to be an integer, and next parts will be considered as
// flag.
type IntContainer[T Integer] struct {
	Value T
	Flags []string
}

type RangeFloat64 struct {
	Min float64
	Max float64
}
