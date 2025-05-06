package test

import (
	"sync"
	"testing"
)

type SafeCounter struct {
	mu    sync.Mutex
	value int
}

func (c *SafeCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func TestConcurrentIncrement(t *testing.T) {
	counter := SafeCounter{}
	iterations := 1000

	var wg sync.WaitGroup
	wg.Add(iterations)

	for i := 0; i < iterations; i++ {
		go func() {
			counter.Increment()
			wg.Done()
		}()
	}

	wg.Wait()

	if counter.Value() != iterations {
		t.Errorf("Expected counter value %d, got %d", iterations, counter.Value())
	}
}

func BenchmarkConcurrentIncrement(b *testing.B) {
	counter := SafeCounter{}
	iterations := b.N

	var wg sync.WaitGroup
	wg.Add(iterations)

	for i := 0; i < iterations; i++ {
		go func() {
			counter.Increment()
			wg.Done()
		}()
	}

	wg.Wait()

	if counter.Value() != iterations {
		b.Errorf("Expected counter value %d, got %d", iterations, counter.Value())
	}
}
