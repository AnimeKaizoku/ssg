// Bot.go Project
// Copyright (C) 2021 Sayan Biswas, ALiwoto
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package strongStringGo

// _setValue will set the bytes value of the StrongString.
func (s *StrongString) _setValue(str string) {
	if s._value == nil {
		s._value = make([]rune, BaseIndex)
	}

	for _, current := range str {
		s._value = append(s._value, current)
	}
}
