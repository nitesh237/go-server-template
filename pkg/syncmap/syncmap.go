package syncmap

import "sync"

// Map is just wrapper of sync.Map from Go's standard library.
// Map struct uses generics to wrap the sync.Map and gives better type safety.
// All Getters and setters associated with the map will enforce type safety for the key and value.
// Please read the documentation of sync.Map before using this Map.
type Map[K any, V any] struct {
	m sync.Map
}

func (m *Map[K, V]) Get(key K) (value V) {
	val, ok := m.m.Load(key)
	if ok {
		return val.(V)
	}
	return value
}
func (m *Map[K, V]) Delete(key K) { m.m.Delete(key) }

func (m *Map[K, V]) Load(key K) (value V, ok bool) {
	val, ok := m.m.Load(key)
	if ok {
		return val.(V), true
	}
	return value, false
}
func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	val, loaded := m.m.LoadAndDelete(key)
	if loaded {
		return val.(V), true
	}
	return value, false
}
func (m *Map[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	val, loaded := m.m.LoadOrStore(key, value)
	if loaded {
		return val.(V), true
	}
	return value, false
}
func (m *Map[K, V]) Range(f func(key K, value V) (continueRange bool)) {
	m.m.Range(func(anyKey any, anyVal any) bool {
		return f(anyKey.(K), anyVal.(V))
	})
}
func (m *Map[K, V]) Store(key K, value V) {
	m.m.Store(key, value)
}
