package core

import (
	"fmt"
	"sync"
)

type StoreIndex = int

type Store[T any] struct {
	data []T
	mu   sync.Mutex
}

func (s *Store[T]) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.data)
}

func (s *Store[T]) Consume(data T) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = append(s.data, data)
	return len(s.data) - 1
}

func (s *Store[T]) Get(index StoreIndex) (T, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if index < 0 || index >= len(s.data) {
		var res T
		return res, fmt.Errorf("Invalid index %d", index)
	}
	return s.data[index], nil
}

func NewStore[T any]() *Store[T] {
	return &Store[T]{}
}
