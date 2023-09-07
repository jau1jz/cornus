package generics_map

import "sync"

func Get[T any](m *sync.Map, key any) (value T, ok bool) {
	load, o := m.Load(key)
	return load.(T), o
}
