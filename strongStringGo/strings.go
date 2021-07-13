// Bot.go Project
// Copyright (C) 2021 Sayan Biswas, ALiwoto
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package strongStringGo

import (
	"strings"
)

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
