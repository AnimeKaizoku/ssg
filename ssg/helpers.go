// ssg Project
// Copyright (C) 2021 ALiwoto
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package ssg

import (
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/AnimeKaizoku/ssg/ssg/internal"
	"github.com/AnimeKaizoku/ssg/ssg/rangeValues"
	"github.com/AnimeKaizoku/ssg/ssg/shellUtils"
	"github.com/AnimeKaizoku/ssg/ssg/strongParser"
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

// MakeSureNum will make sure that when you convert `i`
// to string, its length be the exact same as `count`.
// it will append 0 to the left side of the number to do so.
// for example:
// MakeSureNum(5, 8) will return "00000005"
func MakeSureNum(i, count int) string {
	return MakeSureNumCustom(i, count, "0")
}

// MakeSureNumCustom will make sure that when you convert `i`
// to string, its length be the exact same as `count`.
// it will append 0 to the left side of the number to do so.
// for example:
// MakeSureNum(5, 8) will return "00000005"
func MakeSureNumCustom(i, count int, holder string) string {
	s := strconv.Itoa(i)
	final := count - len(s)
	for ; final > 0; final-- {
		s = holder + s
	}

	return s
}

func GetPrettyTimeDuration(d time.Duration, shorten bool) string {
	var result string
	totalSeconds := int(d.Seconds())

	year := totalSeconds / (60 * 60 * 24 * 365)
	totalSeconds -= year * (60 * 60 * 24 * 365)

	month := totalSeconds / (60 * 60 * 24 * 30)
	totalSeconds -= month * (60 * 60 * 24 * 30)

	day := totalSeconds / (60 * 60 * 24)
	totalSeconds -= day * (60 * 60 * 24)

	hour := totalSeconds / (60 * 60)
	totalSeconds -= hour * (60 * 60)

	minute := totalSeconds / 60
	totalSeconds -= minute * 60

	seconds := totalSeconds

	yBool := year > 0
	mBool := month > 0 || yBool
	shorten = !mBool && shorten
	dBool := day > 0 || mBool
	hBool := hour > 0 || dBool
	if yBool {
		result += strconv.Itoa(year) + " year"
		if year > 1 {
			result += "s"
		}
		result += " "
	}
	if mBool {
		result += " " + strconv.Itoa(month) + " month"
		if month > 1 {
			result += "s"
		}
		result += " "
	}
	if dBool {
		result += strconv.Itoa(day)
		if shorten {
			result += "d"
		} else {
			result += " day"
			if day > 1 {
				result += "s"
			}
		}
		result += " "
	}
	if hBool {
		result += strconv.Itoa(hour)
		if shorten {
			result += "h"
		} else {
			result += " hour"
			if hour > 1 {
				result += "s"
			}
		}
		result += " "
	}
	result += strconv.Itoa(minute)
	if shorten {
		result += "m"
	} else {
		result += " minute"
		if minute > 1 {
			result += "s"
		}
	}

	result += " " + strconv.Itoa(seconds)
	if shorten {
		result += "s"
	} else {
		result += " second"
		if seconds > 1 {
			result += "s"
		}
	}
	return result
}

// SplitWhite splits the string with the given separator
// and will remove the white spaces slices from the results
func SplitWhite(s string, separator ...string) []string {
	return internal.SplitWhite(s, separator...)
}

func SplitN(s string, n int, separators ...string) []string {
	return internal.SplitN(s, n, separators...)
}

func SplitSlice(s string, separators []string) []string {
	return internal.SplitSlice(s, separators)
}

func SplitSliceN(s string, separators []string, n int) []string {
	return internal.SplitSliceN(s, separators, n)
}

func SplitSliceNWhite(s string, separators []string, n int) []string {
	return internal.SplitSliceNWhite(s, separators, n)
}

// FixSplit will fix the bullshit bug in the
// Split function (which is not ignoring the spaces between strings).
func FixSplit(myStrings []string) []string {
	return internal.FixSplit(myStrings)
}

// FixSplit will fix the bullshit bug in the
// Split function (which is not ignoring the spaces between strings).
func FixSplitWhite(myStrings []string) []string {
	return internal.FixSplitWhite(myStrings)
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

func ParseConfig(value interface{}, filename string) error {
	return strongParser.ParseConfig(value, filename)
}

func RunCommand(command string) *ExecuteCommandResult {
	return shellUtils.RunCommand(command)
}

func RunPowerShell(command string) *ExecuteCommandResult {
	return shellUtils.RunPowerShell(command)
}

func RunCommandAsync(command string) *ExecuteCommandResult {
	return shellUtils.RunCommandAsync(command)
}

// RunCommandAsyncWithChan runs the given command in another goroutine and returns back
// from the function immediately. It will sends `true` to `finishedChan` the moment
// it's done executing the command.
func RunCommandAsyncWithChan(command string, finishedChan chan bool) *ExecuteCommandResult {
	return shellUtils.RunCommandAsyncWithChan(command, finishedChan)
}

// RunCommandAsyncWithChan runs the given command in another goroutine and returns back
// from the function immediately. It will sends `true` to `finishedChan` the moment
// it's done executing the command.
func RunPowerShellAsyncWithChan(command string, finishedChan chan bool) *ExecuteCommandResult {
	return shellUtils.RunPowerShellAsyncWithChan(command, finishedChan)
}

func ToBool(str string) bool {
	return strongParser.BoolMapping[strings.ToLower(strings.TrimSpace(str))]
}

func ToBase10[T rangeValues.Integer](value T) string {
	return strconv.FormatInt(int64(value), 10)
}

func ToBase16[T rangeValues.Integer](value T) string {
	return strconv.FormatInt(int64(value), 16)
}

func ToBase18[T rangeValues.Integer](value T) string {
	return strconv.FormatInt(int64(value), 18)
}

func ToBase20[T rangeValues.Integer](value T) string {
	return strconv.FormatInt(int64(value), 20)
}

func ToBase28[T rangeValues.Integer](value T) string {
	return strconv.FormatInt(int64(value), 28)
}

func ToBase30[T rangeValues.Integer](value T) string {
	return strconv.FormatInt(int64(value), 30)
}

func ToBase32[T rangeValues.Integer](value T) string {
	return strconv.FormatInt(int64(value), 32)
}

func ToValidIntegerString(value string) string {
	newValue := ""
	for _, current := range value {
		if unicode.IsNumber(current) || current == '-' {
			newValue += string(current)
		}
	}

	return newValue
}

func Title(value string) string {
	return _titleCaser.String(value)
}

func ToInt64(value string) int64 {
	i, _ := strconv.ParseInt(ToValidIntegerString(value), 10, 64)
	return i
}

// ToInteger converts a specified string value to integer.
func ToInteger[T rangeValues.Integer](value string) T {
	var defaultValue T

	switch interface{}(defaultValue).(type) {
	case int:
		return T(ToInt64(value))
	case int64:
		return T(ToInt64(value))
	case int32:
		return T(ToInt32(value))
	case int16:
		return T(ToInt16(value))
	case int8:
		return T(ToInt8(value))
	case uint:
		return T(ToInt64(value))
	case uint64:
		return T(ToInt64(value))
	case uint32:
		return T(ToInt32(value))
	case uint16:
		return T(ToInt16(value))
	case uint8:
		return T(ToInt8(value))
	}

	return defaultValue
}

func ToInt32(value string) int32 {
	i, _ := strconv.ParseInt(ToValidIntegerString(value), 10, 32)
	return int32(i)
}

func ToInt16(value string) int16 {
	i, _ := strconv.ParseInt(ToValidIntegerString(value), 10, 16)
	return int16(i)
}

func ToInt8(value string) int8 {
	i, _ := strconv.ParseInt(ToValidIntegerString(value), 10, 8)
	return int8(i)
}

func IsMixedCase(value string) bool {
	return strings.ToLower(value) != value && strings.ToUpper(value) != value
}

func RecoverPanic() {
	_ = recover()
}

func GetEmptyList[T comparable]() GenericList[T] {
	return &ListW[T]{}
}

func GetListFromArray[T comparable](array []T) GenericList[T] {
	return &ListW[T]{array}
}

func NewEValue[T any](value T) *ExpiringValue[T] {
	return &ExpiringValue[T]{
		_value: value,
		_t:     time.Now(),
	}
}

func NewSafeMap[TKey comparable, TValue any]() *SafeMap[TKey, TValue] {
	return &SafeMap[TKey, TValue]{
		mut:    &sync.RWMutex{},
		values: make(map[TKey]*TValue),
	}
}

func NewAdvancedMap[TKey comparable, TValue any]() *AdvancedMap[TKey, TValue] {
	return &AdvancedMap[TKey, TValue]{
		mut:           &sync.Mutex{},
		values:        make(map[TKey]*TValue),
		sliceKeyIndex: make(map[TKey]int),
	}
}

func NewSafeEMap[TKey comparable, TValue any]() *SafeEMap[TKey, TValue] {
	return &SafeEMap[TKey, TValue]{
		mut:           &sync.RWMutex{},
		values:        make(map[TKey]*ExpiringValue[*TValue]),
		sliceKeyIndex: make(map[TKey]int),
	}
}

func NewNumIdGenerator[T rangeValues.Integer]() *NumIdGenerator[T] {
	return &NumIdGenerator[T]{
		mut: &sync.Mutex{},
	}
}

func IsAllLower(value string) bool {
	return strings.ToLower(value) == value
}

func IsAllUpper(value string) bool {
	return strings.ToUpper(value) == value
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

// Clone function will make an exact copy of the given
// value. It will return nil if the given value is nil,
// so it's guaranteed to execute without panic.
// Usage of this function is up to the users, but it's recommended
// to not use it on large structures, since it will allocate too much
// memory.
func Clone[T any](value *T) *T {
	if value == nil {
		return nil
	}

	tmp := *value
	return &tmp
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
