package main

import "sync"

// Broadcast sends data sent from source to all destination channels.
type Broadcast[T any] struct {
	source       chan T
	destinations []chan<- T

	mu sync.Mutex
	wg sync.WaitGroup
}

// The caller is responsible for closing source. When source is closed,
// Broadcast will close all destinations.
func NewBroadcast[T any](source chan T) *Broadcast[T] {
	bc := &Broadcast[T]{
		source,
		make([]chan<- T, 0),
		sync.Mutex{},
		sync.WaitGroup{},
	}

	go func() {
		bc.wg.Add(1)
		for v := range bc.source {
			bc.broadcast(v)
		}
		bc.mu.Lock()
		for _, dest := range bc.destinations {
			close(dest)
		}
		bc.mu.Unlock()
		bc.wg.Done()
	}()
	return bc
}

func (bc *Broadcast[T]) AddDestination() <-chan T {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	ch := make(chan T)
	bc.destinations = append(bc.destinations, ch)
	return ch
}

// Wait for the Broadcast to see that source is closed and to close the
// destinations.
func (bc *Broadcast[T]) Wait() {
	bc.wg.Wait()
}

func (bc *Broadcast[T]) broadcast(v T) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	for _, dest := range bc.destinations {
		dest <- v
	}
}
