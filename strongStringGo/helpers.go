// Bot.go Project
// Copyright (C) 2021 Sayan Biswas, ALiwoto
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package strongStringGo

import (
	"strings"

	sin "github.com/ALiwoto/StrongStringGo/strongStringGo/strongInterfaces"
)

// Ss will generate a new StrongString
// with the specified non-encoded string value.
func Ss(s string) StrongString {
	strong := StrongString{}
	strong._setValue(s)
	return strong
}

// SsTrimSpace will trim all whitespaces of the
// specified string value and then will
// generate a new StrongString with that value.
func SsTrimSpace(s string) StrongString {
	s = strings.TrimSpace(s)
	strong := StrongString{}
	strong._setValue(s)
	return strong
}

// SsTrim will trim all leading and
// trailing Unicode code points contained in cutset of the
// specified string value and then will
// generate a new StrongString with that value.
func SsTrim(s, cutset string) StrongString {
	s = strings.Trim(s, cutset)
	strong := StrongString{}
	strong._setValue(s)
	return strong
}

// Qss will generate a new QString
// with the specified non-encoded string value.
func Qss(s string) sin.QString {
	return SsPtr(s)
}

// Sb will generate a new StrongString
// with the specified non-encoded bytes value.
func Sb(b []byte) StrongString {
	return Ss(string(b))
}

// Sb will generate a new StrongString pointer
// with the specified non-encoded bytes value.
func SbPtr(b []byte) *StrongString {
	return SsPtr(string(b))
}

// QSb will generate a new QString
// with the specified non-encoded bytes value.
func Qsb(b []byte) sin.QString {
	return SsPtr(string(b))
}

// SS will generate a new StrongString
// with the specified non-encoded string value.
func SsPtr(s string) *StrongString {
	strong := StrongString{}
	strong._setValue(s)
	return &strong
}

func ToStrSlice(qs []sin.QString) []string {
	tmp := make([]string, len(qs))
	for i, current := range qs {
		tmp[i] = current.GetValue()
	}
	return tmp
}

func ToQSlice(strs []string) []sin.QString {
	tmp := make([]sin.QString, len(strs))
	for i, current := range strs {
		tmp[i] = SsPtr(current)
	}
	return tmp
}

func Split(s string, separator ...string) []string {
	if len(separator) == BaseIndex {
		return nil
	}

	final := s
	for _, myStr := range separator {
		final = strings.ReplaceAll(final, myStr, sepStr)
	}

	return FixSplit(strings.Split(final, sepStr))
}

func SplitN(s string, n int, separator ...string) []string {
	if len(separator) == BaseIndex {
		return nil
	}

	rep := n - BaseOneIndex
	final := s
	done := BaseIndex

	for _, myStr := range separator {
		if done < rep {
			if strings.Contains(final, myStr) {
				final = strings.Replace(final, myStr, sepStr, rep)
				done++
			}
		} else {
			break
		}

	}

	theS := strings.SplitN(final, sepStr, n)

	return FixSplit(theS)
}

func SplitSlice(s string, separator []string) []string {
	if len(separator) == BaseIndex {
		return nil
	}

	final := s
	for _, myStr := range separator {
		final = strings.ReplaceAll(final, myStr, sepStr)
	}

	return FixSplit(strings.Split(final, sepStr))
}

func SplitSliceN(s string, separator []string, n int) []string {
	if len(separator) == BaseIndex {
		return nil
	}

	rep := n - BaseOneIndex
	final := s
	done := BaseIndex

	for _, myStr := range separator {
		if done < rep {
			if strings.Contains(final, myStr) {
				final = strings.Replace(final, myStr, sepStr, rep)
				done++
			}
		} else {
			break
		}

	}

	theS := strings.SplitN(final, sepStr, n)

	return FixSplit(theS)
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

// IsEmpty function will check if the passed-by
// string value is empty or not.
func IsEmpty(s *string) bool {
	return len(*s) == BaseIndex
}

// YesOrNo returns yes if v is true, otherwise no,
func YesOrNo(v bool) string {
	if v {
		return Yes
	} else {
		return No
	}
}

// repairString will go through a string and will repair it.
// In fact, it will escape all special characters between '"' and '"'
// in the specified string.
// For example this value:
//  hello how are you? "I'm = okay"
// will become this:
//  hello how are you? "I'm \= okay"
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

// repairString will go through a string and will repair it.
// In fact, it will escape all special characters between '"' and '"'
// in the specified string.
// For example this value:
//  hello how are you? "I'm = okay"
// will become this:
//  hello how are you? "I'm \= okay"
func repairStringHigh(value *string) *string {
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

		if string(current) == JA_RealStr {
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

			if isSpecialHigh(current) {
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

func isSpecialHigh(r rune) bool {
	switch r {
	case EqualChar, DPointChar,
		BracketOpenChar, BracketcloseChar, CamaChar:
		return true
	default:
		return false
	}
}
