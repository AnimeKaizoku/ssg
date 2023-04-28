package ssg

func (s *SafeMap[TKey, TValue]) lock() {
	s.mut.Lock()
}

func (s *SafeMap[TKey, TValue]) unlock() {
	s.mut.Unlock()
}

func (s *SafeMap[TKey, TValue]) rLock() {
	s.mut.RLock()
}

func (s *SafeMap[TKey, TValue]) rUnlock() {
	s.mut.RUnlock()
}

func (s *SafeMap[TKey, TValue]) Exists(key TKey) bool {
	s.rLock()
	_, b := s.values[key]
	s.rUnlock()
	return b
}

func (s *SafeMap[TKey, TValue]) Add(key TKey, value *TValue) {
	s.lock()
	defer s.unlock()

	if s._disabled {
		return
	}
	s.values[key] = value
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
	s.rLock()
	for _, v := range s.values {
		if v == nil {
			array = append(array, s._default)
			continue
		}

		array = append(array, *v)
	}
	s.rUnlock()

	return array
}

func (s *SafeMap[TKey, TValue]) ToPointerArray() []*TValue {
	var array []*TValue
	s.rLock()
	for _, v := range s.values {
		if v == nil {
			// most likely impossible, this checker is here just for more safety.
			continue
		}

		array = append(array, v)
	}
	s.rUnlock()

	return array
}

func (s *SafeMap[TKey, TValue]) ToList() GenericList[*TValue] {
	list := GetEmptyList[*TValue]()
	s.rLock()
	for _, v := range s.values {
		if v == nil {
			// most likely impossible, this checker is here just for more safety.
			continue
		}

		list.Add(v)
	}
	s.rUnlock()

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
	s.rLock()
	value := s.values[key]
	s.rUnlock()
	return value
}

func (s *SafeMap[TKey, TValue]) GetValue(key TKey) TValue {
	s.rLock()
	value := s.values[key]
	s.rUnlock()

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
	s.rLock()
	l := len(s.values)
	s.rUnlock()

	return l
}

func (s *SafeMap[TKey, TValue]) IsEmpty() bool {
	return s.Length() == 0
}

func (s *SafeMap[TKey, TValue]) ToNormalMap() map[TKey]TValue {
	m := make(map[TKey]TValue)
	s.rLock()
	for k, v := range s.values {
		if v == nil {
			m[k] = s._default
			continue
		}

		m[k] = *v
	}
	s.rUnlock()

	return m
}

func (s *SafeMap[TKey, TValue]) IsThreadSafe() bool {
	return true
}

func (s *SafeMap[TKey, TValue]) IsValid() bool {
	return s.Length() > 0
}

// IsDisabled returns true if this map is disabled.
// Disabled maps won't be able to add new values, but will still be able to
// delete/read values.
func (s *SafeMap[TKey, TValue]) IsDisabled() bool {
	s.rLock()
	defer s.rUnlock()

	return s._disabled
}

// Disable will disable this map, meaning that it won't be able to add new
// values, but will still be able to delete/read values.
func (s *SafeMap[TKey, TValue]) Disable() {
	s.lock()
	s._disabled = true
	s.unlock()
}

// Enable will enable this map, meaning that it will be able to add new values.
func (s *SafeMap[TKey, TValue]) Enable() {
	s.lock()
	s._disabled = false
	s.unlock()
}
