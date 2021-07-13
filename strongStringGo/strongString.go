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
