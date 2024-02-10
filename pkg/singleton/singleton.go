package singleton

import (
	"sync"

	"github.com/tarro-dev/discord-oauth/pkg/option"
)

// Singleton is a thread-safe singleton.
type Singleton[T any] struct {
	value option.Option[T]
	mu    *sync.RWMutex
}

// New returns a new singleton.
func New[T any]() Singleton[T] {
	return Singleton[T]{
		value: option.None[T](),
		mu:    &sync.RWMutex{},
	}
}

// Get returns the value.
func (s *Singleton[T]) Get() option.Option[T] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.value
}

// GetOr returns the value or sets it if it's not set.
func (s *Singleton[T]) GetOr(constructor func() (T, error)) (T, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	v, ok := s.value.Unwrap()
	if ok {
		return v, nil
	}

	value, err := constructor()
	if err != nil {
		return value, err
	}

	s.value = option.Some(value)

	return value, nil
}
