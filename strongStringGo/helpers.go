// Bot.go Project
// Copyright (C) 2021 Sayan Biswas, ALiwoto
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package strongStringGo

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
