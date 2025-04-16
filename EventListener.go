package Binance

import "sync"

// Event[T] holds a set of listeners taking a T and lets you Emit to all of them.
type Event[T any] struct {
	mu        sync.RWMutex
	listeners map[int]func(T)
	nextID    int
}

// New creates an empty Event[T].
func New[T any]() *Event[T] {
	return &Event[T]{
		listeners: make(map[int]func(T)),
	}
}

// Subscribe adds fn as a listener.  Returns an ID you can use to Unsubscribe.
func (e *Event[T]) Subscribe(fn func(T)) int {
	e.mu.Lock()
	defer e.mu.Unlock()
	id := e.nextID
	e.nextID++
	e.listeners[id] = fn
	return id
}

// Unsubscribe removes the listener with the given ID.
func (e *Event[T]) Unsubscribe(id int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.listeners, id)
}

// Emit calls all current listeners with the given data.
func (e *Event[T]) Emit(data T) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, fn := range e.listeners {
		fn(data)
	}
}
