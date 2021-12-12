// StrongStringGo Project
// Copyright (C) 2021 ALiwoto
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package strongStringGo

import (
	"reflect"
	"strings"
)

// _setValue will set the bytes value of the StrongString.
func (s *StrongString) _setValue(str string) {
	if s._value == nil {
		s._value = make([]rune, BaseIndex)
	}

	for _, current := range str {
		s._value = append(s._value, current)
	}
}

// GetValue will give you the real value of this StrongString.
func (s *StrongString) GetValue() string {
	return string(s._value)
}

// length method, will give you the length-as-int of this StrongString.
func (s *StrongString) Length() int {
	return len(s._value)
}

// isEmpty will check if this StrongString is empty or not.
func (s *StrongString) IsEmpty() bool {
	return s._value == nil || len(s._value) == BaseIndex
}

// isEqual will check if the passed-by-value in the arg is equal to this
// StrongString or not.
func (s *StrongString) IsEqual(q QString) bool {
	if reflect.TypeOf(q) != reflect.TypeOf(s) {
		return q.GetValue() == s.GetValue()
	}

	strong, _ok := q.(*StrongString)
	if !_ok {
		return false
	}
	// check if the length of them are equal or not.
	if len(s._value) != len(strong._value) {
		return false
	}
	for i := 0; i < len(s._value); i++ {
		if s._value[i] != strong._value[i] {
			return false
		}
	}
	return true
}

// GetIndexV method will give you the rune in index.
// if the passed-by `index` is out of range, it will return
// zero index of the value.
func (s *StrongString) GetIndexV(index int) rune {
	if s.IsEmpty() {
		return BaseIndex
	}

	l := len(s._value)

	if index >= l || l < BaseIndex {
		return s._value[BaseIndex]
	}

	return s._value[index]
}

// HasSuffix will check if at least there is one suffix is
// presents in this StrongString not.
// the StrongString should ends with at least one of these suffixes.
func (s *StrongString) HasSuffix(values ...string) bool {
	for _, str := range values {
		if strings.HasSuffix(s.GetValue(), str) {
			return true
		}
	}

	return false
}

// HasSuffixes will check if all of the suffixes are
// present in this StrongString or not.
// the StrongString should ends with all of these suffixes.
// usage of this method is not recommended, since you can use
// HasSuffix method with only one string (the longest string).
// this way you will just use too much cpu resources.
func (s *StrongString) HasSuffixes(values ...string) bool {
	for _, str := range values {
		if !strings.HasSuffix(s.GetValue(), str) {
			return false
		}
	}

	return true
}

// HasPrefix will check if at least there is one prefix is
// presents in this StrongString or not.
// the StrongString should starts with at least one of these prefixes.
func (s *StrongString) HasPrefix(values ...string) bool {
	v := s.GetValue()
	if len(v) == BaseIndex {
		return false
	}

	for _, str := range values {
		if strings.HasPrefix(v, str) {
			return true
		}
	}

	return false
}

// HasRunePrefix will check if at least there is one prefix is
// presents in this StrongString or not.
// the StrongString should starts with at least one of these prefixes.
func (s *StrongString) HasRunePrefix(values ...rune) bool {
	if len(s._value) == BaseIndex {
		return false
	}

	for _, r := range values {
		if s._value[BaseIndex] == r {
			return true
		}
	}

	return false
}

// HasPrefixes will check if all of the prefixes are
// present in this StrongString or not.
// the StrongString should starts with all of these suffixes.
// usage of this method is not recommended, since you can use
// HasSuffix method with only one string (the longest string).
// this way you will just use too much cpu resources.
func (s *StrongString) HasPrefixes(values ...string) bool {
	v := s.GetValue()
	if len(v) == BaseIndex {
		return false
	}

	for _, str := range values {
		if !strings.HasPrefix(v, str) {
			return false
		}
	}

	return true
}

// HasRunePrefixes will check if all of the prefixes are
// present in this StrongString or not.
// the StrongString should starts with all of these suffixes.
// usage of this method is not recommended, since you can use
// HasSuffix method with only one string (the longest string).
// this way you will just use too much cpu resources.
func (s *StrongString) HasRunePrefixes(values ...rune) bool {
	if len(s._value) == BaseIndex {
		return false
	}

	for _, r := range values {
		if s._value[BaseIndex] != r {
			return false
		}
	}

	return true
}

func (s *StrongString) Split(qs ...QString) []QString {
	strs := SplitSlice(s.GetValue(), ToStrSlice(qs))
	return ToQSlice(strs)
}

func (s *StrongString) SplitN(n int, qs ...QString) []QString {
	strs := SplitSliceN(s.GetValue(), ToStrSlice(qs), n)
	return ToQSlice(strs)
}

func (s *StrongString) SplitFirst(qs ...QString) []QString {
	strs := SplitSliceN(s.GetValue(), ToStrSlice(qs),
		BaseTwoIndex)
	return ToQSlice(strs)
}

func (s *StrongString) SplitStr(qs ...string) []QString {
	strs := SplitSlice(s.GetValue(), qs)
	return ToQSlice(strs)
}

func (s *StrongString) SplitStrN(n int, qs ...string) []QString {
	strs := SplitSliceN(s.GetValue(), qs, n)
	return ToQSlice(strs)
}

func (s *StrongString) SplitStrFirst(qs ...string) []QString {
	strs := SplitSliceN(s.GetValue(), qs, BaseTwoIndex)
	return ToQSlice(strs)
}

func (s *StrongString) ToQString() QString {
	return s
}

func (s *StrongString) Contains(qs ...QString) bool {
	v := s.GetValue()
	for _, current := range qs {
		if strings.Contains(v, current.GetValue()) {
			return true
		}
	}

	return false
}

func (s *StrongString) ContainsStr(str ...string) bool {
	v := s.GetValue()
	for _, current := range str {
		if strings.Contains(v, current) {
			return true
		}
	}

	return false
}

func (s *StrongString) ContainsAll(qs ...QString) bool {
	v := s.GetValue()
	for _, current := range qs {
		if !strings.Contains(v, current.GetValue()) {
			return false
		}
	}

	return true
}

func (s *StrongString) ContainsStrAll(str ...string) bool {
	v := s.GetValue()
	for _, current := range str {
		if !strings.Contains(v, current) {
			return false
		}
	}

	return true
}

func (s *StrongString) TrimPrefix(qs ...QString) QString {
	final := s.GetValue()
	for _, current := range qs {
		final = strings.TrimPrefix(final, current.GetValue())
	}

	return SsPtr(final)
}

func (s *StrongString) TrimPrefixStr(qs ...string) QString {
	final := s.GetValue()
	for _, current := range qs {
		final = strings.TrimPrefix(final, current)
	}

	return SsPtr(final)
}

func (s *StrongString) TrimSuffix(qs ...QString) QString {
	final := s.GetValue()
	for _, current := range qs {
		final = strings.TrimSuffix(final, current.GetValue())
	}

	return SsPtr(final)
}

func (s *StrongString) TrimSuffixStr(qs ...string) QString {
	final := s.GetValue()
	for _, current := range qs {
		final = strings.TrimSuffix(final, current)
	}

	return SsPtr(final)
}

func (s *StrongString) Trim(qs ...QString) QString {
	final := s.GetValue()
	for _, current := range qs {
		final = strings.Trim(final, current.GetValue())
	}

	return SsPtr(final)
}

func (s *StrongString) ToBool() bool {
	return ToBool(s.GetValue())
}

func (s *StrongString) TrimStr(qs ...string) QString {
	final := s.GetValue()
	for _, current := range qs {
		final = strings.Trim(final, current)
	}

	return SsPtr(final)
}

func (s *StrongString) Replace(qs, newS QString) QString {
	return s.ReplaceStr(qs.GetValue(), newS.GetValue())
}

func (s *StrongString) ReplaceStr(qs, newS string) QString {
	final := s.GetValue()
	final = strings.ReplaceAll(final, qs, newS)
	return SsPtr(final)
}

// LockSpecial will lock all the defined special characters.
// This way, you don't actually have to be worry about
// some normal mistakes in spliting strings, cut them out,
// check them. join them, etc...
// WARNING: this method is so dangerous, it's really
// dangerous. we can't say that it's unsafe actually,
// but still it's really dangerous, so if you don't know what the
// fuck are you doing, then please don't use this method.
// this method will not return you a new value, it will effect the
// current value. please consider using it carefully.
func (s *StrongString) LockSpecial() {
	final := s.GetValue()
	// replacing escaped string characters
	// (I mean escaped double question mark) is necessary before
	// repairing value.
	final = strings.ReplaceAll(final, BACK_STR, JA_STR)

	// let it repair the string.
	// this function is for repairing these special characters
	// and strings:
	// '=', ':' and "=="
	// it will escape them.
	// if it wasn't for this function, members had to
	// escape all of these bullshits themselves...
	// hahaha, you see, it's actually usefull.
	final = *repairString(&final)

	final = strings.ReplaceAll(final, BACK_FLAG, JA_FLAG)
	final = strings.ReplaceAll(final, BACK_EQUALITY, JA_EQUALITY)
	final = strings.ReplaceAll(final, BACK_DDOT, JA_DDOT)

	s._value = make([]rune, BaseIndex)
	for _, c := range final {
		if c != BaseIndex {
			s._value = append(s._value, c)
		}
	}
}

// UnlockSpecial will unlock all the defiend special characters.
// it will return them to their normal form.
func (s *StrongString) UnlockSpecial() {
	final := s.GetValue()
	final = strings.ReplaceAll(final, JA_FLAG, FLAG_PREFIX)
	final = strings.ReplaceAll(final, JA_STR, STR_SIGN)
	final = strings.ReplaceAll(final, JA_EQUALITY, EqualStr)
	final = strings.ReplaceAll(final, JA_DDOT, DdotSign)

	s._value = make([]rune, BaseIndex)
	for _, c := range final {
		if c != BaseIndex {
			s._value = append(s._value, c)
		}
	}
}
