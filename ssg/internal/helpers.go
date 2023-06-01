package internal

import (
	"regexp"
	"strings"
)

// SplitWhite splits the string with the given separator
// and will remove the white spaces slices from the results
func SplitWhite(s string, separator ...string) []string {
	return SplitSliceNWhite(s, separator, -1)
}

func SplitN(s string, n int, separator ...string) []string {
	return SplitSliceN(s, separator, n)
}

func SplitSlice(s string, separator []string) []string {
	return SplitSliceN(s, separator, -1)
}

func SplitSliceN(s string, separator []string, n int) []string {
	if len(separator) == 0 {
		return []string{s}
	}

	var m string
	for i, f := range separator {
		if i != len(separator)-1 {
			m += regexp.QuoteMeta(f) + OrRegexp
		} else {
			m += regexp.QuoteMeta(f)
		}
	}

	re, err := regexp.Compile(m)
	if err != nil {
		return []string{s}
	}

	return FixSplit(re.Split(s, n))
}

func SplitSliceNWhite(s string, separator []string, n int) []string {
	if len(separator) == 0 {
		return []string{s}
	}

	var m string
	for i, f := range separator {
		if i != len(separator)-1 {
			m += regexp.QuoteMeta(f) + OrRegexp
		} else {
			m += regexp.QuoteMeta(f)
		}
	}

	re, err := regexp.Compile(m)
	if err != nil {
		return []string{s}
	}

	return FixSplitWhite(re.Split(s, n))
}

// FixSplit will fix the bullshit bug in the
// Split function (which is not ignoring the spaces between strings).
func FixSplit(myStrings []string) []string {
	final := make([]string, 0, cap(myStrings))

	for _, current := range myStrings {
		if current != "" {
			final = append(final, current)
		}
	}

	return final
}

// FixSplit will fix the bullshit bug in the
// Split function (which is not ignoring the spaces between strings).
func FixSplitWhite(myStrings []string) []string {
	final := make([]string, 0, cap(myStrings))

	for _, current := range myStrings {
		if strings.TrimSpace(current) != "" {
			final = append(final, current)
		}
	}

	return final
}

func AppendUnique[T comparable](slice []T, target ...T) []T {
	for _, t := range target {
		if !Contains(slice, t) {
			slice = append(slice, t)
		}
	}

	return slice
}

func Contains[T comparable](slice []T, target T) bool {
	for _, t := range slice {
		if t == target {
			return true
		}
	}

	return false
}
