package utils

import (
	"fmt"
	"sync"
)

type MutexMap[K comparable, V any] struct {
	m   map[K]V
	mut sync.RWMutex
}

func NewMutexMap[K comparable, V any]() *MutexMap[K, V] {
	return &MutexMap[K, V]{m: make(map[K]V), mut: sync.RWMutex{}}
}

func (m *MutexMap[K, V]) Get(key K) (V, bool) {
	m.mut.RLock()
	defer m.mut.RUnlock()
	v, ok := m.m[key]
	return v, ok
}

func (m *MutexMap[K, V]) Delete(key K) error {
	m.mut.Lock()
	defer m.mut.Unlock()
	_, ok := m.m[key]
	if !ok {
		return fmt.Errorf("not found")
	}
	delete(m.m, key)
	return nil
}

func (m *MutexMap[K, V]) Add(key K, value V) error {
	m.mut.Lock()
	defer m.mut.Unlock()
	_, ok := m.m[key]
	if ok {
		return fmt.Errorf("already exists")
	}
	m.m[key] = value
	return nil
}

func (m *MutexMap[K, V]) AddOrReplace(key K, value V) {
	m.mut.Lock()
	defer m.mut.Unlock()
	m.m[key] = value
}
