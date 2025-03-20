// Copyright 2025 Thom Scheeve. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rand

import (
	"sync"
	"testing"
)

const (
	numTestRuns = 10000
)

func TestIntmnLeftOpen(t *testing.T) {
	const a = -10
	const b = 10
	for i := 0; i < numTestRuns; i++ {
		n := IntmnLeftOpen(a, b)
		if n == a {
			t.Fatalf("IntmnLeftOpen() should be in range (%d,%d]. i: %d n: %d", a, b, i, n)
		}
	}
}

func TestIntmnRightOpen(t *testing.T) {
	const a = -10
	const b = 10
	for i := 0; i < numTestRuns; i++ {
		n := IntmnRightOpen(a, b)
		if n == b {
			t.Fatalf("IntmnRightOpen() should be in range [%d,%d). i: %d n: %d", a, b, i, n)
		}
	}
}

func TestIntmnOpen(t *testing.T) {
	const a = -10
	const b = 10
	for i := 0; i < numTestRuns; i++ {
		n := IntmnOpen(a, b)
		if n == a || n == b {
			t.Fatalf("IntmnOpen() should be in range (%d,%d). i: %d n: %d", a, b, i, n)
		}
	}
}

// Benchmarks

func BenchmarkIntmnClosed1000Threadsafe(b *testing.B) {
	for n := b.N; n > 0; n-- {
		IntmnClosed(-1000, 1000)
	}
}

func BenchmarkIntmnClosed1000ThreadsafeParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			IntmnClosed(-1000, 1000)
		}
	})
}

func BenchmarkIntmnLeftOpen1000Threadsafe(b *testing.B) {
	for n := b.N; n > 0; n-- {
		IntmnLeftOpen(-1000, 1000)
	}
}

func BenchmarkIntmnLeftOpen1000ThreadsafeParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			IntmnLeftOpen(-1000, 1000)
		}
	})
}

func BenchmarkIntmnRightOpen1000Threadsafe(b *testing.B) {
	for n := b.N; n > 0; n-- {
		IntmnRightOpen(-1000, 1000)
	}
}

func BenchmarkIntmnRightOpen1000ThreadsafeParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			IntmnRightOpen(-1000, 1000)
		}
	})
}

func BenchmarkIntmnOpen1000Threadsafe(b *testing.B) {
	for n := b.N; n > 0; n-- {
		IntmnOpen(-1000, 1000)
	}
}

func BenchmarkIntmnOpen1000ThreadsafeParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			IntmnOpen(-1000, 1000)
		}
	})
}

func BenchmarkConcurrent(b *testing.B) {
	const goroutines = 4
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for n := b.N; n > 0; n-- {
				IntmnClosed(-1000, 1000)
			}
		}()
	}
	wg.Wait()
}
