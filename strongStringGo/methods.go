// StrongStringGo Project
// Copyright (C) 2021 ALiwoto
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package strongStringGo

import (
	"reflect"
	"strings"
	"sync"
)

//---------------------------------------------------------

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

//---------------------------------------------------------

func (l *ListW[T]) Find(element T) int {
	for i, v := range l._values {
		if v == element {
			return i
		}
	}

	return LIST_INDEX_NOTFOUND
}

func (l *ListW[T]) Count(element T) int {
	count := 0
	for _, v := range l._values {
		if v == element {
			count++
		}
	}

	return count
}

func (l *ListW[T]) Counts(element ...T) int {
	count := 0
	for _, v := range l._values {
		for _, current := range element {
			if v == current {
				count++
			}
		}
	}

	return count
}

func (l *ListW[T]) Contains(element T) bool {
	return l.Find(element) != LIST_INDEX_NOTFOUND
}

func (l *ListW[T]) ContainsAll(elements ...T) bool {
	for _, current := range elements {
		if !l.Contains(current) {
			return false
		}
	}

	return true
}

func (l *ListW[T]) ContainsOne(elements ...T) bool {
	for _, current := range elements {
		if l.Contains(current) {
			return true
		}
	}

	return false
}

func (l *ListW[T]) Change(index int, element T) {
	if index < 0 || index >= len(l._values) {
		return
	}

	l._values[index] = element
}

func (l *ListW[T]) Exists(element T) bool {
	return l.Find(element) != LIST_INDEX_NOTFOUND
}

func (l *ListW[T]) Append(elements ...T) {
	l._values = append(l._values, elements...)
}

func (l *ListW[T]) Add(elements ...T) {
	l._values = append(l._values, elements...)
}

func (l *ListW[T]) RemoveAt(index int) {
	l._values = append(l._values[:index], l._values[index+1:]...)
}

func (l *ListW[T]) RemoveOnce(element T) {
	index := l.Find(element)
	if index != LIST_INDEX_NOTFOUND {
		l.RemoveAt(index)
	}
}

func (l *ListW[T]) RemoveAll(element ...T) {
	var newVal []T
	for _, current := range element {
		for _, v := range l._values {
			if v != current {
				newVal = append(newVal, v)
			}
		}
	}

	l._values = newVal
}

func (l *ListW[T]) Remove(element T) {
	l.RemoveOnce(element)
}

func (l *ListW[T]) AsArray() []T {
	var arr = make([]T, len(l._values))
	copy(arr, l._values)
	return arr
}

func (l *ListW[T]) Clear() {
	l._values = nil
}

func (l *ListW[T]) Get(index int) T {
	return l._values[index]
}

func (l *ListW[T]) IsThreadSafe() bool {
	return true
}

func (l *ListW[T]) IsEmpty() bool {
	return len(l._values) == 0
}

func (l *ListW[T]) Length() int {
	return len(l._values)
}

func (l *ListW[T]) IsValid() bool {
	return len(l._values) > 0
}

//---------------------------------------------------------

func (s *SafeMap[TKey, TValue]) lock() {
	if s.isLocked {
		return
	}

	if s.mut == nil {
		s.mut = &sync.Mutex{}
		s.values = make(map[TKey]*TValue)
	}

	s.isLocked = true
	s.mut.Lock()
}
func (s *SafeMap[TKey, TValue]) unlock() {
	if !s.isLocked {
		return
	}

	if s.mut == nil {
		s.mut = &sync.Mutex{}
	}

	s.isLocked = false
	s.mut.Unlock()
}

func (s *SafeMap[TKey, TValue]) Exists(key TKey) bool {
	s.lock()
	b := len(s.values) != 0 && s.values[key] != nil
	s.unlock()
	return b
}

func (s *SafeMap[TKey, TValue]) Add(key TKey, value *TValue) {
	s.lock()
	s.values[key] = value
	s.unlock()
}

func (s *SafeMap[TKey, TValue]) Delete(key TKey) {
	s.lock()
	delete(s.values, key)
	s.unlock()
}

func (s *SafeMap[TKey, TValue]) Get(key TKey) *TValue {
	s.lock()
	value := s.values[key]
	s.unlock()
	return value
}

func (s *SafeMap[TKey, TValue]) GetValue(key TKey) TValue {
	s.lock()
	value := s.values[key]
	s.unlock()

	if value == nil {
		return s._default
	}

	return *value
}

func (s *SafeMap[TKey, TValue]) SetDefault(value TValue) {
	s._default = value
}

// Set function sets the key of type TKey in this safe map to the value.
// the value should be of type TValue or *TValue, otherwise this function won't
// do anything at all.
func (s *SafeMap[TKey, TValue]) Set(key TKey, value any) {
	correctValue, ok := value.(*TValue)
	if !ok {
		anotherValue, ok := value.(TValue)
		if !ok {
			return
		}

		correctValue = &anotherValue
	}

	s.Add(key, correctValue)
}

// Clear will clear the whole map.
func (s *SafeMap[TKey, TValue]) Clear() {
	s.lock()
	if len(s.values) != 0 {
		s.values = make(map[TKey]*TValue)
	}
	s.unlock()
}

func (s *SafeMap[TKey, TValue]) Length() int {
	s.lock()
	l := len(s.values)
	s.unlock()

	return l
}
func (s *SafeMap[TKey, TValue]) IsEmpty() bool {
	return s.Length() == 0
}
