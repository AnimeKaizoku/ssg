// StrongStringGo Project
// Copyright (C) 2021 ALiwoto
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package strongStringGo

import (
	"regexp"
	"strings"
	"unicode"
)

// Ss will generate a new StrongString
// with the specified non-encoded string value.
func Ss(s string) StrongString {
	_strong := StrongString{}
	_strong._setValue(s)
	return _strong
}

// Qss will generate a new QString
// with the specified non-encoded string value.
func Qss(s string) QString {
	str := Ss(s)
	return &str
}

// Sb will generate a new StrongString
// with the specified non-encoded bytes value.
func Sb(b []byte) StrongString {
	return Ss(string(b))
}

// QSb will generate a new QString
// with the specified non-encoded bytes value.
func Qsb(b []byte) QString {
	str := Ss(string(b))
	return &str
}

// SS will generate a new StrongString
// with the specified non-encoded string value.
func SsPtr(s string) *StrongString {
	strong := StrongString{}
	strong._setValue(s)
	return &strong
}

func ToStrSlice(qs []QString) []string {
	tmp := make([]string, len(qs))
	for i, current := range qs {
		tmp[i] = current.GetValue()
	}
	return tmp
}

func ToQSlice(strs []string) []QString {
	tmp := make([]QString, len(strs))
	for i, current := range strs {
		tmp[i] = SsPtr(current)
	}
	return tmp
}

func Split(s string, separator ...string) []string {
	return SplitSliceN(s, separator, -1)
}

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
	if len(separator) == BaseIndex {
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
	if len(separator) == BaseIndex {
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
	final := make([]string, BaseIndex, cap(myStrings))

	for _, current := range myStrings {
		if !IsEmpty(&current) {
			final = append(final, current)
		}
	}

	return final
}

// FixSplit will fix the bullshit bug in the
// Split function (which is not ignoring the spaces between strings).
func FixSplitWhite(myStrings []string) []string {
	final := make([]string, BaseIndex, cap(myStrings))

	for _, current := range myStrings {
		if strings.TrimSpace(current) != "" {
			final = append(final, current)
		}
	}

	return final
}

// IsEmpty function will check if the passed-by
// string value is empty or not.
func IsEmpty(s *string) bool {
	return s == nil || len(*s) == BaseIndex
}

// AreEqual will check if two string ptr are equal to each other or not.
func AreEqual(s1, s2 *string) bool {
	if s1 == nil && s2 != nil {
		return len(*s2) == 0
	} else if s1 != nil && s2 == nil {
		return len(*s1) == 0
	}

	return s1 == s2 || *s1 == *s2
}

// YesOrNo returns yes if v is true, otherwise no.
func YesOrNo(v bool) string {
	if v {
		return Yes
	} else {
		return No
	}
}

func ToArray(strs ...string) []string {
	return strs
}

func IsAllNumber(str string) bool {
	for _, s := range str {
		if !IsRuneNumber(s) {
			return false
		}
	}

	return true
}

func IsAllNumbers(str ...string) bool {
	for _, ss := range str {
		if !IsAllNumber(ss) {
			return false
		}
	}

	return true
}

func IsRuneNumber(r rune) bool {
	if r <= unicode.MaxLatin1 {
		return '0' <= r && r <= '9'
	}

	return false
}

func repairString(value *string) *string {
	entered := false
	ignoreNext := false
	final := EMPTY
	last := len(*value) - BaseIndex
	next := BaseIndex
	for i, current := range *value {
		if ignoreNext {
			ignoreNext = false
			continue
		}

		if current == CHAR_STR {
			if !entered {
				entered = true
			} else {
				entered = false
			}

			final += string(current)
			continue
		} else {
			if !entered {
				final += string(current)
				continue
			}

			if isSpecial(current) {
				final += BackSlash + string(current)
				continue
			} else {
				if current == LineChar {
					if i != last {
						next = i + BaseOneIndex
						if (*value)[next] == LineChar {
							final += BackSlash +
								string(current) + string(current)
							ignoreNext = true
							continue
						}
					}
				}
			}
		}

		final += string(current)
	}

	return &final
}

func isSpecial(r rune) bool {
	switch r {
	case EqualChar, DPointChar:
		return true
	default:
		return false
	}

}
