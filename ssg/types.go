// ssg Project
// Copyright (C) 2021 ALiwoto
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package ssg

import (
	"hash"
	"sync"
	"time"
)

type ExpiringValue[T any] struct {
	_value T
	_t     time.Time
}

// the StrongString used in the program for additional usage.
type StrongString struct {
	_value []rune
}

type ListW[T comparable] struct {
	_values []T
}

// SafeMap is a safe map of type TIndex to pointers of type TValue.
// this map is completely thread safe and is using internal lock when
// getting and setting variables.
type SafeMap[TKey comparable, TValue any] struct {
	mut      *sync.Mutex
	values   map[TKey]*TValue
	_default TValue
}

// SafeEMap is a safe map of type TIndex to pointers of type TValue.
// this map is completely thread safe and is using internal lock when
// getting and setting variables.
// the difference of SafeEMap and SafeMap is that SafeEMap is using a checker loop
// for removing the expired values from itself.
type SafeEMap[TKey comparable, TValue any] struct {
	checkingEnabled bool
	checkInterval   time.Duration
	expiration      time.Duration
	mut             *sync.Mutex
	checkerMut      *sync.Mutex
	values          map[TKey]*ExpiringValue[*TValue]
	_default        TValue

	// onExpired is the event function that will be called when a value with the certain
	// key on the map is expired. this event function will be called in a new goroutine.
	onExpired func(key TKey, value TValue)
}

// EndpointResponse is the generalized form of a response from a HTTP API.
//  T field is already a pointer in this struct, please avoid passing a pointer
// type, to prevent from `Result` field being a pointer to a pointer.
type EndpointResponse[T any] struct {
	Success bool           `json:"success"`
	Result  *T             `json:"result"`
	Error   *EndpointError `json:"error"`
}

// EndpointError is the generalized form of an error from a HTTP API.
// this type should be used as a pointer in EndpointResponse.
type EndpointError struct {
	ErrorCode int    `json:"code"`
	ErrorType int    `json:"error_type"`
	Message   string `json:"message"`
	Origin    string `json:"origin"`
	Date      string `json:"date"`
}

type safeList[T any] SafeMap[int, T]

type GenericList[T comparable] interface {
	BasicObject
	Validator

	Find(element T) int
	Count(element T) int
	Counts(element ...T) int
	Contains(element T) bool
	ContainsAll(elements ...T) bool
	ContainsOne(elements ...T) bool
	Change(index int, element T)
	Exists(element T) bool
	Append(elements ...T)
	Add(elements ...T)
	RemoveAt(index int)
	RemoveOnce(element T)
	RemoveAll(element ...T)
	Remove(element T)
	AsArray() []T
	ToArray() []T
	Clear()
	Get(index int) T
}

type BytesObject interface {
	ToBytes() ([]byte, error)
	Length() int
}

type IntegerRepresent interface {
	ToInt64() int64
	ToUInt64() uint64
	ToInt32() int32
	ToUInt32() uint32
}

type BitsBlocks interface {
	GetBitsSize() int
}

type BasicObject interface {
	IsEmpty() bool
	Length() int
}

type QString interface {
	BasicObject

	GetValue() string
	GetIndexV(int) rune
	IsEqual(QString) bool
	Split(...QString) []QString
	SplitN(int, ...QString) []QString
	SplitFirst(qs ...QString) []QString
	SplitStr(...string) []QString
	SplitStrN(int, ...string) []QString
	SplitStrFirst(...string) []QString
	Contains(...QString) bool
	ContainsStr(...string) bool
	ContainsAll(...QString) bool
	ContainsStrAll(...string) bool
	TrimPrefix(...QString) QString
	TrimPrefixStr(...string) QString
	TrimSuffix(...QString) QString
	TrimSuffixStr(...string) QString
	Trim(qs ...QString) QString
	TrimStr(...string) QString
	Replace(qs, newS QString) QString
	ReplaceStr(qs, newS string) QString
	LockSpecial()
	UnlockSpecial()
	ToBool() bool
}

type Serializer interface {
	Serialize() ([]byte, error)
	StrSerialize() string
}

type Validator interface {
	IsValid() bool
}

type SignatureContainer interface {
	GetSignature() string
	SetSignature(signature string) bool
	SetSignatureByBytes(data []byte) bool
	SetSignatureByFunc(h func() hash.Hash) bool
}
