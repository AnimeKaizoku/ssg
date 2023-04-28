package ssg

import (
	"math/rand"
	"sync"
	"time"
)

func (s *SafeEMap[TKey, TValue]) lock() {
	s.mut.Lock()
}

func (s *SafeEMap[TKey, TValue]) unlock() {
	s.mut.Unlock()
}

func (s *SafeEMap[TKey, TValue]) rLock() {
	s.mut.RLock()
}

func (s *SafeEMap[TKey, TValue]) rUnlock() {
	s.mut.RUnlock()
}

func (s *SafeEMap[TKey, TValue]) Exists(key TKey) bool {
	s.rLock()
	_, b := s.values[key]
	s.rUnlock()
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
	s.rLock()
	for _, v := range s.values {
		if v == nil {
			// most likely impossible, this checker is here just for more safety.
			continue
		}

		list.Add(v.GetValue())
	}
	s.rUnlock()

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
	defer s.unlock()

	if s._disabled {
		return
	}

	old := s.values[key]
	if old != nil {
		// don't allocate new memory if we already have the expiring-value struct in
		// the map... just set the new value and reset the time
		old.SetValue(value)
		old.Reset()
		return
	} else {
		s.values[key] = NewEValue(value)
	}

	s.keys = append(s.keys, key)

	// store the index of the map key
	index := len(s.keys) - 1
	s.sliceKeyIndex[key] = index
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

	s.rLock()
	randomIndex := rand.Intn(len(s.keys))
	key := s.keys[randomIndex]
	value := s.values[key]
	s.rUnlock()

	return value.GetValue()
}

func (s *SafeEMap[TKey, TValue]) GetRandomValue() TValue {
	if s.IsEmpty() {
		return s._default
	}

	s.rLock()
	randomIndex := rand.Intn(len(s.keys))
	key := s.keys[randomIndex]
	value := s.values[key]
	s.rUnlock()

	return s.getRealValue(value)
}

func (s *SafeEMap[TKey, TValue]) GetRandomKey() (key TKey, ok bool) {
	if s.IsEmpty() {
		return
	}
	ok = true

	s.rLock()
	key = s.keys[rand.Intn(len(s.keys))]
	s.rUnlock()

	return
}

func (s *SafeEMap[TKey, TValue]) Get(key TKey) *TValue {
	s.rLock()
	value := s.values[key]
	s.rUnlock()
	if value == nil {
		return nil
	}

	return value.GetValue()
}

func (s *SafeEMap[TKey, TValue]) GetValue(key TKey) TValue {
	s.rLock()
	value := s.values[key]
	s.rUnlock()

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
	s.rLock()
	l := len(s.values)
	s.rUnlock()

	return l
}

func (s *SafeEMap[TKey, TValue]) IsEmpty() bool {
	return s.Length() == 0
}

func (s *SafeEMap[TKey, TValue]) ToNormalMap() map[TKey]TValue {
	m := make(map[TKey]TValue)
	s.rLock()
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
	s.rUnlock()

	return m
}

func (s *SafeEMap[TKey, TValue]) ToArray() []TValue {
	var array []TValue
	s.rLock()
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
	s.rUnlock()

	return array
}

func (s *SafeEMap[TKey, TValue]) ToPointerArray() []*TValue {
	var array []*TValue
	s.rLock()
	for _, v := range s.values {
		if v == nil {
			// most likely impossible, this checker is here just for more safety.
			continue
		}

		array = append(array, v._value)
	}
	s.rUnlock()

	return array
}

func (s *SafeEMap[TKey, TValue]) IsThreadSafe() bool {
	return true
}

func (s *SafeEMap[TKey, TValue]) IsValid() bool {
	return s.Length() > 0 && s.values != nil && s.HasValidTimings()
}

// IsDisabled returns true if this map is disabled.
// Disabled maps won't be able to add new values, but will still be able to
// delete/read values.
func (s *SafeEMap[TKey, TValue]) IsDisabled() bool {
	s.rLock()
	defer s.rUnlock()

	return s._disabled
}

// Disable will disable this map, meaning that it won't be able to add new
// values, but will still be able to delete/read values.
func (s *SafeEMap[TKey, TValue]) Disable() {
	s.lock()
	s._disabled = true
	s.unlock()
}

// Enable will enable this map, meaning that it will be able to add new values.
func (s *SafeEMap[TKey, TValue]) Enable() {
	s.lock()
	s._disabled = false
	s.unlock()
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
