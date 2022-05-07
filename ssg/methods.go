// ssg Project
// Copyright (C) 2021 ALiwoto
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package ssg

import (
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
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

// AsArray returns a copy of the value of this list as an array.
// please do notice that if you make changes to the underlying values of
// that array, change won't be applied to the list.
func (l *ListW[T]) AsArray() []T {
	var arr = make([]T, len(l._values))
	copy(arr, l._values)
	return arr
}

// ToArray is equivalent to AsArray method in any way.
// it returns a copy of the value of this list as an array.
// please do notice that if you make changes to the underlying values of
// that array, change won't be applied to the list.
func (l *ListW[T]) ToArray() []T {
	return l.AsArray()
}

// Clear method clears the whole list.
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

func (s *AdvancedMap[TKey, TValue]) lock() {
	s.mut.Lock()
}
func (s *AdvancedMap[TKey, TValue]) unlock() {
	s.mut.Unlock()
}

func (s *AdvancedMap[TKey, TValue]) Exists(key TKey) bool {
	s.lock()
	b := len(s.values) != 0 && s.values[key] != nil
	s.unlock()
	return b
}

func (s *AdvancedMap[TKey, TValue]) Add(key TKey, value *TValue) {
	s.lock()
	_, exists := s.values[key]
	s.values[key] = value
	if exists {
		s.unlock()
		return
	}

	s.keys = append(s.keys, key)

	// store the index of the map key
	index := len(s.keys) - 1
	s.sliceKeyIndex[key] = index
	s.unlock()
}
func (s *AdvancedMap[TKey, TValue]) GetRandom() *TValue {
	if s.IsEmpty() {
		return nil
	}

	s.lock()
	randomIndex := rand.Intn(len(s.keys))
	key := s.keys[randomIndex]
	value := s.values[key]
	s.unlock()

	return value
}

func (s *AdvancedMap[TKey, TValue]) ForEach(fn func(TKey, *TValue) bool) {
	if fn == nil {
		return
	}
	s.lock()

	for key, value := range s.values {
		if fn(key, value) {
			s.delete(key, false)
		}
	}

	s.unlock()
}

func (s *AdvancedMap[TKey, TValue]) GetRandomValue() TValue {
	if s.IsEmpty() {
		return s._default
	}

	s.lock()
	randomIndex := rand.Intn(len(s.keys))
	key := s.keys[randomIndex]
	value := s.values[key]
	s.unlock()

	if value == nil {
		return s._default
	}

	return *value
}

func (s *AdvancedMap[TKey, TValue]) GetRandomKey() (key TKey, ok bool) {
	if s.IsEmpty() {
		return
	}
	ok = true

	s.lock()
	key = s.keys[rand.Intn(len(s.keys))]
	s.unlock()

	return
}

func (s *AdvancedMap[TKey, TValue]) ToArray() []TValue {
	var array []TValue
	s.lock()
	for _, v := range s.values {
		if v == nil {
			array = append(array, s._default)
			continue
		}

		array = append(array, *v)
	}
	s.unlock()

	return array
}

func (s *AdvancedMap[TKey, TValue]) ToPointerArray() []*TValue {
	var array []*TValue
	s.lock()
	for _, v := range s.values {
		if v == nil {
			// most likely impossible, this checker is here just for more safety.
			continue
		}

		array = append(array, v)
	}
	s.unlock()

	return array
}

func (s *AdvancedMap[TKey, TValue]) ToList() GenericList[*TValue] {
	list := GetEmptyList[*TValue]()
	s.lock()
	for _, v := range s.values {
		if v == nil {
			// most likely impossible, this checker is here just for more safety.
			continue
		}

		list.Add(v)
	}
	s.unlock()

	return list
}

func (s *AdvancedMap[TKey, TValue]) AddList(keyGetter func(*TValue) TKey, elements ...TValue) {
	if len(elements) == 0 || keyGetter == nil {
		return
	}

	for _, current := range elements {
		s.Add(keyGetter(&current), &current)
	}
}

func (s *AdvancedMap[TKey, TValue]) AddPointerList(keyGetter func(*TValue) TKey, elements ...*TValue) {
	if len(elements) == 0 || keyGetter == nil {
		return
	}

	for _, current := range elements {
		s.Add(keyGetter(current), current)
	}
}

func (s *AdvancedMap[TKey, TValue]) delete(key TKey, useLock bool) {
	if useLock {
		s.lock()
	}

	// get index in key slice for key
	index, exists := s.sliceKeyIndex[key]
	if !exists {
		if useLock {
			s.unlock()
		}
		// item does not exist
		return
	}

	delete(s.sliceKeyIndex, key)

	wasLastIndex := len(s.keys)-1 == index

	// remove key from slice of keys
	s.keys[index] = s.keys[len(s.keys)-1]
	s.keys = s.keys[:len(s.keys)-1]

	// we just swapped the last element to another position.
	// so we need to update it's index (if it was not in last position)
	if !wasLastIndex {
		otherKey := s.keys[index]
		s.sliceKeyIndex[otherKey] = index
	}

	delete(s.values, key)
	if useLock {
		s.unlock()
	}
}

func (s *AdvancedMap[TKey, TValue]) Delete(key TKey) {
	s.delete(key, true)
}

func (s *AdvancedMap[TKey, TValue]) Get(key TKey) *TValue {
	s.lock()
	value := s.values[key]
	s.unlock()
	return value
}

func (s *AdvancedMap[TKey, TValue]) GetValue(key TKey) TValue {
	s.lock()
	value := s.values[key]
	s.unlock()

	if value == nil {
		return s._default
	}

	return *value
}

func (s *AdvancedMap[TKey, TValue]) SetDefault(value TValue) {
	s._default = value
}

// Set function sets the key of type TKey in this safe map to the value.
// the value should be of type TValue or *TValue, otherwise this function won't
// do anything at all.
func (s *AdvancedMap[TKey, TValue]) Set(key TKey, value any) {
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
func (s *AdvancedMap[TKey, TValue]) Clear() {
	s.lock()
	if len(s.values) != 0 {
		s.values = make(map[TKey]*TValue)
	}
	s.unlock()
}

func (s *AdvancedMap[TKey, TValue]) Length() int {
	s.lock()
	l := len(s.values)
	s.unlock()

	return l
}

func (s *AdvancedMap[TKey, TValue]) IsEmpty() bool {
	return s.Length() == 0
}

func (s *AdvancedMap[TKey, TValue]) ToNormalMap() map[TKey]TValue {
	m := make(map[TKey]TValue)
	s.lock()
	for k, v := range s.values {
		if v == nil {
			m[k] = s._default
			continue
		}

		m[k] = *v
	}
	s.unlock()

	return m
}

func (s *AdvancedMap[TKey, TValue]) IsThreadSafe() bool {
	return true
}

func (s *AdvancedMap[TKey, TValue]) IsValid() bool {
	return s.Length() > 0
}

//---------------------------------------------------------

func (s *SafeMap[TKey, TValue]) lock() {
	s.mut.Lock()
}
func (s *SafeMap[TKey, TValue]) unlock() {
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

func (s *SafeMap[TKey, TValue]) ForEach(fn func(TKey, *TValue) bool) {
	if fn == nil {
		return
	}

	s.lock()
	for key, value := range s.values {
		if fn(key, value) {
			s.delete(key, false)
		}
	}

	s.unlock()
}

func (s *SafeMap[TKey, TValue]) ToArray() []TValue {
	var array []TValue
	s.lock()
	for _, v := range s.values {
		if v == nil {
			array = append(array, s._default)
			continue
		}

		array = append(array, *v)
	}
	s.unlock()

	return array
}

func (s *SafeMap[TKey, TValue]) ToPointerArray() []*TValue {
	var array []*TValue
	s.lock()
	for _, v := range s.values {
		if v == nil {
			// most likely impossible, this checker is here just for more safety.
			continue
		}

		array = append(array, v)
	}
	s.unlock()

	return array
}

func (s *SafeMap[TKey, TValue]) ToList() GenericList[*TValue] {
	list := GetEmptyList[*TValue]()
	s.lock()
	for _, v := range s.values {
		if v == nil {
			// most likely impossible, this checker is here just for more safety.
			continue
		}

		list.Add(v)
	}
	s.unlock()

	return list
}

func (s *SafeMap[TKey, TValue]) AddList(keyGetter func(*TValue) TKey, elements ...TValue) {
	if len(elements) == 0 || keyGetter == nil {
		return
	}

	for _, current := range elements {
		s.Add(keyGetter(&current), &current)
	}
}

func (s *SafeMap[TKey, TValue]) AddPointerList(keyGetter func(*TValue) TKey, elements ...*TValue) {
	if len(elements) == 0 || keyGetter == nil {
		return
	}

	for _, current := range elements {
		s.Add(keyGetter(current), current)
	}
}

func (s *SafeMap[TKey, TValue]) delete(key TKey, useLock bool) {
	if useLock {
		s.lock()
		defer s.unlock()
	}

	delete(s.values, key)
}

func (s *SafeMap[TKey, TValue]) Delete(key TKey) {
	s.delete(key, true)
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

func (s *SafeMap[TKey, TValue]) ToNormalMap() map[TKey]TValue {
	m := make(map[TKey]TValue)
	s.lock()
	for k, v := range s.values {
		if v == nil {
			m[k] = s._default
			continue
		}

		m[k] = *v
	}
	s.unlock()

	return m
}

func (s *SafeMap[TKey, TValue]) IsThreadSafe() bool {
	return true
}

func (s *SafeMap[TKey, TValue]) IsValid() bool {
	return s.Length() > 0
}

//---------------------------------------------------------

func (s *SafeEMap[TKey, TValue]) lock() {
	s.mut.Lock()
}
func (s *SafeEMap[TKey, TValue]) unlock() {
	s.mut.Unlock()
}

func (s *SafeEMap[TKey, TValue]) Exists(key TKey) bool {
	s.lock()
	b := len(s.values) != 0 && s.values[key] != nil
	s.unlock()
	return b
}

func (s *SafeEMap[TKey, TValue]) AddList(keyGetter func(*TValue) TKey, elements ...TValue) {
	if len(elements) == 0 || keyGetter == nil {
		return
	}

	for _, current := range elements {
		s.Add(keyGetter(&current), &current)
	}
}

func (s *SafeEMap[TKey, TValue]) ToList() GenericList[*TValue] {
	list := GetEmptyList[*TValue]()
	s.lock()
	for _, v := range s.values {
		if v == nil {
			// most likely impossible, this checker is here just for more safety.
			continue
		}

		list.Add(v.GetValue())
	}
	s.unlock()

	return list
}

func (s *SafeEMap[TKey, TValue]) AddPointerList(keyGetter func(*TValue) TKey, elements ...*TValue) {
	if len(elements) == 0 || keyGetter == nil {
		return
	}

	for _, current := range elements {
		s.Add(keyGetter(current), current)
	}
}

func (s *SafeEMap[TKey, TValue]) Add(key TKey, value *TValue) {
	s.lock()
	old := s.values[key]
	if old != nil {
		// don't allocate new memory if we already have the expiring-value struct in
		// the map... just set the new value and reset the time
		old.SetValue(value)
		old.Reset()
		s.unlock()
		return
	} else {
		s.values[key] = NewEValue(value)
	}

	s.keys = append(s.keys, key)

	// store the index of the map key
	index := len(s.keys) - 1
	s.sliceKeyIndex[key] = index
	s.unlock()
}

func (s *SafeEMap[TKey, TValue]) delete(key TKey, useLock bool) {
	if useLock {
		s.lock()
	}
	// get index in key slice for key
	index, exists := s.sliceKeyIndex[key]
	if !exists {
		if useLock {
			s.unlock()
		}
		// item does not exist
		return
	}

	delete(s.sliceKeyIndex, key)

	wasLastIndex := len(s.keys)-1 == index

	// remove key from slice of keys
	s.keys[index] = s.keys[len(s.keys)-1]
	s.keys = s.keys[:len(s.keys)-1]

	// we just swapped the last element to another position.
	// so we need to update it's index (if it was not in last position)
	if !wasLastIndex {
		otherKey := s.keys[index]
		s.sliceKeyIndex[otherKey] = index
	}

	delete(s.values, key)
	if useLock {
		s.unlock()
	}
}

func (s *SafeEMap[TKey, TValue]) Delete(key TKey) {
	s.delete(key, true)
}

func (s *SafeEMap[TKey, TValue]) ForEach(fn func(TKey, *TValue) bool) {
	if fn == nil {
		return
	}
	s.lock()

	var tmpValue *TValue

	for key, value := range s.values {
		if value == nil {
			tmpValue = nil
		} else {
			tmpValue = value.GetValue()
		}

		if fn(key, tmpValue) {
			s.delete(key, false)
		}
	}

	s.unlock()
}

func (s *SafeEMap[TKey, TValue]) GetRandom() *TValue {
	if s.IsEmpty() {
		return nil
	}

	s.lock()
	randomIndex := rand.Intn(len(s.keys))
	key := s.keys[randomIndex]
	value := s.values[key]
	s.unlock()

	return value.GetValue()
}

func (s *SafeEMap[TKey, TValue]) GetRandomValue() TValue {
	if s.IsEmpty() {
		return s._default
	}

	s.lock()
	randomIndex := rand.Intn(len(s.keys))
	key := s.keys[randomIndex]
	value := s.values[key]
	s.unlock()

	return s.getRealValue(value)
}

func (s *SafeEMap[TKey, TValue]) GetRandomKey() (key TKey, ok bool) {
	if s.IsEmpty() {
		return
	}
	ok = true

	s.lock()
	key = s.keys[rand.Intn(len(s.keys))]
	s.unlock()

	return
}

func (s *SafeEMap[TKey, TValue]) Get(key TKey) *TValue {
	s.lock()
	value := s.values[key]
	s.unlock()
	if value == nil {
		return nil
	}

	return value.GetValue()
}

func (s *SafeEMap[TKey, TValue]) GetValue(key TKey) TValue {
	s.lock()
	value := s.values[key]
	s.unlock()

	return s.getRealValue(value)
}

func (s *SafeEMap[TKey, TValue]) SetDefault(value TValue) {
	s._default = value
}

// Set function sets the key of type TKey in this safe map to the value.
// the value should be of type TValue or *TValue, otherwise this function won't
// do anything at all.
func (s *SafeEMap[TKey, TValue]) Set(key TKey, value any) {
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
func (s *SafeEMap[TKey, TValue]) Clear() {
	s.lock()
	if len(s.values) != 0 {
		s.values = make(map[TKey]*ExpiringValue[*TValue])
	}
	s.unlock()
}

func (s *SafeEMap[TKey, TValue]) Length() int {
	s.lock()
	l := len(s.values)
	s.unlock()

	return l
}

func (s *SafeEMap[TKey, TValue]) IsEmpty() bool {
	return s.Length() == 0
}

func (s *SafeEMap[TKey, TValue]) ToNormalMap() map[TKey]TValue {
	m := make(map[TKey]TValue)
	s.lock()
	for k, v := range s.values {
		if v == nil {
			m[k] = s._default
			continue
		}

		realValue := v.GetValue()
		if realValue == nil {
			m[k] = s._default
			continue
		}

		m[k] = *realValue
	}
	s.unlock()

	return m
}

func (s *SafeEMap[TKey, TValue]) ToArray() []TValue {
	var array []TValue
	s.lock()
	for _, v := range s.values {
		if v == nil {
			array = append(array, s._default)
			continue
		}

		realValue := v.GetValue()
		if realValue == nil {
			array = append(array, s._default)
			continue
		}

		array = append(array, *realValue)
	}
	s.unlock()

	return array
}

func (s *SafeEMap[TKey, TValue]) ToPointerArray() []*TValue {
	var array []*TValue
	s.lock()
	for _, v := range s.values {
		if v == nil {
			// most likely impossible, this checker is here just for more safety.
			continue
		}

		array = append(array, v._value)
	}
	s.unlock()

	return array
}

func (s *SafeEMap[TKey, TValue]) IsThreadSafe() bool {
	return true
}

func (s *SafeEMap[TKey, TValue]) IsValid() bool {
	return s.Length() > 0 && s.values != nil && s.HasValidTimings()
}

func (s *SafeEMap[TKey, TValue]) HasValidTimings() bool {
	return s.expiration > time.Microsecond && s.checkInterval > time.Second
}

func (s *SafeEMap[TKey, TValue]) EnableChecking() {
	if s.checkerMut == nil {
		s.checkerMut = &sync.Mutex{}
	}

	// this lock here makes sure that only 1 checkLoop is running at a time.
	s.checkerMut.Lock()
	defer s.checkerMut.Unlock()

	if s.checkingEnabled {
		return
	}

	s.checkingEnabled = true
	go s.checkLoop()
}

func (s *SafeEMap[TKey, TValue]) DisableChecking() {
	if !s.checkingEnabled {
		return
	}

	s.checkingEnabled = false
}

func (s *SafeEMap[TKey, TValue]) IsChecking() bool {
	return s.checkingEnabled
}

func (s *SafeEMap[TKey, TValue]) SetExpiration(duration time.Duration) {
	s.expiration = duration
}

func (s *SafeEMap[TKey, TValue]) SetOnExpired(event func(key TKey, value TValue)) {
	s.onExpired = event
}

func (s *SafeEMap[TKey, TValue]) SetInterval(duration time.Duration) {
	s.checkInterval = duration
}

func (s *SafeEMap[TKey, TValue]) getRealValue(eValue *ExpiringValue[*TValue]) TValue {
	if eValue == nil {
		return s._default
	}

	realValue := eValue.GetValue()
	if realValue == nil {
		return s._default
	}

	return *realValue
}

// DoCheck iterates over the map and checks for expired variables and removes them.
// if the `onExpired` member of the map is set, it will call them.
func (s *SafeEMap[TKey, TValue]) DoCheck() {
	if !s.IsValid() {
		return
	}

	s.lock()
	for i, current := range s.values {
		if current == nil || current.IsExpired(s.expiration) {
			delete(s.values, i)
			if s.onExpired != nil {
				go s.onExpired(i, s.getRealValue(current))
			}
		}
	}
	s.unlock()
}

func (s *SafeEMap[TKey, TValue]) checkLoop() {
	for {
		time.Sleep(s.checkInterval)

		if !s.checkingEnabled {
			return
		}

		if !s.IsValid() {
			if s.values != nil && s.Length() == 0 {
				// don't break out if the map is actually valid but it's empty...
				continue
			}

			s.checkingEnabled = false
			return
		}

		s.lock()
		for i, current := range s.values {
			if current == nil || current.IsExpired(s.expiration) {
				delete(s.values, i)
				if s.onExpired != nil {
					go s.onExpired(i, s.getRealValue(current))
				}
			}
		}
		s.unlock()
	}
}

//---------------------------------------------------------

func (e *ExpiringValue[T]) SetTime(t time.Time) {
	e._t = t
}

func (e *ExpiringValue[T]) GetTime() time.Time {
	return e._t
}

func (e *ExpiringValue[T]) Reset() {
	e.SetTime(time.Now())
}

func (e *ExpiringValue[T]) IsExpired(duration time.Duration) bool {
	return time.Since(e._t) > duration
}

func (e *ExpiringValue[T]) SetValue(value T) {
	e._value = value
}

func (e *ExpiringValue[T]) GetValue() T {
	return e._value
}

//---------------------------------------------------------

func (e *EndpointError) Error() string {
	return strconv.Itoa(e.ErrorCode) + ": " + e.Message
}

//---------------------------------------------------------
