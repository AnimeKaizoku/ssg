// Bot.go Project
// Copyright (C) 2021 Sayan Biswas, ALiwoto
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package strongStringGo

import (
	"reflect"
	"strings"

	tf "github.com/ALiwoto/StrongStringGo/strongStringGo/strongInterfaces"
)

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
	return s._value == nil || len(s._value) == 0
}

// isEqual will check if the passed-by-value in the arg is equal to this
// StrongString or not.
func (s *StrongString) IsEqual(_q tf.QString) bool {
	if reflect.TypeOf(_q) != reflect.TypeOf(s) {
		return _q.GetValue() == s.GetValue()
	}

	_strong, _ok := _q.(*StrongString)
	if !_ok {
		return false
	}
	// check if the length of them are equal or not.
	if len(s._value) != len(_strong._value) {
		//fmt.Println(len(_s._value), len(_strong._value))
		return false
	}
	for i := 0; i < len(s._value); i++ {
		if s._value[i] != _strong._value[i] {
			//fmt.Println(_s._value[i], _strong._value[i])
			return false
		}
	}
	return true
}

// GetIndexV method will give you the rune in _index.
func (s *StrongString) GetIndexV(_index int) rune {
	if s.IsEmpty() {
		return 0
	}

	l := len(s._value)

	if _index >= l || l < 0 {

		return s._value[0]
	}

	return s._value[_index]
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
	for _, str := range values {
		if strings.HasPrefix(s.GetValue(), str) {
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
	for _, str := range values {
		if !strings.HasPrefix(s.GetValue(), str) {
			return false
		}
	}

	return true
}

// Split will split the StrongString using the
// given separator values and will return the results
// as a QString array.
func (s *StrongString) Split(qs ...tf.QString) []tf.QString {
	strs := SplitSlice(s.GetValue(), ToStrSlice(qs))
	return ToQSlice(strs)
}

func (s *StrongString) SplitN(n int, qs ...tf.QString) []tf.QString {
	strs := SplitSliceN(s.GetValue(), ToStrSlice(qs), n)
	return ToQSlice(strs)
}

func (s *StrongString) SplitFirst(qs ...tf.QString) []tf.QString {
	strs := SplitSliceN(s.GetValue(), ToStrSlice(qs), 2)
	return ToQSlice(strs)
}

func (s *StrongString) SplitStr(qs ...string) []tf.QString {
	strs := SplitSlice(s.GetValue(), qs)
	return ToQSlice(strs)
}

func (s *StrongString) SplitStrN(n int, qs ...string) []tf.QString {
	strs := SplitSliceN(s.GetValue(), qs, n)
	return ToQSlice(strs)
}

func (s *StrongString) SplitStrFirst(qs ...string) []tf.QString {
	strs := SplitSliceN(s.GetValue(), qs, 2)
	return ToQSlice(strs)
}

func (s *StrongString) ToQString() tf.QString {
	return s
}

func (s *StrongString) Contains(qs ...tf.QString) bool {
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

func (s *StrongString) ContainsAll(qs ...tf.QString) bool {
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

func (s *StrongString) TrimPrefix(qs ...tf.QString) tf.QString {
	final := s.GetValue()
	for _, current := range qs {
		final = strings.TrimPrefix(final, current.GetValue())
	}

	return SsPtr(final)
}

func (s *StrongString) TrimPrefixStr(qs ...string) tf.QString {
	final := s.GetValue()
	for _, current := range qs {
		final = strings.TrimPrefix(final, current)
	}

	return SsPtr(final)
}

func (s *StrongString) TrimSuffix(qs ...tf.QString) tf.QString {
	final := s.GetValue()
	for _, current := range qs {
		final = strings.TrimSuffix(final, current.GetValue())
	}

	return SsPtr(final)
}

func (s *StrongString) TrimSuffixStr(qs ...string) tf.QString {
	final := s.GetValue()
	for _, current := range qs {
		final = strings.TrimSuffix(final, current)
	}

	return SsPtr(final)
}

func (s *StrongString) Trim(qs ...tf.QString) tf.QString {
	final := s.GetValue()
	for _, current := range qs {
		final = strings.Trim(final, current.GetValue())
	}

	return SsPtr(final)
}

func (s *StrongString) TrimStr(qs ...string) tf.QString {
	final := s.GetValue()
	for _, current := range qs {
		final = strings.Trim(final, current)
	}

	return SsPtr(final)
}

func (s *StrongString) Replace(qs, newS tf.QString) tf.QString {
	return s.ReplaceStr(qs.GetValue(), newS.GetValue())
}

func (s *StrongString) ReplaceStr(qs, newS string) tf.QString {
	final := s.GetValue()
	final = strings.ReplaceAll(final, qs, newS)
	return SsPtr(final)
}

// LockSpecial will lock all the defiend special characters.
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
	// (I mean escaped double quetion mark) is necessary before
	// repairing value.
	final = strings.ReplaceAll(final, BACK_Str, JA_Str)

	// let it repair the string.
	// this function is for repairing these special characters
	// and strings:
	// '=', ':' and "=="
	// it will escape them.
	// if it wasn't for this function, members had to
	// escape all of these bullshits themselves...
	// hahaha, you see, it's actually usefull.
	final = *repairString(&final)

	final = strings.ReplaceAll(final, BACK_Flag, JA_Flag)
	final = strings.ReplaceAll(final, BACK_Equality, JA_Equality)
	final = strings.ReplaceAll(final, BACK_Ddot, JA_Ddot)

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
	final = strings.ReplaceAll(final, JA_Flag, FLAG_PREFIX)
	final = strings.ReplaceAll(final, JA_Str, STR_SIGN)
	final = strings.ReplaceAll(final, JA_Equality, EqualStr)
	final = strings.ReplaceAll(final, JA_Ddot, DdotSign)

	s._value = make([]rune, BaseIndex)
	for _, c := range final {
		if c != BaseIndex {
			s._value = append(s._value, c)
		}
	}
}

// LockSpecial will lock all the defiend special characters.
// This way, you don't actually have to be worry about
// some normal mistakes in spliting strings, cut them out,
// check them. join them, etc...
// WARNING: this method is so dangerous, it's really
// dangerous. we can't say that it's unsafe actually,
// but still it's really dangerous, so if you don't know what the
// fuck are you doing, then please don't use this method.
// this method will not return you a new value, it will effect the
// current value. please consider using it carefully.
func (s *StrongString) LockSpecialHigh() {
	final := s.GetValue()
	// replacing escaped string characters
	// (I mean escaped double quetion mark) is necessary before
	// repairing value.
	final = strings.ReplaceAll(final, BACK_Str, JA_RealStr)

	// let it repair the string.
	// this function is for repairing these special characters
	// and strings:
	// '=', ':' and "=="
	// it will escape them.
	// if it wasn't for this function, members had to
	// escape all of these bullshits themselves...
	// hahaha, you see, it's actually usefull.
	final = *repairStringHigh(&final)

	final = strings.ReplaceAll(final, BACK_Flag, JA_Flag)
	final = strings.ReplaceAll(final, BACK_Equality, JA_Equality)
	final = strings.ReplaceAll(final, BACK_Ddot, JA_Ddot)
	final = strings.ReplaceAll(final, BACK_Cama, JA_Cama)
	final = strings.ReplaceAll(final, BACK_BrOpen, JA_BrOpen)
	final = strings.ReplaceAll(final, BACK_BrClose, JA_BrClose)

	s._value = make([]rune, BaseIndex)
	for _, c := range final {
		if c != BaseIndex {
			s._value = append(s._value, c)
		}
	}
}

// UnlockSpecial will unlock all the defiend special characters.
// it will return them to their normal form.
func (s *StrongString) UnlockSpecialHigh() {
	final := s.GetValue()
	final = strings.ReplaceAll(final, JA_Flag, FLAG_PREFIX)
	final = strings.ReplaceAll(final, JA_RealStr, BACK_Str)
	final = strings.ReplaceAll(final, JA_Equality, EqualStr)
	final = strings.ReplaceAll(final, JA_Ddot, DdotSign)
	final = strings.ReplaceAll(final, JA_Cama, CAMA)
	final = strings.ReplaceAll(final, JA_BrOpen, BracketOpen)
	final = strings.ReplaceAll(final, JA_BrOpen, Bracketclose)

	s._value = make([]rune, BaseIndex)
	for _, c := range final {
		if c != BaseIndex {
			s._value = append(s._value, c)
		}
	}
}
