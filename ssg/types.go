// ssg Project
// Copyright (C) 2021 ALiwoto
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package ssg

import (
	"hash"
	"sync"
	"time"

	"github.com/AnimeKaizoku/ssg/ssg/rangeValues"
	"github.com/AnimeKaizoku/ssg/ssg/shellUtils"
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

// AdvancedMap is a safe map of type TIndex to pointers of type TValue with
// extra advanced features that you can't find in safe-map types.
// obviously, because of its extra features, it's slightly slower than other
// normal safe-map types.
// this map is completely thread safe and is using internal lock when
// getting and setting variables.
type AdvancedMap[TKey comparable, TValue any] struct {
	mut    *sync.Mutex
	values map[TKey]*TValue
	// keys field is a slice of the map keys used in the map above. We put them in a slice
	// so that we can get a random key by choosing a random index.
	keys []TKey
	// We store the index of each key, so that when we remove an item, we can
	// quickly remove it from the slice above.
	sliceKeyIndex map[TKey]int
	// _default field is the default value this map has to return in GetValue
	// method when the key is not found.
	_default TValue
}

// SafeMap is a safe map of type TIndex to pointers of type TValue.
// this map is completely thread safe and is using internal lock when
// getting and setting variables.
type SafeMap[TKey comparable, TValue any] struct {
	mut    *sync.RWMutex
	values map[TKey]*TValue
	// _default field is the default value this map has to return in GetValue
	// method when the key is not found.
	_default TValue
}

type NumIdGenerator[T rangeValues.Integer] struct {
	current T
	mut     *sync.Mutex
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
	mut             *sync.RWMutex
	checkerMut      *sync.Mutex
	values          map[TKey]*ExpiringValue[*TValue]
	// keys field is a slice of the map keys used in the map above. We put them in a slice
	// so that we can get a random key by choosing a random index.
	keys []TKey
	// We store the index of each key, so that when we remove an item, we can
	// quickly remove it from the slice above.
	sliceKeyIndex map[TKey]int
	_default      TValue

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

type RangeInt = rangeValues.IntegerRange[int]
type RangeInt32 = rangeValues.IntegerRange[int32]
type RangeInt64 = rangeValues.IntegerRange[int64]
type RangeFloat64 = rangeValues.RangeFloat64

type ExecuteCommandResult = shellUtils.ExecuteCommandResult

//type safeList[T any] #TODO: implement safe-list

type StringUniqueIdContainer = UniqueIdContainer[string]
type Int64UniqueIdContainer = UniqueIdContainer[int64]

type UniqueIdContainer[T comparable] interface {
	GetUniqueId() T
	SetAsUniqueId(value T)
	HasValidUniqueId() bool
}

// MetaDataProvider interface provides useful methods for getting/setting
// metadata for a data structure.
type MetaDataProvider interface {
	Get(key string) (string, error)
	GetInt(key string) (int, error)
	GetInt8(key string) (int8, error)
	GetInt16(key string) (int16, error)
	GetInt32(key string) (int32, error)
	GetInt64(key string) (int64, error)
	GetUInt(key string) (uint, error)
	GetUInt8(key string) (uint8, error)
	GetUInt16(key string) (uint16, error)
	GetUInt32(key string) (uint32, error)
	GetUInt64(key string) (uint64, error)
	GetBool(key string) (bool, error)

	GetNoErr(key string) string
	GetIntNoErr(key string) int
	GetInt8NoErr(key string) int8
	GetInt16NoErr(key string) int16
	GetInt32NoErr(key string) int32
	GetInt64NoErr(key string) int64
	GetUIntNoErr(key string) uint
	GetUInt8NoErr(key string) uint8
	GetUInt16NoErr(key string) uint16
	GetUInt32NoErr(key string) uint32
	GetUInt64NoErr(key string) uint64
	GetBoolNoErr(key string) bool

	Set(key, value string)
	SetInt(key string, value int)
	SetInt8(key string, value int8)
	SetInt16(key string, value int16)
	SetInt32(key string, value int32)
	SetInt64(key string, value int64)
	SetUInt(key string, value uint)
	SetUInt8(key string, value uint8)
	SetUInt16(key string, value uint16)
	SetUInt32(key string, value uint32)
	SetUInt64(key string, value uint64)
	SetBool(key string, value bool)
}

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
